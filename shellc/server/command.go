package server

import (
	"bufio"
	"log"
	"net"
	"strings"

	"github.com/u35s/shell/cmd"
	"github.com/u35s/shell/lib"

	"github.com/u35s/gmod/lib/gnet"
	"github.com/u35s/gmod/mods/gconf"
	"github.com/u35s/gmod/mods/gsrvs"
)

type MCmd struct {
	c string
	f func(m MCmdRuler, mid uint, params map[string]string) bool
}

var MCmdArr = []MCmd{
	{"client", MCmdRuler.client},
}

type MCmdRuler struct {
	connStart uint
	connMap   map[uint]net.Conn
}

func (this *MCmdRuler) connWrite(mid uint, w string) {
	conn, ok := this.connMap[mid]
	if ok {
		conn.Write([]byte(w))
	}
}

func (this *MCmdRuler) newConn(conn net.Conn) {
	this.connStart++
	mid := this.connStart
	this.connMap[mid] = conn
	reader := bufio.NewReader(conn)
	for {
		bts, x, err := reader.ReadLine()
		if err == nil {
			strVec := strings.Split(string(bts), " ")
			params := make(map[string]string)
			for i := 1; i < len(strVec); i++ {
				params[lib.Itoa(i)] = strVec[i]
			}
			if len(strVec) > 0 {
				c := strVec[0]
				for i := range MCmdArr {
					if c == MCmdArr[i].c {
						if MCmdArr[i].f(*this, mid, params) {
							log.Printf("[M指令],执行%v", string(bts))
						} else {
							log.Printf("[M指令],执行%v错误", string(bts))
						}
					}
				}
			}
		} else {
			log.Printf("[连接],%v,%v", x, err)
			break
		}
	}
}

func (this *MCmdRuler) init() {
	this.connMap = make(map[uint]net.Conn)
	listener, err := gnet.Listen(gconf.String("telnet_manager_addr"), "tcp")
	if err == nil {
		go gnet.Accept(listener, this.newConn)
	} else {
		log.Printf("[command],init err,%v", err)
	}
}

func (this MCmdRuler) client(mid uint, params map[string]string) bool {
	switch params["1"] {
	case "proxy":
		raddr := params["2"]
		laddr := params["3"]
		net := params["4"]
		var send cmd.CmdServer_proxy
		send.From = srv.ID
		send.LAddr = laddr
		send.RAddr = raddr
		send.Net = net
		gsrvs.SendCmdToServerWithID(srv.shelld.ID, &send)
	default:
		return false
	}

	return true
}
