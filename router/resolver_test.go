package router_test

import (
	"testing"

	"github.com/gu-io/gu/router"
	"github.com/influx6/faux/tests"
)

func TestResolver(t *testing.T) {
	rx := router.NewResolver("/:id")
	params, _, state := rx.Test("12")

	if !state {
		tests.Failed("Should have matched giving path")
	}
	tests.Passed("Should have matched giving path")

	val, ok := params["id"]
	if !ok {
		tests.Failed("Should have retrieve parameter :id => %s", val)
	}
	tests.Passed("Should have retrieve parameter :id => %s", val)

	rx.Done(func(px router.PushEvent) {
		tests.Passed("Should have notified with PushEvent %#v", px)
	})

	rx.Failed(func(px router.PushEvent) {
		tests.Failed("Should have notified with PushEvent %#v", px)
	})

	rx.Resolve(router.UseLocation("/12"))
}

func TestRootRoute(t *testing.T) {
	home := router.NewResolver("/*")

	home.Done(func(px router.PushEvent) {
		tests.Passed("Should have notified with PushEvent %#v", px)
	})

	home.Failed(func(px router.PushEvent) {
		tests.Failed("Should have notified with PushEvent %#v", px)
	})

	home.Resolve(router.UseLocationHash("/"))
	home.Resolve(router.UseLocationHash("/#home"))
}

func TestResolverLevels(t *testing.T) {
	home := router.NewResolver("/home/*")
	rx := router.NewResolver("/:id")

	home.Register(rx)

	rx.Done(func(px router.PushEvent) {
		tests.Passed("Should have notified with PushEvent %#v", px)
	})

	rx.Failed(func(px router.PushEvent) {
		tests.Failed("Should have notified with PushEvent %#v", px)
	})

	home.Resolve(router.UseLocation("home/12"))
}

func TestResolverFailed(t *testing.T) {
	rx := router.NewResolver("/:id")
	rx.Done(func(px router.PushEvent) {
		tests.Failed("Should have notified with failed PushEvent %#v", px)
	})

	rx.Failed(func(px router.PushEvent) {
		tests.Passed("Should have notified with failed PushEvent %#v", px)
	})

	rx.Resolve(router.UseLocation("/home/12"))
}
