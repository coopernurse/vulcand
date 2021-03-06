package main

import (
	"os"

	"github.com/coopernurse/vulcand/Godeps/_workspace/src/github.com/mailgun/log"
	"github.com/coopernurse/vulcand/plugin/registry"
	"github.com/coopernurse/vulcand/vctl/command"
)

var vulcanUrl string

func main() {
	log.Init([]*log.LogConfig{&log.LogConfig{Name: "console"}})

	cmd := command.NewCommand(registry.GetRegistry())
	err := cmd.Run(os.Args)
	if err != nil {
		log.Errorf("error: %s\n", err)
	}
}
