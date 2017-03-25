package provisioning

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nextcloud/nextcloudgo/ocs"
)

var (
	// ErrGroupDoesNotExist when the group does not exist
	ErrGroupDoesNotExist = errors.New("Group does not exist")
	// ErrGroupAlreadyExists when the group already exists
	ErrGroupAlreadyExists = errors.New("Group already exists")
)

// CreateGroup creates the group on the server if it does not exist yet
// Returns ErrGroupAlreadyExists when the group already exists
func (api *Provisioning) CreateGroup(groupid string) error {
	body := map[string]string{"groupid": groupid}
	reader := new(bytes.Buffer)
	json.NewEncoder(reader).Encode(body)

	url := endpoint + "/groups?format=json"
	content, status, err := api.ocs.RequestWithBody(http.MethodPost, url, reader, true)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		if ocs.ValidateStatusCode(content, 101) {
			return errors.New("Provided group name is invalid")
		}
		if ocs.ValidateStatusCode(content, 102) {
			return ErrGroupAlreadyExists
		}
		return errors.New("An error occured while creating the group")
	}

	return nil
}

// DeleteGroup deletes the given group from the server
// The group admin can not be deleted
// Returns ErrGroupDoesNotExist when the group does not exist
func (api *Provisioning) DeleteGroup(groupid string) error {
	if groupid == "admin" {
		return errors.New("Admin group can not be deleted")
	}

	url := endpoint + "/groups/" + groupid + "?format=json"
	content, status, err := api.ocs.Request(http.MethodDelete, url, true)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		if ocs.ValidateStatusCode(content, 101) {
			return ErrGroupDoesNotExist
		}
		return errors.New("An error occured while deleting the group")
	}

	return nil
}

// GetGroups returns all groups from the server matching the given search
// The search is with wildcards on both ends
func (api *Provisioning) GetGroups(search string) ([]string, error) {
	var url string
	if search != "" {
		url = endpoint + "/groups?format=json&search=" + search
	} else {
		url = endpoint + "/groups?format=json"
	}

	content, status, err := api.ocs.Request(http.MethodGet, url, true)
	if err != nil {
		return []string{}, err
	}

	if status != http.StatusOK {
		return []string{}, errors.New("An error occured while searching for groups")
	}

	return ocs.GetStringList(content, []string{"ocs", "data", "groups"})
}

// GetGroupMembers returns all users that are members of the given group
// Returns ErrGroupDoesNotExist when the group does not exist
func (api *Provisioning) GetGroupMembers(groupid string) ([]string, error) {
	url := endpoint + "/groups/" + groupid + "?format=json"
	content, status, err := api.ocs.Request(http.MethodGet, url, true)
	if err != nil {
		return []string{}, err
	}

	if status == http.StatusNotFound {
		return []string{}, ErrGroupDoesNotExist
	}
	if status != http.StatusOK {
		return []string{}, errors.New("An error occured while getting the members of the group")
	}

	return ocs.GetStringList(content, []string{"ocs", "data", "users"})
}

// GroupExists checks whether a group exists on the server
func (api *Provisioning) GroupExists(groupid string) (bool, error) {
	groups, err := api.GetGroups(groupid)
	if err != nil {
		return false, err
	}

	for _, group := range groups {
		if group == groupid {
			return true, nil
		}
	}

	return false, nil
}
