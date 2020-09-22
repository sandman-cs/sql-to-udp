package main

import (
	"time"
)

// Configuration File Opjects
type configuration struct {
	DbSrvList      []dbSrv
	DefaultSrvList dbSrv
	LocalEcho      bool
	ServerName     string
	AppName        string
	AppVer         string
	MaskMatch      []string
}

type dbSrv struct {
	DbServer   string
	DbUsr      string
	DbPwd      string
	DbDatabase string
	DbStatment string
	SysLogSrv  string
	SysLogPort string
	WorkDelay  time.Duration
}
