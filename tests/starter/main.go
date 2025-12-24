package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

func checkDBConnection(ctx context.Context, host string, port int, user string, password string, dbname string) error {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return fmt.Errorf("database connections: %w", err)
	}
	defer conn.Close(ctx)

	if err = conn.Ping(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	return nil
}

func checkKafkaConnection(brokerAddress string, deadline time.Duration) error {
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka broker: %w", err)
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(deadline)); err != nil {
		return fmt.Errorf("failed to set deadline: %w", err)
	}
	if _, err = conn.Brokers(); err != nil {
		return fmt.Errorf("failed to fetch brokers: %w", err)
	}
	return nil
}

type dbConfig struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
}

var (
	masterDB = dbConfig{
		host:     "postgres-master",
		port:     5432,
		user:     "loms-user",
		password: "loms-password",
		dbName:   "loms_db",
	}
	slavesDB = dbConfig{
		host:     "postgres-replica",
		port:     5432,
		user:     "loms-user",
		password: "loms-password",
		dbName:   "loms_db",
	}
	shared1DB = dbConfig{
		host:     "postgres-comments-shard-1",
		port:     5432,
		user:     "comments-user-1",
		password: "comments-password-1",
		dbName:   "comments_db",
	}
	shared2DB = dbConfig{
		host:     "postgres-comments-shard-2",
		port:     5432,
		user:     "comments-user-2",
		password: "comments-password-2",
		dbName:   "comments_db",
	}
	dbs = []dbConfig{
		masterDB,
		slavesDB,
		shared1DB,
		shared2DB,
	}
	kafkaHosts = []string{"kafka:29092"}
)

func waitPostgres(ctx context.Context) error {
	var (
		attempts = 10
	)
	eg, ctx := errgroup.WithContext(ctx)
	for _, db := range dbs {
		db := db
		eg.Go(func() error {
			if err := retry.Do(func() error {
				return checkDBConnection(ctx, db.host, db.port, db.user, db.password, db.dbName)
			}, retry.Attempts(uint(attempts)), retry.Delay(time.Millisecond*100), retry.Context(ctx)); err != nil {
				return fmt.Errorf("failed to connect postgres (host: %s): %w", db.host, err)
			}
			fmt.Printf("Postgres connection established, host: %s, port: %d\n", db.host, db.port)
			return nil
		})
	}
	return eg.Wait()
}

func waitKafka(ctx context.Context) error {
	var (
		attempts int64 = 10
	)
	eg, ctx := errgroup.WithContext(ctx)
	for _, broker := range kafkaHosts {
		broker := broker
		eg.Go(func() error {
			if err := retry.Do(func() error {
				return checkKafkaConnection(broker, time.Second*10)
			}, retry.Attempts(uint(attempts)), retry.Delay(time.Millisecond*100)); err != nil {
				return fmt.Errorf("failed to connect kafka (host: %s): %w", broker, err)
			}
			fmt.Printf("Kafka connection established, broker: %s\n", broker)
			return nil
		})
	}
	return eg.Wait()
}

func waitResources(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return waitPostgres(ctx)
	})
	eg.Go(func() error {
		return waitKafka(ctx)
	})
	return eg.Wait()
}

type binary struct {
	path string
	env  []string
}

var (
	cart = binary{
		path: "cart/bin/cart",
		env:  []string{"CONFIG_FILE=cart/configs/values_ci.yaml"},
	}
	loms = binary{
		path: "loms/bin/loms",
		env:  []string{"CONFIG_FILE=loms/configs/values_ci.yaml"},
	}
	notifier = binary{
		path: "notifier/bin/notifier",
		env:  []string{"CONFIG_FILE=notifier/configs/values_ci.yaml"},
	}
	comments = binary{
		path: "comments/bin/comments",
		env:  []string{"CONFIG_FILE=comments/configs/values_ci.yaml"},
	}
)

var (
	bins = []binary{cart, loms, comments, notifier}
)

func main() {

	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := waitResources(ctx); err != nil {
		panic("failed to wait resources: " + err.Error())
	}
	hasStart := false
	eg := errgroup.Group{}
	startedProcess := make([]int, 0)
	spLock := sync.Mutex{}
	for _, bin := range bins {
		bin := bin
		eg.Go(func() error {
			if _, err := os.Stat(bin.path); os.IsNotExist(err) {
				fmt.Printf("skip binary: %s\n", bin.path)
				return nil
			}
			cmd := exec.CommandContext(context.Background(), bin.path)
			cmd.Env = append(os.Environ(), bin.env...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.SysProcAttr = &syscall.SysProcAttr{
				//Setsid: true,
			}
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("failed to start %s: %w", bin.path, err)
			}
			fmt.Printf("Started %s\n", bin.path)
			done := make(chan error, 1)
			go func() {
				done <- cmd.Wait()
			}()
			select {
			case err := <-done:
				if err != nil {
					return fmt.Errorf("process %v (PID %d) returned error: %w", bin.path, cmd.Process.Pid, err)
				} else {
					return fmt.Errorf("Process %v (PID %d) finished before 2 second timeout.\n", bin.path, cmd.Process.Pid)
				}
			case <-time.After(2 * time.Second):
			}
			spLock.Lock()
			startedProcess = append(startedProcess, cmd.Process.Pid)
			spLock.Unlock()
			fmt.Printf("Process %v (PID %d) works normally.\n", bin.path, cmd.Process.Pid)
			hasStart = true
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		for _, pid := range startedProcess {
			process, err := os.FindProcess(pid)
			if err != nil {
				fmt.Printf("Failed to find process: %v\n", err)
				continue
			}
			if err = process.Kill(); err != nil {
				fmt.Printf("Failed to kill process: %v\n", err)
				continue
			}
			fmt.Printf("Process: %v killed\n", pid)
		}
		fmt.Printf("RUN BINARIES: %v\n", err)
		os.Exit(1)
	}
	if !hasStart {
		fmt.Println("no binaries to run")
		os.Exit(1)
	}
	fmt.Printf("start complete in %s\n", time.Since(now))
}
