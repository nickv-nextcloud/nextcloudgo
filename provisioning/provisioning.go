package provisioning

import (
	"github.com/nextcloud/nextcloudgo"
	"github.com/nextcloud/nextcloudgo/ocs"
)

var (
	endpoint = "ocs/v2.php/cloud"
)

type Provisioning struct {
	sdk nextcloudgo.NextcloudGo
	ocs ocs.Request
}

func New(sdk nextcloudgo.NextcloudGo) Provisioning {
	ocs := ocs.New(sdk)
	return Provisioning{sdk: sdk, ocs: ocs}
}
