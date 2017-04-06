package router_test

import (
	"net/http"
	"testing"
)

type server struct{}

func (server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "collections/count":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	case "collections":
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func TestRouter(t *testing.T) {

}
