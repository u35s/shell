package server

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	"github.com/u35s/shell/cmd"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
	"github.com/u35s/gmod/lib/gtime"
	"github.com/u35s/gmod/mods/gconf"
	"github.com/u35s/gmod/mods/gsrvs"
	"github.com/u35s/muxtcp"
)

type uint = uint64

type shelld struct {
	gmod.ModBase
	gsrvs.ServerBase
	serverMsgChannel chan interface{}

	tenSec gtime.Timer
}

func (this *shelld) Init() {
	this.Type = "shelld"

	this.ID = gconf.Uint("shell_daemon_id")
	this.Name = gconf.String("shell_daemon_name")

	this.serverMsgChannel = make(chan interface{}, 1<<16)
	gsrvs.AddToListenAddr(gconf.String("shell_daemon_addr"), gconf.String("shell_daemon_network"))

	this.tenSec.Init(10 * gtime.SecondN)
	addServerRoute()
}

func (this *shelld) Wait() bool {
	gsrvs.EachToListenAddr(func(s *gsrvs.ToListenAddr) {
		if !s.Ok && this.listenTo(s.Net, s.Addr) == nil {
			s.Ok = true
		}
	})
	return true
}

func (this *shelld) listenTo(network, addr string) error {
	listener, err := gnet.Listen(network, addr)
	if err != nil {
		log.Printf("[shelld],listen %v err:%v", addr, err)
		return err
	}
	go gnet.Accept(listener, this.handleConn)
	return nil
}

func (this *shelld) handleConn(conn net.Conn) {
	srv := &gsrvs.ConnectedServer{}
	mux := muxtcp.NewMuxTcp(conn)
	session := mux.Open(1)
	okChan := make(chan uint, 1)
	srv.Agent = gnet.NewAgent(session, gcmd.NewProcessor(), func(itfc interface{}) {
		if msg, ok := itfc.(*gcmd.CmdMessage); ok {
			if msg.GetCmd() == cmd.CmdServer &&
				msg.GetParam() == cmd.CmdServerParam_establishConnection {
				var rev cmd.CmdServer_establishConnection
				json.Unmarshal(msg.Data, &rev)
				srv.ID = rev.ID
				srv.Type, srv.Name = rev.Type, rev.Name
				gsrvs.Add(srv)
				muxtcp.Add(srv.ID, mux)
			} else {
				return
			}
		} else {
			return
		}
		okChan <- 1
		srv.Agent.SetOnMessage(func(itfc interface{}) {
			this.serverMsgChannel <- itfc
		})
	}, func(err error) {
		gsrvs.Remove(srv)
		muxtcp.Remove(srv.ID)
		log.Printf("[shelld],server %v,%v remote addr %v error %v",
			srv.Type, srv.Name, srv.Agent.Conn.RemoteAddr(), err)
	})

	select {
	case <-okChan:
	case <-time.After(2 * time.Second):
		err := errors.New("connection verify time out")
		log.Printf("[shelld],%v", err)
		srv.Agent.Close(err)
	}
}

func (this *shelld) Run() {
	this.dealServerMsg()

	nano := gtime.TimeNano()
	if !this.tenSec.TimeUp(nano) {
		return
	}

	var send cmd.CmdServer_ping
	send.Time = gtime.Time()
	send.From = this.ID
	gsrvs.SendCmdToServer("shellc", "shellc", &send)
	gsrvs.SendCmdToServer("webc", "webc", &send)
}

func (this *shelld) dealServerMsg() {
	for {
		select {
		case msg := <-this.serverMsgChannel:
			gcmd.DeliverMsg(msg)
		default:
			return
		}
	}
}
