package ca

import (
	"github.com/cloudslit/cloudslit/ca/api/helper"
	logic "github.com/cloudslit/cloudslit/ca/logic/ca"
)

func init() {
	// load type...
	logic.DoNothing()
}

// WorkloadUnits CA Units
// 	UniqueID as unit
// @Tags CA
// @Summary (p1)Service unit
// @Description CA Units
// @Produce json
// @Param page query int false "Number of pages, default 1"
// @Param limit_num query int false "Page limit, default 20"
// @Param unique_id query string false "UniqueID Query"
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=helper.MSPNormalizeList{list=[]logic.WorkloadUnit}} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /ca/workload_units [get]
func (a *API) WorkloadUnits(c *helper.HTTPWrapContext) (interface{}, error) {
	var req = struct {
		helper.MSPNormalizeListPaginateParams
		UniqueId string `form:"unique_id"`
	}{
		MSPNormalizeListPaginateParams: helper.DefaultMSPNormalizeListPaginateParams,
	}
	c.BindG(&req)

	data, total, err := a.logic.WorkloadUnits(&logic.WorkloadUnitsParams{
		Page:     req.Page,
		PageSize: req.LimitNum,
		UniqueId: req.UniqueId,
	})
	if err != nil {
		return nil, err
	}

	list := helper.MSPNormalizeList{
		List: data,
		Paginate: helper.MSPNormalizePaginate{
			Total:    total,
			Current:  req.Page,
			PageSize: req.LimitNum,
		},
	}
	return list, nil
}
