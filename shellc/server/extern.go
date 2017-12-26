package server

import "github.com/u35s/gmod"

var srv *ShellC

func Mod() gmod.Moder {
	if srv == nil {
		srv = new(ShellC)
	}
	return srv
}
