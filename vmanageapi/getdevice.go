package vmanageapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Device struct {
	UUID  string
	Token string
}

func (app *VManageAPI) GetDevices() []*Device {
	result := make([]*Device, 0, 0)
	client := app.Client
	devicesURL := app.BaseURL + "/dataservice/system/device/vedges"
	response, err := client.Get(devicesURL)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(response.Body)

	var data1 map[string]interface{}
	jsonErr := json.Unmarshal(body, &data1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	devices := data1["data"].([]interface{})

	for value := range devices {
		// Assigning unknown
		current_device := devices[value].(map[string]interface{})
		if (current_device["deviceModel"] == "vedge-C8000V") && (current_device["vedgeCertificateState"] == "tokengenerated") {
			c8k := &Device{
				UUID:  current_device["chasisNumber"].(string),
				Token: current_device["serialNumber"].(string),
			}
			result = append(result, c8k)
		}
	}
	return result

}

func (app *VManageAPI) GetvBond() (string, string) {

	client := app.Client
	devicesURL := app.BaseURL + "/dataservice/settings/configuration/device"
	response, err := client.Get(devicesURL)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(response.Body)

	var data1 map[string]interface{}
	jsonErr := json.Unmarshal(body, &data1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	devices := data1["data"].([]interface{})

	for value := range devices {
		current_device := devices[value].(map[string]interface{})

		if (current_device["domainIp"] != nil) && (current_device["port"] != nil) {
			return current_device["domainIp"].(string), current_device["port"].(string)
		}
	}
	return "", ""
}

func (app *VManageAPI) GetOrgnazation() string {

	client := app.Client
	devicesURL := app.BaseURL + "/dataservice/settings/configuration/organization"
	response, err := client.Get(devicesURL)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(response.Body)

	var data1 map[string]interface{}
	jsonErr := json.Unmarshal(body, &data1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	devices := data1["data"].([]interface{})

	for value := range devices {
		current_device := devices[value].(map[string]interface{})

		if current_device["org"] != nil {
			return current_device["org"].(string)
		}
	}
	return ""
}

/* ADD Vbond

POST /dataservice/system/device

{"deviceIP":"1.2.3.4","username":"admin","password":"adin","personality":"vbond","generateCSR":true}
*/
