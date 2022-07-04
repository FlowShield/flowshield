package ca

import (
	"github.com/cloudSlit/cloudslit/ca/api/helper"
	logic "github.com/cloudSlit/cloudslit/ca/logic/ca"
)

func init() {
	// load type...
	logic.DoNothing()
}

// IntermediateTopology Sub-CA topology
// @Tags CA
// @Summary Sub-CA topology
// @Description Sub-CA topology
// @Produce json
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=[]logic.IntermediateObject} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /ca/intermediate_topology [get]
func (a *API) IntermediateTopology(c *helper.HTTPWrapContext) (interface{}, error) {
	return a.logic.IntermediateTopology()
}

// UpperCaIntermediateTopology Upper CA topology
// @Tags CA
// @Summary Upper CA topology
// @Description Upper CA topology
// @Produce json
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=[]logic.IntermediateObject} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /ca/upper_ca_intermediate_topology [get]
func (a *API) UpperCaIntermediateTopology(c *helper.HTTPWrapContext) (interface{}, error) {
	body, err := a.logic.UpperCaIntermediateTopology()
	if err != nil {
		return nil, err
	}

	return body, nil
}
