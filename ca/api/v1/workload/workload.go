package workload

import (
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/cloudSlit/cloudslit/ca/api/helper"
	"github.com/cloudSlit/cloudslit/ca/core"
	"github.com/cloudSlit/cloudslit/ca/database/mysql/cfssl-model/dao"
	"github.com/cloudSlit/cloudslit/ca/database/mysql/cfssl-model/model"
	"github.com/cloudSlit/cloudslit/ca/logic/schema"
	logic "github.com/cloudSlit/cloudslit/ca/logic/workload"
)

type API struct {
	logic  *logic.Logic
	logger *zap.SugaredLogger
}

func NewAPI() *API {
	return &API{
		logic:  logic.NewLogic(),
		logger: logger.Named("api").SugaredLogger,
	}
}

// CertList Certificate list
// @Tags Workload
// @Summary (p3)List
// @Description Certificate list
// @Produce json
// @Param role query string false "Certificate type default"
// @Param unique_id query string false "Query by unique ID"
// @Param cert_sn query string false "Query by certificate serial number"
// @Param status query string false "Certificate status good/revoked"
// @Param order query string false "Sort, default issued_at desc"
// @Param expiry_start_time query string false "Expiration, starting point"
// @Param expiry_end_time query string false "Expiration, end time point"
// @Param limit_num query int false "Paging parameters, default 20"
// @Param page query int false "Number of pages, default 1"
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=helper.MSPNormalizeList{list=[]schema.SampleCert}} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/certs [get]
func (a *API) CertList(c *helper.HTTPWrapContext) (interface{}, error) {
	var req = struct {
		Role            string `form:"role"`
		UniqueID        string `form:"unique_id"`
		Status          string `form:"status"`
		Order           string `form:"order"`
		CertSN          string `form:"cert_sn"`
		ExpiryStartTime string `form:"expiry_start_time"`
		ExpiryEndTime   string `form:"expiry_end_time"`
		helper.MSPNormalizeListPaginateParams
	}{
		MSPNormalizeListPaginateParams: helper.DefaultMSPNormalizeListPaginateParams,
	}
	c.BindG(&req)

	data, err := a.logic.CertList(&logic.CertListParams{
		Role:            req.Role,
		UniqueID:        req.UniqueID,
		CertSN:          req.CertSN,
		Page:            req.Page,
		PageSize:        req.LimitNum,
		Status:          req.Status,
		Order:           req.Order,
		ExpiryStartTime: req.ExpiryStartTime,
		ExpiryEndTime:   req.ExpiryEndTime,
	})
	if err != nil {
		return nil, err
	}

	result := helper.MSPNormalizeList{
		List: data.CertList,
		Paginate: helper.MSPNormalizePaginate{
			Total:    data.Total,
			Current:  req.Page,
			PageSize: req.LimitNum,
		},
	}
	return result, nil
}

// CertDetail Certificate details
// @Tags Workload
// @Summary Detail
// @Description Certificate details
// @Produce json
// @Param sn query string true "Certificate sn"
// @Param aki query string true "Certificate aki"
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=schema.FullCert} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/cert [get]
func (a *API) CertDetail(c *helper.HTTPWrapContext) (interface{}, error) {
	var req struct {
		SN  string `form:"sn" binding:"required"`
		AKI string `form:"aki" binding:"required"`
	}
	c.BindG(&req)

	data, err := a.logic.CertDetail(&logic.CertDetailParams{
		SN:  req.SN,
		AKI: req.AKI,
	})
	if err != nil {
		return nil, err
	}

	result := struct {
		Cert     interface{} `json:"cert"`
		CertInfo string      `json:"cert_info"`
	}{
		Cert: data,
	}

	result.CertInfo = data.CertStr
	return result, nil
}

// UnitsForbidQuery Query unique_ Is ID prohibited from applying for certificate
// @Tags Workload
// @Summary Prohibit applying for certificate query
// @Description Query unique_id Is it forbidden to apply for certificate
// @Produce json
// @Param unique_ids query []string true "Query unique_ID array" collectionFormat(multi)
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=logic.UnitsForbidQueryResult} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/units_forbid_query [get]
func (a *API) UnitsForbidQuery(c *helper.HTTPWrapContext) (interface{}, error) {
	var req struct {
		UniqueIds []string `form:"unique_ids" binding:"required"`
	}
	c.BindG(&req)

	return a.logic.UnitsForbidQuery(&logic.UnitsForbidQueryParams{UniqueIds: req.UniqueIds})
}

