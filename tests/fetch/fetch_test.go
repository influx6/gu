package fetch_test

import (
	"testing"

	"github.com/gu-io/gu/shell"
	"github.com/gu-io/gu/shell/cache/localcache"
	"github.com/gu-io/gu/shell/fetch"
	"github.com/gu-io/gu/tests"
)

func TestFetch(t *testing.T) {
	ft := fetch.New(localcache.New("test.v1"))

	_, _, err := ft.Get("http://google.com", shell.DefaultStrategy)
	if err != nil {
		tests.Failed(t, "Should have retrieved giving path 'http://google.com': %q", err)
	}
	tests.Passed(t, "Should have retrieved giving path 'http://google.com'")

	_, _, err = ft.Get("http://google.com", shell.NoCacheStrategy)
	if err != nil {
		tests.Failed(t, "Should have retrieved giving path 'http://google.com': %q", err)
	}
	tests.Passed(t, "Should have retrieved giving path 'http://google.com'")
}
