// Package workload Certificate Lifecycle Management
package workload

import (
	"github.com/cloudSlit/cloudslit/ca/api/helper"
	logic "github.com/cloudSlit/cloudslit/ca/logic/workload"
)

// RevokeCerts revoked certificate
// @Tags Workload
// @Summary (p3)Revoke
// @Description revoked certificate
// @Produce json
// @Param body body logic.RevokeCertsParams true "sn+aki / unique_id pick one of two"
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/lifecycle/revoke [post]
func (a *API) RevokeCerts(c *helper.HTTPWrapContext) (interface{}, error) {
	var req logic.RevokeCertsParams
	c.BindG(&req)

	err := a.logic.RevokeCerts(&req)
	if err != nil {
		return nil, err
	}

	return "revoked", nil
}

// RecoverCerts Restore certificate
// @Tags Workload
// @Summary (p3)Recover
// @Description Restore certificate
// @Produce json
// @Param body body logic.RecoverCertsParams true "sn+aki / unique_id either-or"
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/lifecycle/recover [post]
func (a *API) RecoverCerts(c *helper.HTTPWrapContext) (interface{}, error) {
	var req logic.RecoverCertsParams
	c.BindG(&req)

	err := a.logic.RecoverCerts(&req)
	if err != nil {
		return nil, err
	}

	return "recovered", nil
}

// ForbidNewCerts Prohibit a uniqueID from requesting a certificate
// @Tags Workload
// @Summary Application for certificate is prohibited
// @Description Prohibit a uniqueID from requesting a certificate
// @Produce json
// @Param body body logic.ForbidNewCertsParams true " "
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/lifecycle/forbid_new_certs [post]
func (a *API) ForbidNewCerts(c *helper.HTTPWrapContext) (interface{}, error) {
	var req logic.ForbidNewCertsParams
	c.BindG(&req)

	err := a.logic.ForbidNewCerts(&req)
	if err != nil {
		return nil, err
	}

	return "success", nil
}

// RecoverForbidNewCerts Recovery allows a uniqueID to request a certificate
// @Tags Workload
// @Summary Resume application certificate
// @Description Recovery allows a uniqueID to request a certificate
// @Produce json
// @Param body body logic.ForbidNewCertsParams true " "
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/lifecycle/recover_forbid_new_certs [post]
func (a *API) RecoverForbidNewCerts(c *helper.HTTPWrapContext) (interface{}, error) {
	var req logic.ForbidNewCertsParams
	c.BindG(&req)

	err := a.logic.RecoverForbidNewCerts(&req)
	if err != nil {
		return nil, err
	}

	return "success", nil
}

type ForbidUnitParams struct {
	UniqueID string `json:"unique_id" binding:"required"`
}

// ForbidUnit Revoke and prohibit service certificates
// @Tags Workload
// @Summary (p1)Revoke and prohibit service certificates
// @Description Revoke and prohibit service certificates
// @Produce json
// @Param json body ForbidUnitParams true " "
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/lifecycle/forbid_unit [post]
func (a *API) ForbidUnit(c *helper.HTTPWrapContext) (interface{}, error) {
	var req ForbidUnitParams
	c.BindG(&req)

	err := a.logic.ForbidNewCerts(&logic.ForbidNewCertsParams{
		UniqueIds: []string{req.UniqueID},
	})
	if err != nil {
		a.logger.With("req", req).Errorf("Failed to prohibit certificate application: %s", err)
		return nil, err
	}

	// 2021.04.15 (functional adjustment) certificate enabling and disabling will affect certificate communication, OCSP authentication and MTLs use, and the certificate will not be revoked
	// err = a.logic.RevokeCerts(&logic.RevokeCertsParams{
	//    UniqueId: req.UniqueID,
	// })
	// if err != nil {
	//    a.logger.With("req", req).Errorf("Revocation of service certificate failed: %s", err)
	//    return nil, err
	// }

	return "success", nil
}

// RecoverUnit Restore and allow service certificates
// @Tags Workload
// @Summary (p1)Restore and allow service certificates
// @Description Restore and allow service certificates
// @Produce json
// @Param json body ForbidUnitParams true " "
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/lifecycle/recover_unit [post]
func (a *API) RecoverUnit(c *helper.HTTPWrapContext) (interface{}, error) {
	var req ForbidUnitParams
	c.BindG(&req)

	err := a.logic.RecoverForbidNewCerts(&logic.ForbidNewCertsParams{
		UniqueIds: []string{req.UniqueID},
	})
	if err != nil {
		a.logger.With("req", req).Errorf("Failed to restore the requested certificate: %s", err)
		return nil, err
	}

	// err = a.logic.RecoverCerts(&logic.RecoverCertsParams{
	//    UniqueId: req.UniqueID,
	// })
	// if err != nil {
	//    a.logger.With("req", req).Errorf("Failed to restore service certificate: %s", err)
	//    return nil, err
	// }

	return "success", nil
}
