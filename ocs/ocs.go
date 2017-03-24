package ocs

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/nextcloud/nextcloudgo"
)

type Request struct {
	sdk nextcloudgo.NextcloudGo
}

func New(sdk nextcloudgo.NextcloudGo) Request {
	return Request{sdk: sdk}
}

func (ocs *Request) NewRequest(method, url string, auth bool) (map[string]interface{}, error) {
	response, err := ocs.sdk.Request(method, url, auth)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var mixed interface{}
	json.Unmarshal(contents, &mixed)

	data, ok := mixed.(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid JSON response")
	}
	return data, nil
}

func ValidateStatusCode(data map[string]interface{}, accepted int) bool {

	ocs, ok := data["ocs"].(map[string]interface{})
	if !ok {
		return ok
	}

	meta := ocs["meta"].(map[string]interface{})
	if !ok {
		return ok
	}

	status := int(meta["statuscode"].(float64))
	return accepted == status
}
