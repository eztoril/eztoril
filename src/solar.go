package main

import (
    //"bufio"
	"encoding/json"
    "fmt"
	"io/ioutil"
    "net/http"
	"log"
	"os"
	"time"
)

const (
	invAddr = "http://192.168.1.34"
	apiVersionURL = invAddr + "/solar_api/GetAPIVersion.cgi"
	invRTData = invAddr + "/solar_api/v1/GetInverterRealtimeData.cgi"
)

type cumInvRTDataReqType struct {
	url string
	params [6]string
	body cumInvRTDataReqResp
}

type commonInvRTDataReqType struct {
	url string
	params [6]string
	body cumInvRTDataReqResp
}

type cumInvRTDataReqResp struct {
	Body struct {
		Data struct {
			DayEnergy struct {
				Unit   string `json:"Unit"`
				Values struct {
					Num1 int `json:"1"`
				} `json:"Values"`
			} `json:"DAY_ENERGY"`
			Pac struct {
				Unit   string `json:"Unit"`
				Values struct {
					Num1 int `json:"1"`
				} `json:"Values"`
			} `json:"PAC"`
			TotalEnergy struct {
				Unit   string `json:"Unit"`
				Values struct {
					Num1 int `json:"1"`
				} `json:"Values"`
			} `json:"TOTAL_ENERGY"`
			YearEnergy struct {
				Unit   string `json:"Unit"`
				Values struct {
					Num1 int `json:"1"`
				} `json:"Values"`
			} `json:"YEAR_ENERGY"`
		} `json:"Data"`
	} `json:"Body"`
	Head struct {
		RequestArguments struct {
			DeviceClass string `json:"DeviceClass"`
			Scope       string `json:"Scope"`
		} `json:"RequestArguments"`
		Status struct {
			Code        int    `json:"Code"`
			Reason      string `json:"Reason"`
			UserMessage string `json:"UserMessage"`
		} `json:"Status"`
		Timestamp time.Time `json:"Timestamp"`
	} `json:"Head"`
}


// Defining an interface
type DataTypes interface {
	createHttpRequest() (string, *http.Response)
	parseData(string) cumInvRTDataReqResp
}

func (d cumInvRTDataReqType) parseData(bodyString string) cumInvRTDataReqResp {
	fmt.Println("cumInvRTDataReqType:parseData called")
	var body cumInvRTDataReqResp
	json.Unmarshal([]byte(bodyString), &body)
	return body
}

func (d cumInvRTDataReqType) createHttpRequest() (string, *http.Response) {
	fmt.Println("cumInvRTDataReqType:createHttpRequest called")
	req, err := http.NewRequest("GET", d.url, nil)
    if err != nil {
        log.Print(err)
        os.Exit(1)
    }
	
	// Add request info to the GET request
    q := req.URL.Query()
    q.Add(d.params[0], d.params[1])
	q.Add(d.params[2], d.params[3])
	q.Add(d.params[4], d.params[5])

	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
    if err != nil {
        panic(err)
    }

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	
    fmt.Println("Response status:", resp.Status)

	fmt.Println("cumInvRTDataReqType:createHttpRequest End")
	return bodyString, resp
}

//################################################################

func main() {
	var emptyBody cumInvRTDataReqResp
	invRTDataReq := cumInvRTDataReqType{invRTData, [6]string{
		"Scope", "System", "DeviceId", "1", "DataCollection", "CumulationInverterData1"},
		emptyBody}

	bodyString, resp := invRTDataReq.createHttpRequest()

	invRTDataReq.body = invRTDataReq.parseData(bodyString)
	fmt.Printf("%+v\n", invRTDataReq.body)

	//fmt.Println(bodyString)

	defer resp.Body.Close()
}

//################################################################
