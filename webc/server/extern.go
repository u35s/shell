package server

import "github.com/u35s/gmod"

var srv *WebC

func Mod() gmod.Moder {
	if srv == nil {
		srv = new(WebC)
	}
	return srv
}
