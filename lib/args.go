package lib

import (
	"flag"
	"os"
	"os/exec"

	"github.com/u35s/glog"
)

func DealArgs() {
	daemon()
}

func daemon() {
	var godaemon = flag.Bool("d", false, "run app as a daemon with -d .")
	if !flag.Parsed() {
		flag.Parse()
	}
	if *godaemon {
		args := os.Args[1:]
		for i := 0; i < len(args); i++ {
			if args[i] == "-d" {
				args = append(args[:i], args[i+1:]...)
				break
			}
		}
		cmd := exec.Command(os.Args[0], args...)
		cmd.Stderr = glog.Dump()
		cmd.Start()
		os.Exit(0)
	}
}
