package resources

import (
	"fmt"

	"github.com/fission/fission/pkg/controller/client"
	"github.com/fission/fission/pkg/fission-cli/util"
)

//FissionVersionLabel is label for FissionVersion check
const FissionVersionLabel = "determine fission server version"

// FissionVersion is used to check if it is possible to determine fission server version
type FissionVersion struct {
	client client.Interface
}

// NewFissionVersion is used to create new FissionVersion instance
func NewFissionVersion(client client.Interface) Resource {
	return FissionVersion{client: client}
}

// Check performs the check and returns the result
func (res FissionVersion) Check() Results {
	ver := util.GetVersion(res.client)
	if ver.Server["fission/core"].Version == "" {
		return getResults("not able to determine fission server version", false)
	}
	return getResults(fmt.Sprintf("able to determine fission server version: %s", ver.Server["fission/core"].Version), true)
}

// GetLabel returns the label for check
func (res FissionVersion) GetLabel() string {
	return FissionVersionLabel
}
