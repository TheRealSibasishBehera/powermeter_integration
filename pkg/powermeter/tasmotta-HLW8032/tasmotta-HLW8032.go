package powermeter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EnergyStats struct {
	StatusSNS struct {
		Time   string `json:"Time"`
		Energy struct {
			TotalStartTime string  `json:"TotalStartTime"`
			Total          float64 `json:"Total"`
			Yesterday      float64 `json:"Yesterday"`
			Today          float64 `json:"Today"`
			Power          int     `json:"Power"`
			ApparentPower  int     `json:"ApparentPower"`
			ReactivePower  int     `json:"ReactivePower"`
			Factor         float64 `json:"Factor"`
			Voltage        int     `json:"Voltage"`
			Current        float64 `json:"Current"`
		} `json:"ENERGY"`
	} `json:"StatusSNS"`
}

type SwitchState struct {
	Power string `json:"POWER"`
}

func SwitchPowerMeter() []byte {
	resp, err := http.Get("http://192.168.43.51/cm?cmnd=POWER+TOGGLE")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	return body
}
func GetSwitchStateJSON() SwitchState {
	var result SwitchState
	body := SwitchPowerMeter()
	if err := json.Unmarshal(
		body,
		&result,
	); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	fmt.Println(PrettyPrint(result))
	return result
}

func GetEnergyStats() EnergyStats {
	body := GetEnergyStatBody()
	result := UnMarshalEnergyStatJSON(body)
	fmt.Println(PrettyPrint(result))
	return result
}

func UnMarshalEnergyStatJSON(body []byte) EnergyStats {
	var result EnergyStats
	if err := json.Unmarshal(
		body,
		&result,
	); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	return result
}

func GetEnergyStatBody() []byte {
	resp, err := http.Get("http://192.168.43.51/cm?cmnd=STATUS+8")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	return body
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(
		i,
		"",
		"\t",
	)
	return string(s)
}
