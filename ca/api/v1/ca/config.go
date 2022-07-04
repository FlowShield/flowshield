package ca

import (
	"strings"

	"github.com/cloudSlit/cloudslit/ca/api/helper"
	logic "github.com/cloudSlit/cloudslit/ca/logic/ca"
)

func init() {
	// Load type...
	logic.DoNothing()
}

// RoleProfiles Environmental isolation type
// @Tags CA
// @Summary (p1)Environmental isolation type
// @Description Environmental isolation type
// @Produce json
// @Param short query bool false "Only a list of types is returned for search criteria"
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=logic.RoleProfile} " "
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=[]string} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /ca/role_profiles [get]
func (a *API) RoleProfiles(c *helper.HTTPWrapContext) (interface{}, error) {
	profiles, err := a.logic.RoleProfiles()
	if err != nil {
		return nil, err
	}
	if c.G.Query("short") == "true" {
		roles := make([]string, 0, len(profiles)-1)
		for _, profile := range profiles {
			if strings.ToLower(profile.Name) == "default" {
				continue
			}
			roles = append(roles, strings.ToLower(profile.Name))
		}
		return roles, nil
	}
	return profiles, nil
}
