package cmd

import "github.com/u35s/gmod/lib/gcmd"

const CmdServerParam_establishConnection = 1

type CmdServer_establishConnection struct {
	gcmd.Cmd
	ID   uint
	Type string
	Name string
}

func (m *CmdServer_establishConnection) Init() {
	m.SetBase(CmdServer,
		CmdServerParam_establishConnection)
}

const CmdServerParam_ping = 2

type CmdServer_ping struct {
	gcmd.Cmd
	From uint
	Time uint
}

func (m *CmdServer_ping) Init() {
	m.SetBase(CmdServer,
		CmdServerParam_ping)
}

const CmdServerParam_forward = 3

type CmdServer_forward struct {
	gcmd.Cmd
	To       uint
	SubCmd   uint8
	SubParam uint8
	Data     []byte
}

func (m *CmdServer_forward) Init() {
	m.SetBase(CmdServer,
		CmdServerParam_forward)
}

const CmdServerParam_proxy = 4

type CmdServer_proxy struct {
	gcmd.Cmd
	Net   string
	From  uint
	RAddr string
	LAddr string
	Err   string
}

func (m *CmdServer_proxy) Init() {
	m.SetBase(CmdServer,
		CmdServerParam_proxy)
}
