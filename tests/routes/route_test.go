//+build ignore

package routes

import (
	"testing"

	"github.com/gu-io/gu/router"
	"github.com/influx6/faux/tests"
)

func TestRoute(t *testing.T) {
	rm := router.NewRouteManager()

	home := rm.L("/home/*")
	if _, _, pass := home.Test("/home/models/12"); !pass {
		tests.Failed("Should have validated path /home/models/12")
	}
	tests.Passed("Should have validated path /home/models/12")

	index := rm.L("/index/*")
	if _, _, pass := index.Test("/index"); !pass {
		tests.Failed("Should have validated path /index")
	}
	tests.Passed("Should have validated path /index")

	getModel := home.N("/models/*")
	if _, _, pass := getModel.Test("/models"); !pass {
		tests.Failed("Should have validated path /models")
	}
	tests.Passed("Should have validated path /models")

	if _, _, pass := getModel.Test("/models/12"); !pass {
		tests.Failed("Should have validated path /models/12")
	}
	tests.Passed("Should have validated path /models/12")

	id := getModel.N("/:id")
	param, _, pass := id.Test("/12")
	if !pass {
		tests.Failed("Should have validated path /12")
	}
	tests.Passed("Should have validated path /12: %#v", param)

	home.Done(func(px router.PushEvent) {
		tests.Passed("Should have validated path /home/models/12:  /home")
	}).Failed(func(px router.PushEvent) {
		tests.Failed("Should have validated path /home/models/12:  /home")
	})

	getModel.Done(func(px router.PushEvent) {
		tests.Passed("Should have validated path /home/models/12:  /models")
	}).Failed(func(px router.PushEvent) {
		tests.Failed("Should have validated path /home/models/12:  /models")
	})

	id.Done(func(px router.PushEvent) {
		tests.Passed("Should have validated path /home/models/12:  /id")
	}).Failed(func(px router.PushEvent) {
		tests.Failed("Should have validated path /home/models/12:  /id")
	})

	home.Resolve(router.UseLocation("/home/models/12"))
	home.Resolve(router.UseLocationHash("http://thunderhouse.com/#home/models/12"))
}
