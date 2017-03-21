// Package nextcloudgo.share contains utility functions for working with shares.
package sharing

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nextcloud/nextcloudgo"
	"github.com/nextcloud/nextcloudgo/ocs"
)

const TypeUser = 0
const TypeGroup = 1
const TypeLink = 3
const TypeMail = 4
const TypeRemote = 6

const PermissionRead = 1
const PermissionUpdate = 2
const PermissionCreate = 4
const PermissionDelete = 8
const PermissionShare = 16
const PermissionAll = 31

var (
	endpoint = "ocs/v2.php/apps/files_sharing/api/v1"
)

type Share struct {
	Id   int
	Type int

	Owner                string
	OwnerDisplayName     string
	Initiator            string
	InitiatorDisplayName string

	Path        string
	Permissions int
	Time        time.Time

	// UserShare
	// GroupShare
	// FederatedShare
	// MailShare
	With            string
	WithDisplayName string

	// FederatedShare
	// MailShare
	// LinkShare
	Token string

	// LinkShare
	Expiration time.Time
}

type Sharing struct {
	sdk nextcloudgo.NextcloudGo
	ocs ocs.Request
}

func New(sdk nextcloudgo.NextcloudGo) Sharing {
	ocs := ocs.New(sdk)
	return Sharing{sdk: sdk, ocs: ocs}
}

func (sharing *Sharing) GetShareById(id int) (Share, error) {
	url := sharing.sdk.GetServerUrl()
	url += endpoint + "/shares/" + strconv.Itoa(id)

	content, err := sharing.ocs.NewRequest(http.MethodGet, url+"?format=json")
	if err != nil {
		return Share{}, err
	}

	if !ocs.ValidateStatusCode(content, 200) {
		return Share{}, errors.New("Status code was invalid")
	}

	ocs := content["ocs"].(map[string]interface{})
	data := ocs["data"].([]interface{})
	share := data[0].(map[string]interface{})
	return sharing.createShareFromMap(share)
}

func (sharing *Sharing) createShareFromMap(share map[string]interface{}) (Share, error) {
	log.Println(share)
	s := Share{}
	s.Id, _ = share["id"].(int)
	s.Type, _ = share["share_type"].(int)

	s.Owner = share["uid_file_owner"].(string)
	s.OwnerDisplayName = share["displayname_file_owner"].(string)
	s.Initiator = share["uid_owner"].(string)
	s.InitiatorDisplayName = share["displayname_owner"].(string)
	s.With = share["share_with"].(string)
	s.WithDisplayName = share["share_with_displayname"].(string)

	s.Path = share["path"].(string)
	if share["token"] != nil {
		s.Token = share["token"].(string)
	}
	s.Permissions, _ = share["permissions"].(int)

	s.Time = time.Unix(int64(share["stime"].(float64)), 0)
	if share["expiration"] != nil {
		s.Expiration = time.Unix(int64(share["expiration"].(float64)), 0)
	}

	return s, nil

}
