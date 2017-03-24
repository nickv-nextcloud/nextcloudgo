package provisioning

import (
	"errors"
	"net/http"

	"github.com/nextcloud/nextcloudgo/ocs"
)

func (api *Provisioning) IsAppEnabled(appid string) (bool, error) {
	return api.isAppInArray(appid, endpoint+"/apps?filter=enabled&format=json")
}

func (api *Provisioning) IsAppAvailable(appid string) (bool, error) {
	return api.isAppInArray(appid, endpoint+"/apps?format=json")
}

func (api *Provisioning) isAppInArray(appid, url string) (bool, error) {
	content, err := api.ocs.NewRequest(http.MethodGet, url, true)
	if err != nil {
		return false, err
	}

	if !ocs.ValidateStatusCode(content, 200) {
		return false, errors.New("Status code was invalid")
	}

	ocs := content["ocs"].(map[string]interface{})
	data := ocs["data"].(map[string]interface{})
	apps := data["apps"].([]interface{})

	for _, app := range apps {
		if app == appid {
			return true, nil
		}
	}

	return false, nil
}

func (api *Provisioning) EnableApp(appid string) error {
	return api.changeAppState(appid, http.MethodPost)
}

func (api *Provisioning) DisableApp(appid string) error {
	return api.changeAppState(appid, http.MethodDelete)
}

func (api *Provisioning) changeAppState(appid, method string) error {
	url := endpoint + "/apps/" + appid + "?format=json"

	content, err := api.ocs.NewRequest(method, url, true)
	if err != nil {
		return err
	}

	if !ocs.ValidateStatusCode(content, 200) {
		if method == http.MethodPost {
			return errors.New("An error occured while enabling the app")
		}

		return errors.New("An error occured while disabling the app")
	}

	return nil
}
