package workload

import (
	"github.com/cloudslit/cloudslit/ca/database/mysql/cfssl-model/dao"
)

type UnitsForbidQueryParams struct {
	UniqueIds []string `json:"unique_ids"`
}

type UnitForbidQueryItem struct {
	UniqueId string `json:"unique_id"`
	Forbid   bool   `json:"forbid"`
}

type UnitsForbidQueryResult struct {
	Status map[string]UnitForbidQueryItem `json:"status"`
}

// UnitsForbidQuery Query unique_id Is it forbidden to apply for certificate
func (l *Logic) UnitsForbidQuery(params *UnitsForbidQueryParams) (*UnitsForbidQueryResult, error) {
	db := l.db.Where("unique_id IN ?", params.UniqueIds).
		Where("deleted_at IS NULL")
	list, _, err := dao.GetAllForbid(db, 1, 1000, "id desc")
	if err != nil {
		l.logger.Errorf("Database query error: %s", err)
		return nil, err
	}
	result := UnitsForbidQueryResult{
		Status: make(map[string]UnitForbidQueryItem),
	}

	l.logger.Debugf("Query results: %v", list)

	for _, uid := range params.UniqueIds {
		result.Status[uid] = UnitForbidQueryItem{
			UniqueId: uid,
			Forbid:   false,
		}
	}

	for _, row := range list {
		result.Status[row.UniqueID] = UnitForbidQueryItem{
			UniqueId: row.UniqueID,
			Forbid:   true,
		}
	}

	return &result, nil
}
