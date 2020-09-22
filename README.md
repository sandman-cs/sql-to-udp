# sql-to-udp
Utility to call a SQL and send the results as JSON to a UDP listener like Splunk


Configuration Example:

DefaultSrvList is filled in as default for any blank values in the DbSrvList array to save time entering/updating common information.

```
{
	"DbSrvList": [{
		"DbServer": "serv100.peoplenetonline.com"
	}],
	"DefaultSrvList": {
		"DbDatabase": "dba",
		"DbUsr": "user",
		"DbPwd": "password",
		"DbStatment": "execute get_failed_logins_into_splunk",
		"SysLogSrv": "localhost",
		"SysLogPort": "8514",
		"WorkDelay": 60
	},
	"LocalEcho": true
}
```