package main

import (
	"log"

	"github.com/u35s/shell/lib"
	"github.com/u35s/shell/shellc/server"

	"github.com/u35s/glog"
	"github.com/u35s/gmod"
	"github.com/u35s/gmod/lib/utils"
	"github.com/u35s/gmod/mods/gconf"
)

func init() {
	lib.DealArgs()
}

func main() {
	log.SetOutput(glog.Dump())

	gconf.ReadFile("shell.conf")
	defer glog.Flush()
	defer utils.DumpStack("shellc", gconf.Uint("shell_client_id"))

	gmod.Run(
		server.Mod(),
	)
}
