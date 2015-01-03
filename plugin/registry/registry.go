// This file will be generated to include all customer specific middlewares
package registry

import (
	"github.com/coopernurse/vulcand/plugin"
	"github.com/coopernurse/vulcand/plugin/cbreaker"
	"github.com/coopernurse/vulcand/plugin/connlimit"
	"github.com/coopernurse/vulcand/plugin/ratelimit"
	"github.com/coopernurse/vulcand/plugin/rewrite"
)

func GetRegistry() *plugin.Registry {
	r := plugin.NewRegistry()

	if err := r.AddSpec(ratelimit.GetSpec()); err != nil {
		panic(err)
	}

	if err := r.AddSpec(connlimit.GetSpec()); err != nil {
		panic(err)
	}

	if err := r.AddSpec(rewrite.GetSpec()); err != nil {
		panic(err)
	}

	if err := r.AddSpec(cbreaker.GetSpec()); err != nil {
		panic(err)
	}

	return r
}
