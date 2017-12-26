package server

import (
	"encoding/json"
	"log"

	"github.com/u35s/shell/cmd"
	"github.com/u35s/shell/lib/download"

	"github.com/u35s/gmod"
	"github.com/u35s/gmod/lib/gcmd"
	"github.com/u35s/gmod/lib/gnet"
	"github.com/u35s/gmod/lib/gtime"
	"github.com/u35s/gmod/mods/gconf"
	"github.com/u35s/gmod/mods/gsrvs"
	"github.com/u35s/muxtcp"
)

type uint = uint64

type WebC struct {
	gmod.ModBase
	gsrvs.ServerBase
	shelld gsrvs.ServerBase

	serverMsgChannel chan interface{}

	tenSec  gtime.Timer
	fiveMin gtime.Timer
	web     Web
	dm      download.DownloadManager
}

func (this *WebC) Init() {
	this.Type = "webc"

	this.ID = gconf.Uint("web_client_id")
	this.Name = gconf.String("web_client_name")

	this.shelld.Type = "shelld"
	this.shelld.ID = gconf.Uint("shell_daemon_id")
	this.shelld.Name = gconf.String("shell_daemon_name")
	gsrvs.AddToConnectServerWithID(this.shelld.ID, this.shelld.Type,
		this.shelld.Name, gconf.String("shell_daemon_addr"), gconf.String("shell_daemon_network"))

	this.serverMsgChannel = make(chan interface{}, 1<<16)

	this.tenSec.Init(10 * gtime.SecondN)
	this.fiveMin.Init(5 * gtime.MinuteN)
	this.web.init()
	this.dm.Init()
	addServerRoute()
}

func (this *WebC) Run() {
	this.dealServerMsg()

	if gsrvs.ServerSize() == 0 {
		gsrvs.EachToConnectServer(func(s *gsrvs.ToConnectServer) {
			if !s.Ok && this.connectTo(s.Net, s.Addr, s.Type, s.Name, s.ID) == nil {
				s.Ok = true
				var send cmd.CmdServer_proxy
				send.From = this.ID
				send.RAddr = gconf.String("web_client_proxy_addr")
				send.Net = gconf.String("web_client_proxy_network")
				send.LAddr = gconf.String("web_client_addr")
				gsrvs.SendCmdToServerWithID(s.ID, &send)
			}
		})
	}
	nano := gtime.TimeNano()
	if this.tenSec.TimeUp(nano) {
		this.dm.Timer()
		if this.fiveMin.TimeUp(nano) {
			//this.web.writeIPAddr()
		}
	}
}

func (this *WebC) connectTo(network, addr, tp, name string, id uint) error {
	conn, err := gnet.ConnectTo(network, addr)
	if err != nil {
		log.Printf("connect to %v err:%v", addr, err)
		return err
	}
	mux := muxtcp.NewMuxTcp(conn)
	muxtcp.Add(id, mux)
	session := mux.Open(1)
	srv := &gsrvs.ConnectedServer{ServerBase: gsrvs.ServerBase{Type: tp, Name: name, ID: id}}
	srv.Agent = gnet.NewAgent(session, gcmd.NewProcessor(), func(msg interface{}) {
		this.serverMsgChannel <- msg
	}, func(err error) {
		gsrvs.Remove(srv)
		muxtcp.Remove(srv.ID)
		log.Printf("server %v,%v remote addr %v error %v",
			srv.Type, srv.Name, srv.Agent.Conn.RemoteAddr(), err)
	})
	gsrvs.Add(srv)
	var send cmd.CmdServer_establishConnection
	send.ID = this.ID
	send.Type = this.Type
	send.Name = this.Name
	srv.Agent.SendCmd(&send)
	return nil
}

func (this *WebC) dealServerMsg() {
	for {
		select {
		case msg := <-this.serverMsgChannel:
			gcmd.DeliverMsg(msg)
		default:
			return
		}
	}
}

func (this *WebC) ForwardCmdToServerWithID(sid uint, m gcmd.Cmder) {
	m.Init()
	var send cmd.CmdServer_forward
	send.To = sid
	send.SubCmd = m.GetCmd()
	send.SubParam = m.GetParam()
	bts, err := json.Marshal(m)
	if err == nil {
		send.Data = bts
		gsrvs.SendCmdToServerWithID(this.shelld.ID, &send)
	}
}