type UnitsStatusItem struct {
	Active bool `json:"active"`
}

type UnitsStatusMap map[string]*UnitsStatusItem

type UnitsStatusReq struct {
	UniqueIds []string `json:"unique_ids" binding:"required"`
}

// UnitsStatus Service corresponding status query
// @Tags Workload
// @Summary (p1)Service corresponding status query
// @Description Service corresponding status query
// @Produce json
// @Param json body UnitsStatusReq true "Query unique_ID array"
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=object} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/units_status [post]
func (a *API) UnitsStatus(c *helper.HTTPWrapContext) (interface{}, error) {
	var req UnitsStatusReq
	c.BindG(&req)

	statusMap := make(UnitsStatusMap)
	var err error

	fx.From(func(source chan<- interface{}) {
		for _, uid := range req.UniqueIds {
			if len(uid) == 0 {
				continue
			}
			source <- uid
		}
	}).Split(300).ForEach(func(obj interface{}) {
		group := obj.([]interface{})
		var uids []string
		for _, item := range group {
			uids = append(uids, item.(string))
		}
		var sm UnitsStatusMap
		sm, err = a.getUnitsStatus(uids)
		if err != nil {
			return
		}
		for k, v := range sm {
			statusMap[k] = v
		}
	})

	if err != nil {
		return nil, err
	}

	return statusMap, nil
}

func (a *API) getUnitsStatus(uniqueIds []string) (UnitsStatusMap, error) {
	db := core.Is.Db
	query := db.Where("common_name IN ?", uniqueIds).
		Where("status = ?", "good").
		Where("expiry > ?", time.Now()).
		Select("common_name").
		Group("common_name")

	var list []*model.Certificates
	if err := query.Find(&list).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.logger.Errorf("Database query error: %s", err)
		return nil, err
	}

	statusMap := make(UnitsStatusMap)
	for _, uid := range uniqueIds {
		statusMap[uid] = &UnitsStatusItem{
			Active: false,
		}
	}

	for _, row := range list {
		if row.CommonName.Valid {
			uid := row.CommonName.String
			statusMap[uid] = &UnitsStatusItem{
				Active: true,
			}
		}
	}

	return statusMap, nil
}

type UnitsCertsItem struct {
	Certs     []*schema.FullCert `json:"certs"`
	UniqueID  string             `json:"unique_id"`
	Forbidden bool               `json:"forbidden"`
}

