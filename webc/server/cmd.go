package server

import (
	"encoding/json"
	"log"

	"github.com/u35s/shell/cmd"

	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/muxtcp"
)

func serverRoute(cmd gcmd.Cmder, f func(*gcmd.CmdMessage)) {
	gcmd.Route(cmd, func(h func(*gcmd.CmdMessage)) func(*gcmd.CmdMessage, ...interface{}) {
		return func(msg *gcmd.CmdMessage, itfc ...interface{}) {
			h(msg)
		}
	}(f))
}

func addServerRoute() {
	serverRoute(&cmd.CmdServer_proxy{}, func(msg *gcmd.CmdMessage) {
		var rev cmd.CmdServer_proxy
		json.Unmarshal(msg.Data, &rev)
		if rev.Err == "" {
			mux := muxtcp.Get(srv.shelld.ID)
			muxtcp.AcceptDial(mux, rev.Net, rev.LAddr)
		}
		log.Printf("[proxy],local addr %v,remote addr %v,err %v", rev.LAddr, rev.RAddr, rev.Err)
	})
}
