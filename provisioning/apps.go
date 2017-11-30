package provisioning

import (
	"errors"
	"net/http"

	"github.com/nextcloud/nextcloudgo/ocs"
)

var (
	// ErrAppDoesNotExist when the app does not exist
	ErrAppDoesNotExist = errors.New("App does not exist")
)

// GetApps returns the list of apps matching the given filter
// Valid values for filter are: enabled, disabled, all
func (api *Provisioning) GetApps(filter string) ([]string, error) {
	if filter != "enabled" && filter != "disabled" && filter != "all" {
		return []string{}, errors.New("Invalid filter given")
	}

	url := endpoint + "/apps"
	if filter == "enabled" || filter == "disabled" {
		url = url + "?filter=" + filter
	}

	content, status, err := api.ocs.Request(http.MethodGet, url, true)
	if err != nil {
		return []string{}, err
	}

	if status != http.StatusOK {
		return []string{}, errors.New("An error occured while search for apps")
	}

	return ocs.GetStringList(content, []string{"ocs", "data", "apps"})
}

// IsAppEnabled returns true when the app is enabled, false otherwise
func (api *Provisioning) IsAppEnabled(appid string) (bool, error) {
	return api.isAppInArray(appid, "enabled")
}

// IsAppDisabled returns true when the app is disabled but available, false otherwise
func (api *Provisioning) IsAppDisabled(appid string) (bool, error) {
	return api.isAppInArray(appid, "disabled")
}

// IsAppAvailable returns true when the app is available, false otherwise
func (api *Provisioning) IsAppAvailable(appid string) (bool, error) {
	return api.isAppInArray(appid, "all")
}

func (api *Provisioning) isAppInArray(appid, filter string) (bool, error) {
	apps, err := api.GetApps(filter)
	if err != nil {
		return false, err
	}

	for _, app := range apps {
		if app == appid {
			return true, nil
		}
	}

	return false, nil
}

// EnableApp enables an app when it is available
// Returns ErrAppDoesNotExist when the app does not exist
func (api *Provisioning) EnableApp(appid string) error {
	return api.changeAppState(appid, http.MethodPost)
}

// DisableApp disables an app when it is available
// Returns ErrAppDoesNotExist when the app does not exist
func (api *Provisioning) DisableApp(appid string) error {
	return api.changeAppState(appid, http.MethodDelete)
}

func (api *Provisioning) changeAppState(appid, method string) error {
	url := endpoint + "/apps/" + appid

	_, status, err := api.ocs.Request(method, url, true)
	if err != nil {
		return err
	}

	if status == http.StatusNotFound {
		return ErrAppDoesNotExist
	}

	if status != http.StatusOK {
		if method == http.MethodPost {
			return errors.New("An error occured while enabling the app")
		}

		return errors.New("An error occured while disabling the app")
	}

	return nil
}
