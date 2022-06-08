package invdata

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	meterRtDataUrl = "/solar_api/v1/GetMeterRealtimeData.cgi"
)

type meterRtDataReq struct {
	url        string
	params     [6]string
	bodyString string
}

type meterRtDataResp struct {
	Body struct {
		Data struct {
			Num0 struct {
				CurrentACPhase1 float64 `json:"Current_AC_Phase_1"`
				CurrentACPhase2 float64 `json:"Current_AC_Phase_2"`
				CurrentACPhase3 float64 `json:"Current_AC_Phase_3"`
				Details         struct {
					Manufacturer string `json:"Manufacturer"`
					Model        string `json:"Model"`
					Serial       string `json:"Serial"`
				} `json:"Details"`
				Enable                         int     `json:"Enable"`
				EnergyReactiveVArACSumConsumed int     `json:"EnergyReactive_VArAC_Sum_Consumed"`
				EnergyReactiveVArACSumProduced int     `json:"EnergyReactive_VArAC_Sum_Produced"`
				EnergyRealWACMinusAbsolute     int     `json:"EnergyReal_WAC_Minus_Absolute"`
				EnergyRealWACPlusAbsolute      int     `json:"EnergyReal_WAC_Plus_Absolute"`
				EnergyRealWACSumConsumed       int     `json:"EnergyReal_WAC_Sum_Consumed"`
				EnergyRealWACSumProduced       int     `json:"EnergyReal_WAC_Sum_Produced"`
				FrequencyPhaseAverage          int     `json:"Frequency_Phase_Average"`
				MeterLocationCurrent           int     `json:"Meter_Location_Current"`
				PowerApparentSPhase1           float64 `json:"PowerApparent_S_Phase_1"`
				PowerApparentSPhase2           float64 `json:"PowerApparent_S_Phase_2"`
				PowerApparentSPhase3           float64 `json:"PowerApparent_S_Phase_3"`
				PowerApparentSSum              int     `json:"PowerApparent_S_Sum"`
				PowerFactorPhase1              int     `json:"PowerFactor_Phase_1"`
				PowerFactorPhase2              float64 `json:"PowerFactor_Phase_2"`
				PowerFactorPhase3              int     `json:"PowerFactor_Phase_3"`
				PowerFactorSum                 float64 `json:"PowerFactor_Sum"`
				PowerReactiveQPhase1           float64 `json:"PowerReactive_Q_Phase_1"`
				PowerReactiveQPhase2           float64 `json:"PowerReactive_Q_Phase_2"`
				PowerReactiveQPhase3           int     `json:"PowerReactive_Q_Phase_3"`
				PowerReactiveQSum              float64 `json:"PowerReactive_Q_Sum"`
				PowerRealPPhase1               int     `json:"PowerReal_P_Phase_1"`
				PowerRealPPhase2               float64 `json:"PowerReal_P_Phase_2"`
				PowerRealPPhase3               float64 `json:"PowerReal_P_Phase_3"`
				PowerRealPSum                  float64 `json:"PowerReal_P_Sum"`
				TimeStamp                      int     `json:"TimeStamp"`
				Visible                        int     `json:"Visible"`
				VoltageACPhaseToPhase12        float64 `json:"Voltage_AC_PhaseToPhase_12"`
				VoltageACPhaseToPhase23        float64 `json:"Voltage_AC_PhaseToPhase_23"`
				VoltageACPhaseToPhase31        float64 `json:"Voltage_AC_PhaseToPhase_31"`
				VoltageACPhase1                float64 `json:"Voltage_AC_Phase_1"`
				VoltageACPhase2                float64 `json:"Voltage_AC_Phase_2"`
				VoltageACPhase3                float64 `json:"Voltage_AC_Phase_3"`
			} `json:"0"`
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
			Num2 struct {
				CurrentACPhase1 float64 `json:"Current_AC_Phase_1"`
				CurrentACSum    float64 `json:"Current_AC_Sum"`
				Details         struct {
					Manufacturer string `json:"Manufacturer"`
					Model        string `json:"Model"`
					Serial       string `json:"Serial"`
				} `json:"Details"`
				Enable                            int     `json:"Enable"`
				EnergyReactiveVArACPhase1Consumed int     `json:"EnergyReactive_VArAC_Phase_1_Consumed"`
				EnergyReactiveVArACPhase1Produced int     `json:"EnergyReactive_VArAC_Phase_1_Produced"`
				EnergyReactiveVArACSumConsumed    int     `json:"EnergyReactive_VArAC_Sum_Consumed"`
				EnergyReactiveVArACSumProduced    int     `json:"EnergyReactive_VArAC_Sum_Produced"`
				EnergyRealWACMinusAbsolute        int     `json:"EnergyReal_WAC_Minus_Absolute"`
				EnergyRealWACPhase1Consumed       int     `json:"EnergyReal_WAC_Phase_1_Consumed"`
				EnergyRealWACPhase1Produced       int     `json:"EnergyReal_WAC_Phase_1_Produced"`
				EnergyRealWACPlusAbsolute         int     `json:"EnergyReal_WAC_Plus_Absolute"`
				EnergyRealWACSumConsumed          int     `json:"EnergyReal_WAC_Sum_Consumed"`
				EnergyRealWACSumProduced          int     `json:"EnergyReal_WAC_Sum_Produced"`
				FrequencyPhaseAverage             int     `json:"Frequency_Phase_Average"`
				MeterLocationCurrent              int     `json:"Meter_Location_Current"`
				PowerApparentSPhase1              float64 `json:"PowerApparent_S_Phase_1"`
				PowerApparentSSum                 float64 `json:"PowerApparent_S_Sum"`
				PowerFactorPhase1                 float64 `json:"PowerFactor_Phase_1"`
				PowerFactorSum                    float64 `json:"PowerFactor_Sum"`
				PowerReactiveQPhase1              float64 `json:"PowerReactive_Q_Phase_1"`
				PowerReactiveQSum                 float64 `json:"PowerReactive_Q_Sum"`
				PowerRealPPhase1                  float64 `json:"PowerReal_P_Phase_1"`
				PowerRealPSum                     float64 `json:"PowerReal_P_Sum"`
				TimeStamp                         int     `json:"TimeStamp"`
				Visible                           int     `json:"Visible"`
				VoltageACPhase1                   float64 `json:"Voltage_AC_Phase_1"`
			} `json:"2"`
			Num3 struct {
				Details struct {
					Manufacturer string `json:"Manufacturer"`
					Model        string `json:"Model"`
					Serial       string `json:"Serial"`
				} `json:"Details"`
				Enable                     int `json:"Enable"`
				MeterLocationCurrent       int `json:"Meter_Location_Current"`
				TimeStamp                  int `json:"TimeStamp"`
				EnergyRealWACMinusRelative int `json:"EnergyReal_WAC_Minus_Relative"`
				EnergyRealWACPlusRelative  int `json:"EnergyReal_WAC_Plus_Relative"`
				PowerRealPSum              int `json:"PowerReal_P_Sum"`
				Visible                    int `json:"Visible"`
			} `json:"3"`
			Num4 struct {
				CurrentACPhase1 int `json:"Current_AC_Phase_1"`
				CurrentACPhase2 int `json:"Current_AC_Phase_2"`
				CurrentACPhase3 int `json:"Current_AC_Phase_3"`
				CurrentACSum    int `json:"Current_AC_Sum"`
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
				PowerApparentSPhase1        int     `json:"PowerApparent_S_Phase_1"`
				PowerApparentSPhase2        int     `json:"PowerApparent_S_Phase_2"`
				PowerApparentSPhase3        int     `json:"PowerApparent_S_Phase_3"`
				PowerApparentSSum           int     `json:"PowerApparent_S_Sum"`
				PowerFactorPhase1           int     `json:"PowerFactor_Phase_1"`
				PowerFactorPhase2           int     `json:"PowerFactor_Phase_2"`
				PowerFactorPhase3           int     `json:"PowerFactor_Phase_3"`
				PowerFactorSum              int     `json:"PowerFactor_Sum"`
				PowerReactiveQPhase1        int     `json:"PowerReactive_Q_Phase_1"`
				PowerReactiveQPhase2        int     `json:"PowerReactive_Q_Phase_2"`
				PowerReactiveQPhase3        int     `json:"PowerReactive_Q_Phase_3"`
				PowerReactiveQSum           int     `json:"PowerReactive_Q_Sum"`
				PowerRealPPhase1            int     `json:"PowerReal_P_Phase_1"`
				PowerRealPPhase2            int     `json:"PowerReal_P_Phase_2"`
				PowerRealPPhase3            int     `json:"PowerReal_P_Phase_3"`
				PowerRealPSum               int     `json:"PowerReal_P_Sum"`
				TimeStamp                   int     `json:"TimeStamp"`
				Visible                     int     `json:"Visible"`
				VoltageACPhaseToPhase12     float64 `json:"Voltage_AC_PhaseToPhase_12"`
				VoltageACPhaseToPhase23     int     `json:"Voltage_AC_PhaseToPhase_23"`
				VoltageACPhaseToPhase31     float64 `json:"Voltage_AC_PhaseToPhase_31"`
				VoltageACPhase1             float64 `json:"Voltage_AC_Phase_1"`
				VoltageACPhase2             float64 `json:"Voltage_AC_Phase_2"`
				VoltageACPhase3             float64 `json:"Voltage_AC_Phase_3"`
				VoltageACPhaseAverage       float64 `json:"Voltage_AC_Phase_Average"`
			} `json:"4"`
			Num5 struct {
				CurrentACPhase1 int `json:"Current_AC_Phase_1"`
				CurrentACPhase2 int `json:"Current_AC_Phase_2"`
				CurrentACPhase3 int `json:"Current_AC_Phase_3"`
				Details         struct {
					Manufacturer string `json:"Manufacturer"`
					Model        string `json:"Model"`
					Serial       string `json:"Serial"`
				} `json:"Details"`
				Enable                         int     `json:"Enable"`
				EnergyReactiveVArACSumConsumed int     `json:"EnergyReactive_VArAC_Sum_Consumed"`
				EnergyReactiveVArACSumProduced int     `json:"EnergyReactive_VArAC_Sum_Produced"`
				EnergyRealWACMinusAbsolute     int     `json:"EnergyReal_WAC_Minus_Absolute"`
				EnergyRealWACPlusAbsolute      int     `json:"EnergyReal_WAC_Plus_Absolute"`
				EnergyRealWACSumConsumed       int     `json:"EnergyReal_WAC_Sum_Consumed"`
				EnergyRealWACSumProduced       int     `json:"EnergyReal_WAC_Sum_Produced"`
				FrequencyPhaseAverage          float64 `json:"Frequency_Phase_Average"`
				MeterLocationCurrent           int     `json:"Meter_Location_Current"`
				PowerApparentSPhase1           int     `json:"PowerApparent_S_Phase_1"`
				PowerApparentSPhase2           int     `json:"PowerApparent_S_Phase_2"`
				PowerApparentSPhase3           int     `json:"PowerApparent_S_Phase_3"`
				PowerApparentSSum              int     `json:"PowerApparent_S_Sum"`
				PowerFactorPhase1              float64 `json:"PowerFactor_Phase_1"`
				PowerFactorPhase2              float64 `json:"PowerFactor_Phase_2"`
				PowerFactorPhase3              float64 `json:"PowerFactor_Phase_3"`
				PowerFactorSum                 float64 `json:"PowerFactor_Sum"`
				PowerReactiveQPhase1           int     `json:"PowerReactive_Q_Phase_1"`
				PowerReactiveQPhase2           int     `json:"PowerReactive_Q_Phase_2"`
				PowerReactiveQPhase3           int     `json:"PowerReactive_Q_Phase_3"`
				PowerReactiveQSum              int     `json:"PowerReactive_Q_Sum"`
				PowerRealPPhase1               int     `json:"PowerReal_P_Phase_1"`
				PowerRealPPhase2               int     `json:"PowerReal_P_Phase_2"`
				PowerRealPPhase3               int     `json:"PowerReal_P_Phase_3"`
				PowerRealPSum                  int     `json:"PowerReal_P_Sum"`
				TimeStamp                      int     `json:"TimeStamp"`
				Visible                        int     `json:"Visible"`
				VoltageACPhaseToPhase12        float64 `json:"Voltage_AC_PhaseToPhase_12"`
				VoltageACPhaseToPhase23        float64 `json:"Voltage_AC_PhaseToPhase_23"`
				VoltageACPhaseToPhase31        float64 `json:"Voltage_AC_PhaseToPhase_31"`
				VoltageACPhase1                float64 `json:"Voltage_AC_Phase_1"`
				VoltageACPhase2                float64 `json:"Voltage_AC_Phase_2"`
				VoltageACPhase3                float64 `json:"Voltage_AC_Phase_3"`
			} `json:"5"`
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

func FetchMeterRtData(addr string) (int64, int64) {
	meterReq := meterRtDataReq{"", [6]string{
		"Scope", "System", "DeviceId", "1", "DataCollection", "CumulationInverterData1"}, "",
	}
	resp := meterReq.httpGet(addr)
	defer resp.Body.Close()

	bodyString := meterReq.parseJsonData()
	energyFromGrid, energyToGrid := parseMeterData(bodyString)
	return energyFromGrid, energyToGrid
}

func (d *meterRtDataReq) parseJsonData() string {
	var body meterRtDataResp
	json.Unmarshal([]byte(d.bodyString), &body)
	//fmt.Printf(d.bodyString)
	return d.bodyString
}

func parseMeterData(bodyString string) (energyFromGrid int64, energyToGrid int64) {
	scanner := bufio.NewScanner(strings.NewReader(bodyString))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "EnergyReal_WAC_Sum_Consumed") {
			dataIx := strings.Index(scanner.Text(), " : ")
			if dataIx != -1 {
				data := scanner.Text()[dataIx+3:]
				data = strings.TrimSuffix(data, ",")
				energyFromGrid, _ = strconv.ParseInt(data, 10, 0)
				//fmt.Printf("energyToGrid: %d\n", energyFromGrid)
			}
		}
		if strings.Contains(scanner.Text(), "EnergyReal_WAC_Sum_Produced") {
			dataIx := strings.Index(scanner.Text(), " : ")
			if dataIx != -1 {
				data := scanner.Text()[dataIx+3:]
				data = strings.TrimSuffix(data, ",")
				energyToGrid, _ = strconv.ParseInt(data, 10, 0)
				//fmt.Printf("energyToGrid: %d\n", energyToGrid)
			}
		}
	}
	return energyFromGrid, energyToGrid
}

func (d *meterRtDataReq) httpGet(invAddr string) *http.Response {
	d.url = invAddr + meterRtDataUrl
	d.params = [6]string{
		"Scope", "System", "DeviceId", "1", "DataCollection", "CumulationInverterData1",
	}
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
	d.bodyString = string(bodyBytes)

	//fmt.Printf(d.bodyString)
	//fmt.Println("Fronius HTTP GET Response status:", resp.Status)
	return resp
}
