package provisioning

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/nextcloud/nextcloudgo/ocs"
)

var (
	// ErrUserDoesNotExist when the user does not exist
	ErrUserDoesNotExist = errors.New("User does not exist")
	// ErrUserAlreadyExists when the user already exists
	ErrUserAlreadyExists = errors.New("User already exists")
)

// TODO missing GetUserData
// TODO missing SetUserData

// CreateUser creates the user on the server if it does not exist yet
// Returns ErrUserAlreadyExists when the user already exists
func (api *Provisioning) CreateUser(userid, password string) error {
	body := map[string]string{"userid": userid, "password": password}
	reader := new(bytes.Buffer)
	json.NewEncoder(reader).Encode(body)

	url := endpoint + "/users?format=json"
	content, status, err := api.ocs.RequestWithBody(http.MethodPost, url, reader, true)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		if ocs.ValidateStatusCode(content, 102) {
			return ErrUserAlreadyExists
		}
		return errors.New("An error occured while creating the user")
	}

	return nil
}

// DeleteUser deletes the given user from the server
// The current user can not be deleted
// Returns ErrUserDoesNotExist when the user does not exist
func (api *Provisioning) DeleteUser(userid string) error {
	url := endpoint + "/users/" + userid + "?format=json"
	_, status, err := api.ocs.Request(http.MethodDelete, url, true)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.New("An error occured while deleting the user")
	}

	return nil
}

// GetUsers returns all users from the server matching the given search
// The search is with wildcards on both ends
// Setting start and limit to 0 will return all users
func (api *Provisioning) GetUsers(search string, start, limit int) ([]string, error) {
	url := endpoint + "/users?format=json"
	if search != "" {
		url += "&search=" + search
	}
	if start > 0 {
		url += "&start=" + strconv.Itoa(start)
	}
	if limit > 0 {
		url += "&limit=" + strconv.Itoa(limit)
	}

	content, status, err := api.ocs.Request(http.MethodGet, url, true)
	if err != nil {
		return []string{}, err
	}

	if status != http.StatusOK {
		return []string{}, errors.New("An error occured while searching for groups")
	}

	return ocs.GetStringList(content, []string{"ocs", "data", "users"})
}

// UserExists checks whether a user exists on the server
func (api *Provisioning) UserExists(userid string) (bool, error) {
	users, err := api.GetUsers(userid, 0, 0)
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if user == userid {
			return true, nil
		}
	}

	return false, nil
}

// EnableUser enables a disabled user
func (api *Provisioning) EnableUser(userid string) error {
	return api.changeUserState(userid, "enable")
}

// DisableUser disables an enabled user
func (api *Provisioning) DisableUser(userid string) error {
	return api.changeUserState(userid, "disable")

}

func (api *Provisioning) changeUserState(userid string, state string) error {
	url := endpoint + "/users/" + userid + "/" + state + "?format=json"

	_, status, err := api.ocs.Request(http.MethodPut, url, true)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		if state == "enable" {
			return errors.New("An error occured while enabling the user")
		}

		return errors.New("An error occured while disabling the user")
	}

	return nil
}

// AddUserToGroup adds the user to the given group
// Returns ErrGroupDoesNotExist when the group does not exist
// Returns ErrUserDoesNotExist when the user does not exist
func (api *Provisioning) AddUserToGroup(userid, groupid string) error {
	return api.changeUserGroupMemberState(userid, groupid, http.MethodPost)
}

// RemoveUserFromGroup removes the user from the given group
// Returns ErrGroupDoesNotExist when the group does not exist
// Returns ErrUserDoesNotExist when the user does not exist
func (api *Provisioning) RemoveUserFromGroup(userid, groupid string) error {
	return api.changeUserGroupMemberState(userid, groupid, http.MethodDelete)
}

func (api *Provisioning) changeUserGroupMemberState(userid, groupid, method string) error {
	body := map[string]string{"groupid": groupid}
	reader := new(bytes.Buffer)
	json.NewEncoder(reader).Encode(body)

	url := endpoint + "/users/" + userid + "/groups?format=json"
	content, status, err := api.ocs.RequestWithBody(method, url, reader, true)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		if ocs.ValidateStatusCode(content, 102) {
			return ErrGroupDoesNotExist
		}
		if ocs.ValidateStatusCode(content, 103) {
			return ErrUserDoesNotExist
		}
		if method == http.MethodPost {
			return errors.New("An error occured while adding the user to the group")
		}

		return errors.New("An error occured while removing the user from the group")
	}

	return nil
}

// GetUserGroups returns all groups the user is a member of
func (api *Provisioning) GetUserGroups(userid string) ([]string, error) {
	url := endpoint + "/users/" + userid + "/groups?format=json"

	content, status, err := api.ocs.Request(http.MethodGet, url, true)
	if err != nil {
		return []string{}, err
	}

	if status == http.StatusNotFound {
		return []string{}, ErrUserDoesNotExist
	}

	if status != http.StatusOK {
		return []string{}, errors.New("An error occured while searching for groups")
	}

	return ocs.GetStringList(content, []string{"ocs", "data", "groups"})
}
