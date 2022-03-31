package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/html"
)

const (
	fetchTime     = "23:50"
	addr          = "http://192.168.1.34"
	apiVersionURL = addr + "/solar_api/GetAPIVersion.cgi"
	invRTData     = addr + "/solar_api/v1/GetInverterRealtimeData.cgi"
	meterRTData   = addr + "/solar_api/v1/GetMeterRealtimeData.cgi"
	// Spot price
	spotAddr = "https://www.elbruk.se/timpriser-se3-stockholm"
)

type cumInvRTDataReqType struct {
	url    string
	params [6]string
	body   cumInvRTDataReqResp
}

type meterRTDataReqType struct {
	url    string
	params [4]string
	body   meterRTDataReqResp
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

type meterRTDataReqResp struct {
	Body struct {
		Data struct {
			Num1 struct {
				CurrentACPhase1 float64 `json:"Current_AC_Phase_1"`
				CurrentACPhase2 float64 `json:"Current_AC_Phase_2"`
				CurrentACPhase3 float64 `json:"Current_AC_Phase_3"`
				CurrentACSum    float64 `json:"Current_AC_Sum"`
				Details         struct {
					Manufacturer string `json:"Manufacturer"`
					Model        string `json:"Model"`
					Serial       string `json:"Serial"`
				} `json:"Details"`
				Enable                      int     `json:"Enable"`
				EnergyRealWACMinusAbsolute  int     `json:"EnergyReal_WAC_Minus_Absolute"`
				EnergyRealWACPhase1Consumed int     `json:"EnergyReal_WAC_Phase_1_Consumed"`
				EnergyRealWACPhase1Produced int     `json:"EnergyReal_WAC_Phase_1_Produced"`
				EnergyRealWACPhase2Consumed int     `json:"EnergyReal_WAC_Phase_2_Consumed"`
				EnergyRealWACPhase2Produced int     `json:"EnergyReal_WAC_Phase_2_Produced"`
				EnergyRealWACPhase3Consumed int     `json:"EnergyReal_WAC_Phase_3_Consumed"`
				EnergyRealWACPhase3Produced int     `json:"EnergyReal_WAC_Phase_3_Produced"`
				EnergyRealWACPlusAbsolute   int     `json:"EnergyReal_WAC_Plus_Absolute"`
				EnergyRealWACSumConsumed    int     `json:"EnergyReal_WAC_Sum_Consumed"`
				EnergyRealWACSumProduced    int     `json:"EnergyReal_WAC_Sum_Produced"`
				FrequencyPhaseAverage       float64 `json:"Frequency_Phase_Average"`
				MeterLocationCurrent        int     `json:"Meter_Location_Current"`
				PowerApparentSPhase1        float64 `json:"PowerApparent_S_Phase_1"`
				PowerApparentSPhase2        float64 `json:"PowerApparent_S_Phase_2"`
				PowerApparentSPhase3        float64 `json:"PowerApparent_S_Phase_3"`
				PowerApparentSSum           float64 `json:"PowerApparent_S_Sum"`
				PowerFactorPhase1           int     `json:"PowerFactor_Phase_1"`
				PowerFactorPhase2           int     `json:"PowerFactor_Phase_2"`
				PowerFactorPhase3           int     `json:"PowerFactor_Phase_3"`
				PowerFactorSum              int     `json:"PowerFactor_Sum"`
				PowerReactiveQPhase1        int     `json:"PowerReactive_Q_Phase_1"`
				PowerReactiveQPhase2        int     `json:"PowerReactive_Q_Phase_2"`
				PowerReactiveQPhase3        int     `json:"PowerReactive_Q_Phase_3"`
				PowerReactiveQSum           int     `json:"PowerReactive_Q_Sum"`
				PowerRealPPhase1            float64 `json:"PowerReal_P_Phase_1"`
				PowerRealPPhase2            float64 `json:"PowerReal_P_Phase_2"`
				PowerRealPPhase3            float64 `json:"PowerReal_P_Phase_3"`
				PowerRealPSum               float64 `json:"PowerReal_P_Sum"`
				TimeStamp                   int     `json:"TimeStamp"`
				Visible                     int     `json:"Visible"`
				VoltageACPhaseToPhase12     float64 `json:"Voltage_AC_PhaseToPhase_12"`
				VoltageACPhaseToPhase23     float64 `json:"Voltage_AC_PhaseToPhase_23"`
				VoltageACPhaseToPhase31     float64 `json:"Voltage_AC_PhaseToPhase_31"`
				VoltageACPhase1             float64 `json:"Voltage_AC_Phase_1"`
				VoltageACPhase2             float64 `json:"Voltage_AC_Phase_2"`
				VoltageACPhase3             float64 `json:"Voltage_AC_Phase_3"`
				VoltageACPhaseAverage       float64 `json:"Voltage_AC_Phase_Average"`
			} `json:"1"`
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
	var body cumInvRTDataReqResp
	json.Unmarshal([]byte(bodyString), &body)
	return body
}

func (d cumInvRTDataReqType) createHttpRequest() (string, *http.Response) {
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

	//fmt.Println("Fronius HTTP GET Response status:", resp.Status)
	return bodyString, resp
}

func (d meterRTDataReqType) parseData(bodyString string) meterRTDataReqResp {
	var body meterRTDataReqResp
	json.Unmarshal([]byte(bodyString), &body)
	return body
}

func (d meterRTDataReqType) createHttpRequest() (string, *http.Response) {
	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// Add request info to the GET request
	q := req.URL.Query()
	q.Add(d.params[0], d.params[1])
	q.Add(d.params[2], d.params[3])
	//	q.Add(d.params[4], d.params[5])

	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	//fmt.Println("Fronius HTTP GET Response status:", resp.Status)
	return bodyString, resp
}

// #################################################

func (d cumInvRTDataReqType) writeDB(body cumInvRTDataReqResp) int {
	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "solar:solar@tcp(127.0.0.1:3306)/solardb")
	defer db.Close()

	if err != nil {
		fmt.Println("sql.Open failure: ", err)
		return -1
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("db.Ping failure: ", err)
		return -1
	}
	currentTime := time.Now()
	fmt.Printf("Total power production at %s was: %d\n",
		currentTime.Format("2006-01-02 3:4"), body.Body.Data.DayEnergy.Values.Num1)
	query := fmt.Sprintf("INSERT INTO Production(Day,Power) values(curdate(), %d);",
		body.Body.Data.DayEnergy.Values.Num1)

	_, err = db.Query(query)

	if err != nil {
		fmt.Println("db.Query failure: ", err)
		return -1
	}
	return 0
}

func (d meterRTDataReqType) writeDB(body meterRTDataReqResp) int {
	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "solar:solar@tcp(127.0.0.1:3306)/solardb")
	defer db.Close()

	if err != nil {
		fmt.Println("sql.Open failure: ", err)
		return -1
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("db.Ping failure: ", err)
		return -1
	}
	currentTime := time.Now()
	fmt.Println("time:", currentTime)
	fmt.Printf("SumConsumed at %s was: %s\n",
		currentTime.Format("2006-01-02 3:4"), body.Body.Data.Num1.Details.Manufacturer)
	//query := fmt.Sprintf("INSERT INTO Production(Day,Power) values(curdate(), %d);",
	//	body.Body.Data.Num1.EnergyRealWACSumConsumed)

	//_, err = db.Query(query)

	//if err != nil {
	//	fmt.Println("db.Query failure: ", err)
	//	return -1
	//}
	return 0
}

//################################################################

func storePower() {
	var emptyBody cumInvRTDataReqResp
	invRTDataReq := cumInvRTDataReqType{invRTData, [6]string{
		"Scope", "System", "DeviceId", "1", "DataCollection", "CumulationInverterData1"},
		emptyBody}

	bodyString, resp := invRTDataReq.createHttpRequest()
	defer resp.Body.Close()

	invRTDataReq.body = invRTDataReq.parseData(bodyString)
	//fmt.Printf("%+v\n", invRTDataReq.body)

	_ = invRTDataReq.writeDB(invRTDataReq.body)
}

func storeMeter() {
	var emptyBody meterRTDataReqResp
	meterRTDataReq := meterRTDataReqType{meterRTData, [4]string{
		"Scope", "System", "DeviceId", "0"},
		emptyBody}

	fmt.Println("storeMeter: Enter")
	bodyString, resp := meterRTDataReq.createHttpRequest()
	defer resp.Body.Close()

	meterRTDataReq.body = meterRTDataReq.parseData(bodyString)
	fmt.Printf("%+v\n", meterRTDataReq.body)

	_ = meterRTDataReq.writeDB(meterRTDataReq.body)
}

// ##############################################################################

func getHtmlPage(webPage string) (string, error) {

	//	fmt.Printf("getHtmlPage %s", webPage)

	resp, err := http.Get(webPage)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseMonthlySpot(input string) float64 {
	// Find correct data portion
	var data string
	ix := strings.Index(input, "label: '2022'")
	if ix != -1 {
		s1 := input[ix:]
		ix2 := strings.Index(s1, "data: [")
		s2 := s1[ix2:]
		//fmt.Printf("ix: %d : %s", ix, s2)
		ix3 := strings.Index(s2, " ]")
		data = s2[:ix3+1]
		//fmt.Printf("ix3: %d data: %s", ix3, data)

		re := regexp.MustCompile("[ ]?([0-9]*[.])?[0-9]+ ")
		found := re.FindAllString(data, -1)
		s := strings.TrimSpace(found[len(found)-1])
		f, _ := strconv.ParseFloat(s, 32)
		//fmt.Printf("%.2f\n", f)

		return f

		//	if ix != -1 {
		//		fmt.Printf("ix2: %d", ix2)
		//	}
	}
	return -1
}

func parseHtml(text string) float64 {

	tkn := html.NewTokenizer(strings.NewReader(text))
	var isTd bool

	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return -1
		case tt == html.StartTagToken:
			t := tkn.Token()
			isTd = t.Data == "script"
		case tt == html.TextToken:
			t := tkn.Token()
			if isTd {
				spot := parseMonthlySpot(t.Data)
				if spot != -1 {
					return spot
				}
			}
			isTd = false
		}
	}
}

