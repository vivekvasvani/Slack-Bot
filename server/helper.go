package server

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type STFResponse struct {
	Success bool `json:"success"`
	Devices []struct {
		Abi          string `json:"abi"`
		AirplaneMode bool   `json:"airplaneMode"`
		Battery      struct {
			Health  string  `json:"health"`
			Level   int     `json:"level"`
			Scale   int     `json:"scale"`
			Source  string  `json:"source"`
			Status  string  `json:"status"`
			Temp    float64 `json:"temp"`
			Voltage float64 `json:"voltage"`
		} `json:"battery"`
		Browser struct {
			Apps []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Selected  bool   `json:"selected"`
				System    bool   `json:"system"`
				Type      string `json:"type"`
				Developer string `json:"developer"`
			} `json:"apps"`
			Selected bool `json:"selected"`
		} `json:"browser"`
		Channel   string    `json:"channel"`
		CreatedAt time.Time `json:"createdAt"`
		Display   struct {
			Density  float64 `json:"density"`
			Fps      float64 `json:"fps"`
			Height   int     `json:"height"`
			ID       int     `json:"id"`
			Rotation int     `json:"rotation"`
			Secure   bool    `json:"secure"`
			Size     float64 `json:"size"`
			URL      string  `json:"url"`
			Width    int     `json:"width"`
			Xdpi     float64 `json:"xdpi"`
			Ydpi     float64 `json:"ydpi"`
		} `json:"display"`
		Manufacturer string `json:"manufacturer"`
		Model        string `json:"model"`
		Network      struct {
			Connected bool   `json:"connected"`
			Failover  bool   `json:"failover"`
			Roaming   bool   `json:"roaming"`
			Subtype   string `json:"subtype"`
			Type      string `json:"type"`
		} `json:"network"`
		Notes    string      `json:"notes"`
		Operator interface{} `json:"operator"`
		Owner    struct {
			Email string `json:"email"`
			Group string `json:"group"`
			Name  string `json:"name"`
		} `json:"owner"`
		Phone struct {
			Iccid       interface{} `json:"iccid"`
			Imei        string      `json:"imei"`
			Imsi        interface{} `json:"imsi"`
			Network     string      `json:"network"`
			PhoneNumber interface{} `json:"phoneNumber"`
		} `json:"phone"`
		Platform          string    `json:"platform"`
		PresenceChangedAt time.Time `json:"presenceChangedAt"`
		Present           bool      `json:"present"`
		Product           string    `json:"product"`
		Provider          struct {
			Channel string `json:"channel"`
			Name    string `json:"name"`
		} `json:"provider"`
		Ready            bool          `json:"ready"`
		RemoteConnect    bool          `json:"remoteConnect"`
		RemoteConnectURL interface{}   `json:"remoteConnectUrl"`
		ReverseForwards  []interface{} `json:"reverseForwards"`
		Sdk              string        `json:"sdk"`
		Serial           string        `json:"serial"`
		Status           int           `json:"status"`
		StatusChangedAt  time.Time     `json:"statusChangedAt"`
		Usage            interface{}   `json:"usage"`
		UsageChangedAt   time.Time     `json:"usageChangedAt"`
		Version          string        `json:"version"`
		Using            bool          `json:"using"`
	} `json:"devices"`
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
		if strings.ContainsAny(textData, "${"+strconv.Itoa(i)) {
			textData = strings.Replace(textData, "${"+strconv.Itoa(i)+"}", value, -1)
		}
	}

	//textData = strings.Replace(textData, "${TA}", strconv.Itoa(len(strings.Split(sessionMap[0], "\n"))), -1)
	//textData = strings.Replace(textData, "${TB}", strconv.Itoa(len(strings.Split(sessionMap[1], "\n"))), -1)
	//textData = strings.Replace(textData, "${TD}", strconv.Itoa(len(strings.Split(sessionMap[2], "\n"))), -1)
	return textData
}
