package server

import (

	"encoding/json"
	"github.com/valyala/fasthttp"
	"bitbucket.org/myntra/admission_bot/client"
	"strings"
	"time"
	"strconv"
)

const (
	application_json ="application/json"
	//slackUrl="https://hooks.slack.com/services/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"
	//slackUrl="https://hooks.slack.com/services/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"
	slackUrl="https://hooks.slack.com/services/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"
)


func postAdmissionResults(ctx *fasthttp.RequestCtx) {

	header := make(map[string]string)
	header["Content-Type"] = application_json
	header["Accept"] = application_json
	response := client.HitRequest("http://mordor.myntra.com/admission/context/eors/slot", "GET", header, "")
	var allSlots = client.GetResponse("$.eors+", response)
	allSlots = strings.Replace(allSlots, "[", "", -1)
	allSlots = strings.Replace(allSlots, "]", "", -1)
	allSlots = strings.Replace(allSlots, "\"", "", -1)
	var slotSplit = strings.Split(allSlots, ",")
	var currentUsers int64
	var totalUsers int64
	output := make([]string, 0)
	for _, v := range slotSplit {
		avaSlot := &Slot{}
		slotResponse := client.HitRequest("http://mordor.myntra.com/admission/context/eors/slot/" + v, "GET", header, "")
		var slot = client.GetResponse("$.eors+", slotResponse)
		//log.Println(slot)
		json.Unmarshal([]byte(slot), avaSlot)
		if(avaSlot.Name=="EORS-Slot-1"){
			avaSlot.CurrentUsers=avaSlot.CurrentUsers+128738
			avaSlot.MaxUserPerSlot=249999+avaSlot.MaxUserPerSlot
		}
		currentUsers = currentUsers + avaSlot.CurrentUsers
		totalUsers = totalUsers + avaSlot.MaxUserPerSlot
		percentage := (float64(avaSlot.CurrentUsers) / float64(avaSlot.MaxUserPerSlot)) * 100
		startTime := strconv.Itoa(time.Unix(avaSlot.SlotStartTime / 1000, 0).Hour()) + ":" + strconv.Itoa(time.Unix(avaSlot.SlotStartTime / 1000, 0).Minute())
		output = append(output, "SlotStartTime=" + startTime + ", CurrentUsers=" + strconv.FormatInt(avaSlot.CurrentUsers, 10)+ " ,MaxUsers=" + strconv.FormatInt(avaSlot.MaxUserPerSlot, 10)+ " ,SlotsPercentage="+ strconv.Itoa(int(percentage)) + "%")
		//fmt.Println("SlotStartTime="+startTime+" ,CurrentUsers="+strconv.FormatInt(avaSlot.CurrentUsers, 10),",MaxUsers="+strconv.FormatInt(avaSlot.MaxUserPerSlot, 10),",SlotsPercentage=",strconv.Itoa(int(percentage))+"%")
	}
	ttlPercentage := (float64(currentUsers) / float64(totalUsers)) * 100
	output = append(output, "Total:=" + strconv.Itoa(int(currentUsers))+ " / TotalPercentage:=" + strconv.Itoa(int(ttlPercentage)) + "%")

	//fmt.Println("Total:=="+strconv.Itoa(int(currentUsers)),"TotalPercentage="+strconv.Itoa(int(ttlPercentage))+"%")
	//fmt.Println(output)

	reader := SubstParams(output,getPayload("slackpayload"))
	//log.Println(reader)

	client.HitRequest(slackUrl, "POST", header, reader)

	//api := slack.New("xoxp-2151807574-3267458998-3324373337-2b7a5a")
	//err := api.ChatPostMessage("hrd_admission_summary", "```" + strings.Join(output, "\n") + "```", nil)
	//if err != nil {
	//}
	ctx.Response.SetBodyString("Ok")
}


func prepareSlackPayload(){

}
