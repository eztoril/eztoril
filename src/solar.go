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
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonlvhit/gocron"
)

const (
	fetchTime = "23:50"
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

func main() {
	fmt.Println("Main: Starting solar power production store. Daily fetching time:", fetchTime)
    s := gocron.NewScheduler()
    s.Every(1).Day().At(fetchTime).Do(storePower)
	<- s.Start()
}

// user: solar pw: solar ################################################################
//
// martin@htpc:~/repo/solar$ sudo mysql -u root -p
//
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
