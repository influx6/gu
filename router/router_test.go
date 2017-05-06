package router_test

import (
	"net/http"
	"testing"

	"github.com/gu-io/gu/router"
	"github.com/gu-io/gu/router/cache/memorycache"
	"github.com/influx6/faux/tests"
)

type server struct{}

func (server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/collections/count":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	case "/collections":
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func TestRouter(t *testing.T) {
	router := router.NewRouter(server{}, memorycache.New("inmem"))

	res, err := router.Get("/collections", nil)
	if err != nil {
		tests.Failed("Should have sucessesfully made request to %q", "/collections")
	}
	tests.Passed("Should have sucessesfully made request to %q", "/collections")

	if res.StatusCode != http.StatusNoContent {
		tests.Failed("Should have sucessesfully received expected response: %q", res.Status)
	}
	tests.Passed("Should have sucessesfully received expected response: %q", res.Status)

	res, err = router.Get("/collections/count", nil)
	if err != nil {
		tests.Failed("Should have sucessesfully made request to %q", "/collections/count")
	}
	tests.Passed("Should have sucessesfully made request to %q", "/collections/count")

	if res.StatusCode != http.StatusOK {
		tests.Failed("Should have sucessesfully received expected response: %q", res.Status)
	}
	tests.Passed("Should have sucessesfully received expected response: %q", res.Status)

}
