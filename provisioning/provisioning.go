// Package provisioning allows to manage apps, users and groups on a nextcloud instance.
// Unless specified otherwise the functions can only be used with an admin user.
package provisioning

import (
	"github.com/nextcloud/nextcloudgo"
	"github.com/nextcloud/nextcloudgo/ocs"
)

var (
	endpoint = "ocs/v2.php/cloud"
)

// Provisioning allows to manage apps, users and groups on a nextcloud instance
type Provisioning struct {
	sdk nextcloudgo.NextcloudGo
	ocs ocs.Request
}

// New returns a new Provisioning instance when given the sdk
func New(sdk nextcloudgo.NextcloudGo) Provisioning {
	ocs := ocs.New(sdk)
	return Provisioning{sdk: sdk, ocs: ocs}
}
