package hw5

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	mathrand "math/rand/v2"
	"testing"

	"route256/tests/app/domain"

	toxiproxy "github.com/Shopify/toxiproxy/client"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	"route256/tests/app/clients/cart"
)

type Suite struct {
	suite.Suite

	ctx       context.Context
	randomize *mathrand.Rand
	toxicName string

	totalSkus []int64

	cartClient *cart.Client

	toxiproxyClient *toxiproxy.Client
	proxy           *toxiproxy.Proxy
}

func TestSuite(t *testing.T) {
	suite.RunNamedSuite(t, "Домашнее задание 5", new(Suite))
}

func (s *Suite) BeforeAll(t provider.T) {
	//cfg, err := config.NewConfig()
	//t.Require().NoError(err, "Не удалось создать конфиг")

	s.ctx = context.Background()
	s.toxicName = "latency_down"

	s.cartClient = cart.NewClient("http://localhost:8080")

	random, err := newRand()
	t.Require().NoError(err, "Не удалось создать генератор случайных чисел")
	s.randomize = random

	// SKU из списка должны быть в бд.
	s.totalSkus = []int64{
		1076963, 1148162, 1625903, 2618151, 2956315, 2958025, 3596599, 3618852, 4288068, 4465995, 30816475,
	}

	const (
		proxyName          = "product-service"
		toxiproxyServerURL = "localhost:8474" // "localhost:8474"
		proxyListenURL     = "0.0.0.0:8082"   // "0.0.0.0:8082"
		productServiceURL  = "localhost:8082" // "products:8082"
	)

	s.toxiproxyClient = toxiproxy.NewClient(toxiproxyServerURL)
	s.proxy, err = s.toxiproxyClient.CreateProxy(proxyName, proxyListenURL, productServiceURL)
	t.Require().NoError(err, "Не удалось создать proxy")

	_, err = s.proxy.AddToxic(s.toxicName, "latency", "downstream", 1.0, toxiproxy.Attributes{
		"latency": 600,
	})
	t.Require().NoError(err, "latency down")
}

func (s *Suite) AfterAll(t provider.T) {
	err := s.proxy.Delete()
	t.Require().NoError(err, "Не удалось удалить proxy")
}

func newRand() (*mathrand.Rand, error) {
	var seedBytes [16]byte
	if _, err := rand.Read(seedBytes[:]); err != nil {
		return nil, err
	}

	seed := binary.LittleEndian.Uint64(seedBytes[:8])
	seq := binary.LittleEndian.Uint64(seedBytes[8:])
	src := mathrand.NewPCG(seed, seq)

	return mathrand.New(src), nil
}

func checkSkusMatch(skus []int64, cartItems []domain.CartItem) bool {
	if len(skus) != len(cartItems) {
		return false
	}

	skuMap := make(map[int64]struct{}, len(skus))
	for _, sku := range skus {
		skuMap[sku] = struct{}{}
	}

	for _, item := range cartItems {
		if _, ok := skuMap[int64(item.SKU)]; !ok {
			return false
		}
	}

	return true
}

func (s *Suite) setToxicEnabled(toxicName string, enable bool) error {
	var toxicity float32 = 0.0
	if enable {
		toxicity = 1.0
	}

	toxics, err := s.proxy.Toxics()
	if err != nil {
		return err
	}

	var targetToxic *toxiproxy.Toxic
	for _, toxic := range toxics {
		if toxic.Name == toxicName {
			targetToxic = &toxic
			break
		}
	}

	if targetToxic == nil {
		return errors.New("toxic not found")
	}

	_, err = s.proxy.UpdateToxic(targetToxic.Name, toxicity, targetToxic.Attributes)
	if err != nil {
		return err
	}

	return nil
}
