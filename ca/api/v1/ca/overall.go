package ca

import (
	"strings"
	"time"

	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/cloudslit/cloudslit/ca/api/helper"
	"github.com/cloudslit/cloudslit/ca/core"
	"github.com/cloudslit/cloudslit/ca/database/mysql/cfssl-model/model"
)

type OverallCertsCountItem struct {
	Role       string `json:"role"`
	Total      int64  `json:"total"`       // Total number of certificates
	UnitsCount int64  `json:"units_count"` // number of services
}

type OverallCertsCountResponse struct {
	Total int64                   `json:"total"`
	Certs []OverallCertsCountItem `json:"certs"`
}

// OverallCertsCount Certificate classification
// @Tags CA
// @Summary (p2)Certificate classification
// @Description Total number of certificates, number by classification, number of corresponding services
// @Produce json
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=OverallCertsCountResponse} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /ca/overall_certs_count [get]
func (a *API) OverallCertsCount(c *helper.HTTPWrapContext) (interface{}, error) {
	query := func() *gorm.DB {
		return core.Is.Db.Session(&gorm.Session{}).Model(&model.Certificates{}).
			Where("expiry > ?", time.Now()).
			Where("reason IN ?", []int{0, 2})
	}

	var total int64
	if err := query().Count(&total).Error; err != nil {
		a.logger.Errorf("mysql query err: %s", err)
		return nil, err
	}

	roleProfiles, err := a.logic.RoleProfiles()
	if err != nil {
		a.logger.Errorf("Error getting role profiles: %s", err)
		return nil, errors.New("Error getting role profiles")
	}

	res := &OverallCertsCountResponse{
		Total: total,
		Certs: make([]OverallCertsCountItem, 0),
	}

	for _, roleProfile := range roleProfiles {
		role := roleProfile.Name
		item := OverallCertsCountItem{Role: role}
		if err := query().Where("ca_label = ?", strings.ToLower(role)).Count(&item.Total).Error; err != nil {
			a.logger.Errorf("mysql query err: %s", err)
			return nil, err
		}

		if err := query().Where("ca_label = ?", strings.ToLower(role)).
			Where(`common_name != ""`).Group("common_name").Count(&item.UnitsCount).Error; err != nil {
			a.logger.Errorf("mysql query err: %s", err)
			return nil, err
		}
		res.Certs = append(res.Certs, item)
	}

	return res, nil
}

type OverallExpiryGroup struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type OverallExpiryCertsResponse struct {
	ExpiryTotal int64                `json:"expiry_total"`
	ExpiryCerts []OverallExpiryGroup `json:"expiry_certs"`
}

// OverallExpiryCerts Certificate validity
// @Tags CA
// @Summary (p2)Certificate validity
// @Description Number of certificates expired: within one week, within 1/3 months and after 3 months
// @Produce json
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=OverallExpiryCertsResponse} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /ca/overall_expiry_certs [get]
func (a *API) OverallExpiryCerts(c *helper.HTTPWrapContext) (interface{}, error) {
	query := func() *gorm.DB {
		return core.Is.Db.Model(&model.Certificates{}).
			Where("expiry < ?", time.Now()).
			Where("revoked_at IS NULL")
	}

	var total int64
	if err := query().Count(&total).Error; err != nil {
		a.logger.Errorf("mysql query err: %s", err)
		return nil, err
	}

	res := &OverallExpiryCertsResponse{
		ExpiryTotal: total,
		ExpiryCerts: make([]OverallExpiryGroup, 0),
	}

	// Within a week
	{
		item := OverallExpiryGroup{Name: "1w"}
		count, err := getExpiryCountByDuration(7*24*time.Hour, time.Now())
		if err != nil {
			return nil, err
		}
		item.Count = count
		res.ExpiryCerts = append(res.ExpiryCerts, item)
	}

	// Within one month
	{
		item := OverallExpiryGroup{Name: "1m"}
		count, err := getExpiryCountByDuration(30*24*time.Hour, time.Now())
		if err != nil {
			return nil, err
		}
		item.Count = count
		res.ExpiryCerts = append(res.ExpiryCerts, item)
	}

	// Within three months
	{
		item := OverallExpiryGroup{Name: "3m"}
		count, err := getExpiryCountByDuration(3*30*24*time.Hour, time.Now())
		if err != nil {
			return nil, err
		}
		item.Count = count
		res.ExpiryCerts = append(res.ExpiryCerts, item)
	}

	// Three months later
	{
		item := OverallExpiryGroup{Name: "3m+"}
		count, err := getExpiryCountByDuration(999*30*24*time.Hour, time.Now().AddDate(0, 3, 0))
		if err != nil {
			return nil, err
		}
		item.Count = count
		res.ExpiryCerts = append(res.ExpiryCerts, item)
	}

	return res, nil
}

func getExpiryCountByDuration(period time.Duration, before time.Time) (int64, error) {
	// Within a week
	// Expiration time - current time < = one week
	// Expiration time < = current time + one week
	expiryDate := time.Now().Add(period)
	query := core.Is.Db.Session(&gorm.Session{}).Model(&model.Certificates{}).
		Where("expiry > ?", before).
		Where("expiry < ?", expiryDate).
		Where("reason = 0").
		Where(`common_name != ""`)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		logger.Errorf("mysql query err: %s", err)
		return 0, err
	}

	return count, nil
}

type OverallUnitsEnableItem struct {
	CertsCount int64 `json:"certs_count"`
	UnitsCount int64 `json:"units_count"`
}

type OverallUnitsEnableStatus struct {
	Enable  OverallUnitsEnableItem `json:"enable"`
	Disable OverallUnitsEnableItem `json:"disable"`
}

// OverallUnitsEnableStatus Enabling condition
// @Tags CA
// @Summary (p2)Enabling condition
// @Description Total enabled, total disabled, corresponding services
// @Produce json
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=OverallUnitsEnableStatus} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /ca/overall_units_enable_status [get]
func (a *API) OverallUnitsEnableStatus(c *helper.HTTPWrapContext) (interface{}, error) {
	query := func() *gorm.DB {
		return core.Is.Db.Session(&gorm.Session{}).Model(&model.Certificates{}).
			Where("expiry > ?", time.Now())
	}

	res := &OverallUnitsEnableStatus{}

	{
		if err := query().Where("reason = ?", 0).Count(&res.Enable.CertsCount).Error; err != nil {
			a.logger.Errorf("mysql query err: %s", err)
			return nil, err
		}
		if err := query().Where("reason = ?", 2).Count(&res.Disable.CertsCount).Error; err != nil {
			a.logger.Errorf("mysql query err: %s", err)
			return nil, err
		}
	}

	{
		if err := query().Where("reason = ?", 0).Group("common_name").Count(&res.Enable.UnitsCount).Error; err != nil {
			a.logger.Errorf("mysql query err: %s", err)
			return nil, err
		}
		if err := query().Where("reason = ?", 2).Group("common_name").Count(&res.Disable.UnitsCount).Error; err != nil {
			a.logger.Errorf("mysql query err: %s", err)
			return nil, err
		}
	}

	return res, nil
}
