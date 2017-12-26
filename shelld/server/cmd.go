package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/u35s/shell/cmd"

	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/mods/gsrvs"
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
	serverRoute(&cmd.CmdServer_forward{}, func(msg *gcmd.CmdMessage) {
		var rev cmd.CmdServer_forward
		json.Unmarshal(msg.Data, &rev)
		if srv.ID == rev.To {
			subMsg := new(gcmd.CmdMessage)
			subMsg.Data = rev.Data
			subMsg.SetBase(rev.SubCmd, rev.SubParam)
			gcmd.DeliverMsg(subMsg)
		} else {
			gsrvs.SendCmdToServerWithID(rev.To, &rev)
		}
	})
	serverRoute(&cmd.CmdServer_proxy{}, func(msg *gcmd.CmdMessage) {
		var rev cmd.CmdServer_proxy
		json.Unmarshal(msg.Data, &rev)
		mux := muxtcp.Get(rev.From)
		var err error
		if mux != nil {
			if err = muxtcp.ListenAccept(mux, rev.Net, rev.RAddr); err != nil {
				rev.Err = err.Error()
			}
		} else {
			rev.Err = fmt.Sprintf("%v muxtcp nil", rev.From)
		}
		log.Printf("[proxy],local addr %v,remote addr %v,err %v", rev.LAddr, rev.RAddr, rev.Err)
		gsrvs.SendCmdToServerWithID(rev.From, &rev)
	})
}
