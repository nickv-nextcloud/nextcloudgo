// Package provisioning allows to manage apps, users and groups on a nextcloud instance.
// Unless specified otherwise the functions can only be used with an admin user.
package provisioning

import (
	"github.com/nextcloud/nextcloudgo"
	"github.com/nextcloud/nextcloudgo/ocs"
)

var (
	endpoint = "/ocs/v2.php/cloud"
)

// Provisioning allows to manage apps, users and groups on a nextcloud instance
type Provisioning struct {
	nc  nextcloudgo.NextcloudGo
	ocs ocs.Request
}

// New returns a new Provisioning instance when given the NextcloudGo
func New(nc nextcloudgo.NextcloudGo) Provisioning {
	ocs := ocs.New(nc)
	return Provisioning{nc: nc, ocs: ocs}
}
