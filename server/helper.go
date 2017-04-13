package server

import (
	"io/ioutil"
	"os"
	"strings"
	"strconv"
)

type Slot struct  {
	Id   string				`json:"id,omitempty"`
	Name string  				`json:"name,omitempty"`
	SlotStartTime int64			`json:"slotStartTime,omitempty"`
	SlotEndTime int64			`json:"slotEndTime,omitempty"`
	CurrentUsers int64 			`json:"currentUsers"`
	MaxUserPerSlot int64	        	`json:"maxUserPerSlot,omitempty"`
}

func getPayload(payloadPath string) string {
	if payloadPath != "" {
		dir, _ := os.Getwd()
		templateData, _ := ioutil.ReadFile(dir + "/server/" + payloadPath)
		return string(templateData)
	} else {
		return ""
	}
}


func SubstParams(sessionMap []string, textData string) string {
	for i, value := range sessionMap {
		if strings.ContainsAny(textData, "${" + strconv.Itoa(i)) {
			textData = strings.Replace(textData, "${" + strconv.Itoa(i) + "}", value, -1)
		}
	}
	return textData
}

