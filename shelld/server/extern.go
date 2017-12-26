package server

import "github.com/u35s/gmod"

var srv *shelld

func Mod() gmod.Moder {
	if srv == nil {
		srv = new(shelld)
	}
	return srv
}