// UnitsCertsList List of service certificates
// Deprecated
// @Tags Workload
// @Summary (p1)List of service certificates
// @Description List of service certificates
// @Produce json
// @Param unique_id query string false "Query unique_id"
// @Param role query string false "Certificate type"
// @Param expiry_start_time query string false "Expiration, starting point"
// @Param expiry_end_time query string false "Expiration, end time point"
// @Param is_forbid query int false "Disable, 1 disable, 2 enable"
// @Param limit_num query int false "Paging parameters, default 20"
// @Param page query int false "Number of pages, default 1"
// @Success 200 {object} helper.MSPNormalizeHTTPResponseBody{data=[]UnitsCertsItem} " "
// @Failure 400 {object} helper.HTTPWrapErrorResponse
// @Failure 500 {object} helper.HTTPWrapErrorResponse
// @Router /workload/units_certs_list [get]
func (a *API) UnitsCertsList(c *helper.HTTPWrapContext) (interface{}, error) {
	var req = struct {
		UniqueID        string `form:"unique_id"`
		Role            string `form:"role"`
		ExpiryStartTime string `form:"expiry_start_time"`
		ExpiryEndTime   string `form:"expiry_end_time"`
		IsForbid        int    `form:"is_forbid"`
		helper.MSPNormalizeListPaginateParams
	}{
		MSPNormalizeListPaginateParams: helper.DefaultMSPNormalizeListPaginateParams,
	}
	c.BindG(&req)

	query := core.Is.Db.Session(&gorm.Session{}).
		Select("common_name").
		Where(`common_name != ""`).
		Group("common_name")

	var uniqueIds []string
	var totalNum int64

	if req.UniqueID != "" {
		uniqueIds = []string{req.UniqueID}
		totalNum = 1
	}

	if req.Role != "" {
		query = query.Where("ca_label = ?", strings.ToLower(req.Role))
	}

	var expiryStartDate, expiryEndDate *time.Time

	if req.ExpiryStartTime != "" {
		date, err := dateparse.ParseAny(req.ExpiryStartTime)
		if err != nil {
			return nil, errors.Wrap(err, "Expiration start time error")
		}
		query = query.Where("expiry > ?", date)
		expiryStartDate = &date
	}

	if req.ExpiryEndTime != "" {
		date, err := dateparse.ParseAny(req.ExpiryEndTime)
		if err != nil {
			return nil, errors.Wrap(err, "Expiration end time error")
		}
		query = query.Where("expiry < ?", date)
		expiryEndDate = &date
	}

	switch {
	case req.IsForbid == 1:
		query = query.Where("reason = ?", 2)
	case req.IsForbid == 2:
		query = query.Where("reason = ?", 0)
	default:
		query = query.Where("reason IN ?", []int{0, 2})
	}

	if len(uniqueIds) == 0 {
		list, total, err := dao.GetAllCertificates(query, req.Page, req.LimitNum, "common_name asc")
		if err != nil {
			a.logger.Errorf("Database query error: %s", err)
			return nil, err
		}

		totalNum = total

		for _, row := range list {
			if row.CommonName.Valid {
				uniqueIds = append(uniqueIds, row.CommonName.String)
			}
		}
	}

	query = core.Is.Db.Session(&gorm.Session{}).
		Where("common_name IN ?", uniqueIds)

	if expiryStartDate != nil {
		query = query.Where("expiry > ?", *expiryStartDate)
	}

	if expiryEndDate != nil {
		query = query.Where("expiry < ?", *expiryEndDate)
	}

	switch {
	case req.IsForbid == 1:
		query = query.Where("reason = ?", 2)
	case req.IsForbid == 2:
		query = query.Where("reason = ?", 0)
	default:
		query = query.Where("reason IN ?", []int{0, 2})
	}

	if req.Role != "" {
		query = query.Where("ca_label = ?", strings.ToLower(req.Role))
	}

	list, _, err := dao.GetAllCertificates(query, 1, 100, "issued_at desc")
	if err != nil {
		a.logger.Errorf("Database query error: %s", err)
		return nil, err
	}

	a.logger.Debugf("Number of returned certificates: %v", len(list))

	forbidMap, err := a.logic.UnitsForbidQuery(&logic.UnitsForbidQueryParams{
		UniqueIds: uniqueIds,
	})
	if err != nil {
		a.logger.Errorf("Service prohibition status query error: %s", err)
		return nil, err
	}

	unitsCertsMap := make(map[string]*UnitsCertsItem)
	for _, row := range list {
		uid := row.CommonName.String
		if _, ok := unitsCertsMap[uid]; !ok {
			unitsCertsMap[uid] = &UnitsCertsItem{
				UniqueID:  uid,
				Forbidden: forbidMap.Status[uid].Forbid,
			}
		}

		fullCert, err := schema.GetFullCertByModelCert(row)
		if err != nil {
			a.logger.Errorf("Get full cert error: %s", err)
			continue
		}
		unitsCertsMap[uid].Certs = append(unitsCertsMap[uid].Certs, fullCert)
	}

	var result []*UnitsCertsItem
	for _, v := range unitsCertsMap {
		result = append(result, v)
	}

	a.logger.Debugf("Return service quantity: %v", len(result))

	return helper.MSPNormalizeList{
		List: result,
		Paginate: helper.MSPNormalizePaginate{
			Total:    totalNum,
			Current:  req.Page,
			PageSize: req.LimitNum,
		},
	}, nil
}

//func getExpiryCountByDuration(sign string) (before, after time.Time, err error) {
//	expiryDate := func(du time.Duration) time.Time {
//		return time.Now().Add(du)
//	}
//
//	switch sign {
//	case "1w":
//		before = time.Now()
//		after = expiryDate(7 * 24 * time.Hour)
//	case "1m":
//		before = time.Now()
//		after = expiryDate(30 * 24 * time.Hour)
//	case "3m":
//		before = time.Now()
//		after = expiryDate(3 * 30 * 24 * time.Hour)
//	case "3m+":
//		before = expiryDate(3 * 30 * 24 * time.Hour)
//		after = expiryDate(999 * 30 * 24 * time.Hour)
//	default:
//		return time.Time{}, time.Time{}, err
//	}
//
//	return
//}
