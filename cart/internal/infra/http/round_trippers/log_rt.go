package round_trippers

import (
	"fmt"
	"net/http"
)

type LogRoundTripper struct {
	rt http.RoundTripper
}

func NewLogRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &LogRoundTripper{rt: rt}
}

func (l *LogRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Printf("%s called\n", r.URL.String())

	//return l.rt.RoundTrip(r)
	resp, err := l.rt.RoundTrip(r)
	fmt.Printf("%+v, %v", resp, err)
	return resp, err
}
