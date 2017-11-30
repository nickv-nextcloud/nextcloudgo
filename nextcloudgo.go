package nextcloudgo

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"

	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//
type NextcloudGo struct {
	ServerURL string
	CertPath  string
	User      string
	Password  string
}

// Status object for the server which mirrors the status.php content
type Status struct {
	// Installed reflects the installation state of the server
	Installed bool `json:"installed"`
	// Maintenance is true, when an update is in progress or the instance was put to maintenance manually
	Maintenance bool `json:"maintenance"`
	// Version is the numeric/internal version number of the server, e.g. 13.0.0.6
	Version string `json:"version"`
	// VersionString is the readable version string shown in the admin interface, e.g. 13.0.0 Beta1
	VersionString string `json:"versionstring"`
}

var (
	// ErrNotConnected is returned when the api has no ServerURL set.
	ErrNotConnected = errors.New("Not connected to any server")
	// ErrNoUserOrPassword is returned when the api has no user and/or password set.
	ErrNoUserOrPassword = errors.New("No user/password given")
)

func (nc *NextcloudGo) isConnected() bool {
	return nc.ServerURL != ""
}

func (nc *NextcloudGo) isLoggedIn() bool {
	return nc.isConnected() && nc.User != "" && nc.Password != ""
}

// Capabilities should be moved to ocs
// TODO(nickvergessen) Move out of the base
func (nc *NextcloudGo) Capabilities() interface{} {
	//capabilities := make(map[string]string)
	var capabilities interface{}
	if !nc.isLoggedIn() {
		return capabilities
	}

	response, err := nc.Request(http.MethodGet, "/ocs/v1.php/cloud/capabilities?format=json", nil, true)
	if err != nil {
		log.Fatal(err)
		return capabilities
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return capabilities
	}

	json.Unmarshal(contents, &capabilities)
	return capabilities
}

// Status returns the Status of the server
func (nc *NextcloudGo) Status() (Status, error) {
	status := Status{Maintenance: true}
	if !nc.isConnected() {
		return status, ErrNotConnected
	}

	response, err := nc.Request(http.MethodGet, "/status.php", nil, false)
	if err != nil {
		return status, err
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&status)

	return status, err
}

// Request performs a request to the given url and takes care of the authentication
// Content-Type and everything else. But in general you should not need to use this
// method yourself.
func (nc *NextcloudGo) Request(method, url string, body io.Reader, auth bool) (*http.Response, error) {
	req, err := http.NewRequest(method, nc.ServerURL+url, body)

	if err != nil {
		return nil, err
	}

	if auth {
		if !nc.isLoggedIn() {
			return nil, errors.New(nc.ServerURL) //ErrNoUserOrPassword
		}
		req.SetBasicAuth(nc.User, nc.Password)
	}
	req.Header.Add("OCS-APIRequest", "true")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	if nc.CertPath != "" {
		tlsConfig := &tls.Config{}
		certs := x509.NewCertPool()

		pemData, err := ioutil.ReadFile(nc.CertPath)
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
