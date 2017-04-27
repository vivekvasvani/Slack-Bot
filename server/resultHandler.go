package server

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vivekvasvani/Slack-Bot/client"
	"strconv"
	"strings"
)

const (
	application_json = "application/json"
	slackUrl         = "https://hooks.slack.com/services/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"
)

var (
	header  = make(map[string]string)
	output  = make([]string, 0)
	release = make([]string, 0)
)

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
	reader := SubstParams(output, GetPayload("slackpayload"))
	//fmt.Println(reader)
	client.HitRequest(slackUrl, "POST", header, reader)
	//ctx.Response.SetBodyString("Ok")
}

func sendOptions(ctx *fasthttp.RequestCtx) {
	var (
		responseURL string
	)
	//output slice before making call
	output = make([]string, 0)
	release = make([]string, 0)

	//headers for response
	header["Content-Type"] = application_json
	header["Accept"] = application_json

	//get Response URL from request
	responseURL = string(ctx.PostArgs().Peek("response_url"))
	//channelId = string(ctx.PostArgs().Peek("channel_id"))

	//To get command line args
	text := string(ctx.PostArgs().Peek("text"))
	if len(strings.Split(text, " ")) != 2 {
		client.HitRequest(responseURL, "POST", header, " { \"text\" : \" Command line params are missing Expected /androidbuild [origin] [branchname]\" } ")
	}

	//Saving values of origin and branch name
	output = append(output, strings.Split(text, " ")[0])
	output = append(output, strings.Split(text, " ")[1])

	//Send back the response/options to user
	client.HitRequest(responseURL, "POST", header, SubstParams([]string{"\n" + strings.Join(release[:], ",")}, GetPayload("sendOptions.json")))
}

func sendMoreOptions(ctx *fasthttp.RequestCtx) {
	var (
		requestPayloadButton Button
		requestPayloadSelect Select
	)

	header["Content-Type"] = application_json
	header["Accept"] = application_json
	header["Authorization"] = "Basic dml2ZWt2QGhpa2UuaW46ZjE5NjFlMDVhM2QzMTczNGZhNGQwMmI3ZmNlMTQ2ZGQ="

	error := json.Unmarshal(ctx.Request.PostArgs().Peek("payload"), &requestPayloadButton)
	if error != nil {
		fmt.Println(error)
	}

	//To send interactive elements
	if requestPayloadButton.Actions[0].Type == "button" {
		switch {
		case requestPayloadButton.Actions[0].Type == "button" && requestPayloadButton.Actions[0].Value == "yes":
			client.HitRequest(requestPayloadButton.ResponseURL, "POST", header, SubstParams([]string{"\n Already Selected :" + strings.Join(release[:], ",")}, GetPayload("sendOptions.json")))

		case requestPayloadButton.Actions[0].Type == "button" && requestPayloadButton.Actions[0].Value == "theme":
			client.HitRequest(requestPayloadButton.ResponseURL, "POST", header, SubstParams([]string{strings.Join(release[:], ",\n")}, GetPayload("selectTheme.json")))

		case requestPayloadButton.Actions[0].Type == "button" && requestPayloadButton.Actions[0].Value == "done":
			client.HitRequest(requestPayloadButton.ResponseURL, "POST", header, SubstParams(output, GetPayload("finalResponseAfterSubmit.json")))
			jenkinsUrl := "https://jenkins.im.hike.in:8443/view/Release%20(Integration)/job/android_on_demand_build/buildWithParameters?token=AaBbCcDd12345&origin=${0}&Branch=${1}&Build_Flavour=${2}&Theme=${3}&Slack_Notification=${4}"
			output = append(output, "@"+requestPayloadButton.User.Name)
			payload := SubstParams(output, GetPayload("jenkins.json"))
			jenkinsUrl = SubstParams(output, jenkinsUrl)
			client.HitRequest(jenkinsUrl, "POST", header, payload)
			output = make([]string, 0)
			release = make([]string, 0)
		}
	}

	//To unmarshell payload
	if requestPayloadButton.Actions[0].Type == "select" {
		error := json.Unmarshal(ctx.Request.PostArgs().Peek("payload"), &requestPayloadSelect)
		if error != nil {
			fmt.Println(error)
		}

		if requestPayloadSelect.CallbackID == "release" {
			release = append(release, requestPayloadSelect.Actions[0].SelectedOptions[0].Value)
		}

	}

	//To add values on final array, after getting done
	if requestPayloadButton.Actions[0].Type == "select" && requestPayloadSelect.CallbackID == "theme_selection" {
		output = append(output, strings.Join(release[:], ","))
		output = append(output, requestPayloadSelect.Actions[0].SelectedOptions[0].Value)

		client.HitRequest(requestPayloadButton.ResponseURL, "POST", header, SubstParams(output, GetPayload("response.json")))
	}
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
