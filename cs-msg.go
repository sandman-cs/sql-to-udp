package main

import (
	"encoding/json"
	"fmt"
	"time"
)

//Start of Msg to JSON Library Code

// StringInfoMsgToJSON ...
func StringInfoMsgToJSON(msg string) string {
	return stringMsgToJSON(msg, "info")
}

// StringErrorMsgToJSON ...
func StringErrorMsgToJSON(msg string) string {
	return stringMsgToJSON(msg, "error")
}

// StringDebugMsgToJSON ...
func StringDebugMsgToJSON(msg string) string {
	return stringMsgToJSON(msg, "debug")
}

// StringWarnMsgToJSON ...
func StringWarnMsgToJSON(msg string) string {
	return stringMsgToJSON(msg, "warning")
}

// StringFatalMsgToJSON ...
func StringFatalMsgToJSON(msg string) string {
	return stringMsgToJSON(msg, "fatal")
}

func stringMsgToJSON(msg string, mType string) string {

	//Load Service Bus UserProperties
	szUserProperties := map[string]interface{}{
		"level":     mType,
		"type":      "msg",
		"msg":       msg,
		"app":       conf.AppName,
		"host":      conf.ServerName,
		"ver":       conf.AppVer,
		"timestamp": fmt.Sprintf(time.Now().String()),
	}

	bytes, _ := json.Marshal(szUserProperties)
	return fmt.Sprintf(string(bytes))

}
