package ocs

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/nextcloud/nextcloudgo"
)

// Request is a wrapper around net.http.Request which also takes care of authentication
type Request struct {
	nc nextcloudgo.NextcloudGo
}

// New returns a new OCS request handler when given the NextcloudGo
func New(nc nextcloudgo.NextcloudGo) Request {
	return Request{nc: nc}
}

// Request performs an (authenticated) request and returns the data
func (ocs *Request) Request(method, url string, auth bool) (map[string]interface{}, int, error) {
	return ocs.RequestWithBody(method, url, nil, auth)
}

// RequestWithBody performs an (authenticated) request and returns the data
func (ocs *Request) RequestWithBody(method, url string, body io.Reader, auth bool) (map[string]interface{}, int, error) {
	response, err := ocs.nc.Request(method, url, body, auth)
	if response == nil {
		return nil, 400, err
	}

	if err != nil {
		return nil, response.StatusCode, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	var mixed interface{}
	json.Unmarshal(contents, &mixed)

	data, ok := mixed.(map[string]interface{})
	if !ok {
		return nil, response.StatusCode, errors.New("Invalid JSON response")
	}
	return data, response.StatusCode, nil
}

// ValidateStatusCode checks whether the OCS status code matches the accepted value
func ValidateStatusCode(data map[string]interface{}, accepted int) bool {
	status, err := GetInt(data, []string{"ocs", "meta", "statuscode"})
	if err != nil {
		return false
	}
	return accepted == status
}

// GetInt returns a single int from a given subtree in the OCS response
func GetInt(data map[string]interface{}, keys []string) (int, error) {
	var ok bool
	var element float64

	for k, v := range keys {
		if k == len(keys)-1 {
			element, ok = data[v].(float64)
			if !ok {
				return 0, errors.New("Error while trying to get OCS response subtree")
			}
			return int(element), nil
		}

		data, ok = data[v].(map[string]interface{})
		if !ok {
			return 0, errors.New("Error while trying to get OCS response subtree")
		}
	}

	return 0, errors.New("Error while trying to get OCS response subtree")
}

// GetStringList returns a string array from a given subtree in the OCS response
func GetStringList(data map[string]interface{}, keys []string) ([]string, error) {
	var elements []interface{}
	var ok bool

	for k, v := range keys {
		if k == len(keys)-1 {
			elements, ok = data[v].([]interface{})
			if !ok {
				return []string{}, errors.New("Error while trying to get OCS response subtree")
			}
		} else {
			data, ok = data[v].(map[string]interface{})
			if !ok {
				return []string{}, errors.New("Error while trying to get OCS response subtree")
			}
		}
	}

	var list []string

	for _, element := range elements {
		list = append(list, element.(string))
	}

	return list, nil
}
