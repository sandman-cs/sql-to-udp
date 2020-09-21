package main

import (
	"time"
)

// Configuration File Opjects
type configuration struct {
	DbSrvList    []dbSrv
	LocalEcho    bool
	ServerName   string
	DbServer     string
	DbUsr        string
	DbPwd        string
	DbDatabase   string
	WorkDelay    time.Duration
	SysLogSrv    string
	SysLogPort   string
	AppName      string
	AppVer       string
	SlackToken   string
	SlackChannel string
	MaskMatch    []string
}

type dbSrv struct {
	DbServer   string
	DbUsr      string
	DbPwd      string
	DbDatabase string
	WorkDelay  time.Duration
}
