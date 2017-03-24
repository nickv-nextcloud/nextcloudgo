package nextcloudgo

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"

	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type NextcloudGo struct {
	serverUrl string
	certPath  string
	user      string
	password  string
}

func (sdk *NextcloudGo) SetCustomCertPath(certPath string) {
	sdk.certPath = certPath
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

	response, err := sdk.Request(http.MethodGet, "ocs/v1.php/cloud/capabilities?format=json", true)
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

	response, err := sdk.Request(http.MethodGet, "status.php", false)
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

func (sdk *NextcloudGo) Request(method, url string, auth bool) (*http.Response, error) {
	req, err := http.NewRequest(method, sdk.serverUrl+url, nil)
	if err != nil {
		return nil, err
	}

	if auth {
		if !sdk.IsLoggedIn() {
			return nil, errors.New("No user logged in")
		}
		req.SetBasicAuth(sdk.user, sdk.password)
	}
	req.Header.Add("OCS-APIRequest", "true")

	client := &http.Client{}
	if sdk.certPath != "" {
		tlsConfig := &tls.Config{}
		certs := x509.NewCertPool()

		pemData, err := ioutil.ReadFile(sdk.certPath)
		if err != nil {
			// do error
		}
		certs.AppendCertsFromPEM(pemData)
		tlsConfig.RootCAs = certs

		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		client = &http.Client{Transport: tr}
	}

	return client.Do(req)
}