func getSpotPrice(webPage string) {
	data, err := getHtmlPage(webPage)

	if err != nil {
		log.Fatal(err)
	}
	spot := parseHtml(data)
	fmt.Printf("%.2f\n", spot)
}

func storeData() {
	storePower()
	//storeMeter()
}

func main() {
	fmt.Println("Starting solar power production storage. Daily fetching time:", fetchTime)
	storePower()
	getSpotPrice(spotAddr)
	//	time.Sleep(20 * time.Second)
	//	storeMeter()
	//	s := gocron.NewScheduler()
	//	s.Every(1).Day().At(fetchTime).Do(storeData)
	//	<-s.Start()
}

// user: solar pw: solar ################################################################
//
// martin@htpc:~/repo/solar$ sudo mysql -u root -p              (pw: htpc)
// martin@htpc:~/repo/solar/src$ sudo mysql -u solar -p         (pw: solar)

// MariaDB [(none)]> CREATE DATABASE solardb
// MariaDB [(none)]> use solardb;
// MariaDB [solardb]> SELECT user FROM mysql.user;
// MariaDB [solardb]> CREATE USER 'solar'@'localhost' IDENTIFIED BY 'password';
// MariaDB [solardb]> GRANT ALL PRIVILEGES ON * . * TO 'solar'@'localhost';

// create table DemoTable
// (
//   StudentId int NOT NULL AUTO_INCREMENT PRIMARY KEY,
//   StudentName varchar(20),
//   StudentAdmissionDate DATE
// );
// insert into Production(Day,Power) values(now(), 14);
// insert into Production(Day,Power) values('2021-05-27', '105');

// MariaDB [solardb]> CREATE TABLE Production ( Day DATE PRIMARY KEY, Power int );

// select sum(Power) from Production where Production.Day Between '2021-04-27' and '2022-05-27';
