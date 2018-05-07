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

type ShellC struct {
	gmod.ModBase
	gsrvs.ServerBase
	shelld gsrvs.ServerBase

	serverMsgChannel chan interface{}

	tenSec gtime.Timer
	m      MCmdRuler
	dm     download.DownloadManager
}

func (this *ShellC) Init() {
	this.Type = "shellc"

	this.ID = gconf.Uint("shell_client_id")
	this.Name = gconf.String("shell_client_name")

	this.shelld.Type = "shelld"
	this.shelld.ID = gconf.Uint("shell_daemon_id")
	this.shelld.Name = gconf.String("shell_daemon_name")
	gsrvs.AddToConnectServerWithID(this.shelld.ID, this.shelld.Type,
		this.shelld.Name, gconf.String("shell_daemon_addr"), gconf.String("shell_daemon_network"))

	this.serverMsgChannel = make(chan interface{}, 1<<16)

	this.tenSec.Init(10 * gtime.SecondN)
	this.m.init()
	this.dm.Init()
	addServerRoute()
}

func (this *ShellC) Run() {
	this.dealServerMsg()

	if gsrvs.ServerSize() == 0 {
		gsrvs.EachToConnectServer(func(s *gsrvs.ToConnectServer) {
			if !s.Ok && this.connectTo(s.Net, s.Addr, s.Type, s.Name, s.ID) == nil {
				s.Ok = true

				var send cmd.CmdServer_proxy
				send.From = this.ID
				send.RAddr = gconf.String("shell_client_proxy_addr")
				send.LAddr = gconf.String("shell_client_addr")
				send.Net = gconf.String("shell_client_proxy_network")
				gsrvs.SendCmdToServerWithID(this.shelld.ID, &send)
			}
		})
	}
}

func (this *ShellC) connectTo(network, addr, tp, name string, id uint) error {
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
	var send cmd.CmdServer_establishConnection
	send.ID = this.ID
	send.Type = this.Type
	send.Name = this.Name
	srv.Agent.SendCmd(&send)
	gsrvs.Add(srv)
	return nil
}

func (this *ShellC) dealServerMsg() {
	for {
		select {
		case msg := <-this.serverMsgChannel:
			gcmd.DeliverMsg(msg)
		default:
			return
		}
	}
}

func (this *ShellC) ForwardCmdToServerWithID(sid uint, m gcmd.Cmder) {
	m.Init()
	var send cmd.CmdServer_forward
	send.To = sid
	send.SubCmd = m.GetCmd()
	send.SubParam = m.GetParam()
	bts, err := json.Marshal(m)
	if err == nil {
		send.Data = bts
		gsrvs.SendCmdToServer(this.shelld.Type, this.shelld.Name, &send)
	}
}
