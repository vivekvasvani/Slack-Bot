package server

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vivekvasvani/Slack-Bot/client"
	"strconv"
)

const (
	application_json = "application/json"
	slackUrl         = "https://hooks.slack.com/services/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"
)

var header = make(map[string]string)

func getSTGStatus(ctx *fasthttp.RequestCtx) {
	var (
		responseV    STFResponse
		available    string
		busy         string
		disconnected string
		indexA       int = 1
		indexB       int = 1
		indexD       int = 1
	)
	header["Content-Type"] = application_json
	header["Accept"] = application_json
	header["Authorization"] = "Bearer bb4a20783d034ce684b7f564bb13f2b15a9c80313d914b60b68f42f0fb75746c"
	response := client.HitRequest("http://devicefarm.hikeapp.com/api/v1/devices", "GET", header, "")
	errUnmarshal := json.Unmarshal(response, &responseV)
	if errUnmarshal != nil {
		fmt.Println(errUnmarshal)
	}

	for _, value := range responseV.Devices {
		//Available
		if value.Owner.Email == "" && value.Present == true {
			//available = available + "{ \"title\": \"Model, OS\", \"value\": \"" + value.Model + "," + value.Version + "\", \"short\": true },"
			available = available + strconv.Itoa(indexA) + ".) " + value.Model + ",\t" + value.Version + ",\t" + value.Serial + "\n"
			indexA++
			continue
		}

		//Busy
		if value.Owner.Email != "" && value.Present == true {
			//busy = busy + "{ \"title\": \"Model, OS, User\", \"value\": \"" + value.Model + "," + value.Version + "," + value.Owner.Email + "\", \"short\": true },"
			busy = busy + strconv.Itoa(indexB) + ".) " + value.Model + ",\t" + value.Version + ",\t" + value.Serial + ",\t" + value.Owner.Email + "\n"
			indexB++
			continue
		}

		//disconnect
		if value.Present == false {
			//disconnected = disconnected + "{ \"title\": \"Model\", \"value\": \"" + value.Model + "\", \"short\": true },"
			disconnected = disconnected + strconv.Itoa(indexD) + ".) " + value.Model + ",\t" + value.Serial + "\n"
			indexD++
			continue
		}
	}
	output := appendToSlice(available, busy, disconnected)
	reader := SubstParams(output, getPayload("slackpayload"))
	//fmt.Println(reader)
	client.HitRequest(slackUrl, "POST", header, reader)
	//ctx.Response.SetBodyString("Ok")
}

func appendToSlice(available, busy, disconnected string) []string {
	output := make([]string, 0)
	if len(available) == 0 {
		output = append(output, "NA")
	} else {
		output = append(output, available[:len(available)-1])
	}

	if len(busy) == 0 {
		output = append(output, "NA")
	} else {
		output = append(output, busy[:len(busy)-1])
	}

	if len(disconnected) == 0 {
		output = append(output, "NA")
	} else {
		output = append(output, disconnected[:len(disconnected)-1])
	}
	return output
}
