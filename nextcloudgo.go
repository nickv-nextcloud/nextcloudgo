package nextcloudgo

import (
	"encoding/json"

	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type NextcloudGo struct {
	serverUrl string
	user      string
	password  string
}

func (sdk *NextcloudGo) Connect(serverUrl string) bool {
	sdk.serverUrl = serverUrl
	if !sdk.IsConnected() {
		return false
	}
	status := sdk.Status()
	return status["version"] != ""
}

func (sdk *NextcloudGo) Disconnect() bool {
	sdk.serverUrl = ""
	return !sdk.IsConnected()
}

func (sdk *NextcloudGo) IsConnected() bool {
	return sdk.serverUrl != ""
}

func (sdk *NextcloudGo) Login(user string, password string) bool {
	if !sdk.IsConnected() {
		return false
	}

	sdk.user = user
	sdk.password = password
	return sdk.IsLoggedIn() // checkCapabilities
}

func (sdk *NextcloudGo) Logout() bool {
	if !sdk.IsConnected() {
		return false
	}

	sdk.user = ""
	sdk.password = ""
	return !sdk.IsLoggedIn()
}

func (sdk *NextcloudGo) IsLoggedIn() bool {
	return sdk.IsConnected() && sdk.user != "" && sdk.password != ""
}

func (sdk *NextcloudGo) Capabilities() interface{} {
	//capabilities := make(map[string]string)
	var capabilities interface{}
	if !sdk.IsLoggedIn() {
		return capabilities
	}

	client := &http.Client{}

	/* Authenticate */
	req, err := http.NewRequest(http.MethodGet, sdk.serverUrl+"ocs/v1.php/cloud/capabilities?format=json", nil)
	req.SetBasicAuth(sdk.user, sdk.password)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return capabilities
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
			return capabilities
		}

		json.Unmarshal(contents, &capabilities)
		return capabilities
	}
}

func (sdk *NextcloudGo) Status() map[string]string {
	status := make(map[string]string)
	if !sdk.IsConnected() {
		return status
	}

	client := &http.Client{}
	response, err := client.Get(sdk.serverUrl + "status.php")
	if err != nil {
		log.Fatal(err)
		return status
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
			return status
		}

		json.Unmarshal(contents, &status)
		return status
	}
}

func (sdk *NextcloudGo) GetServerUrl() string {
	return sdk.serverUrl
}

func (sdk *NextcloudGo) GetAuthUser() string {
	return sdk.user
}

func (sdk *NextcloudGo) GetAuthPassword() string {
	return sdk.password
}
