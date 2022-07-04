package dao

import (
	"time"

	"github.com/guregu/null"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/cloudSlit/cloudslit/ca/database/mysql/cfssl-model/model"
)

var (
	_ = time.Second
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// GetAllForbid is a function to get a slice of record(s) from forbid table in the cap database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func GetAllForbid(db *gorm.DB, page, pagesize int, order string) (results []*model.Forbid, totalRows int64, err error) {

	resultOrm := db.Model(&model.Forbid{})
	resultOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
	}

	if order != "" {
		resultOrm = resultOrm.Order(order)
	}

	if err = resultOrm.Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return results, totalRows, nil
		}
		return nil, -1, err
	}

	return results, totalRows, nil
}

// GetForbid is a function to get a single record from the forbid table in the cap database
// error - ErrNotFound, db Find error
func GetForbid(db *gorm.DB) (record *model.Forbid, err error) {
	record = &model.Forbid{}
	if err = db.First(record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return record, err
	}

	return record, nil
}

// AddForbid is a function to add a single record to forbid table in the cap database
// error - ErrInsertFailed, db save call failed
func AddForbid(db *gorm.DB, record *model.Forbid) (result *model.Forbid, RowsAffected int64, err error) {
	query := db.Save(record)
	if err = query.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}

	return record, query.RowsAffected, nil
}

// UpdateForbid is a function to update a single record from forbid table in the cap database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func UpdateForbid(argId uint32, updated *model.Forbid) (result *model.Forbid, RowsAffected int64, err error) {

	result = &model.Forbid{}
	db := DB.First(result, argId)
	if err = db.Error; err != nil {
		return nil, -1, ErrNotFound
	}

	if err = Copy(result, updated); err != nil {
		return nil, -1, ErrUpdateFailed
	}

	db = db.Save(result)
	if err = db.Error; err != nil {
		return nil, -1, ErrUpdateFailed
	}

	return result, db.RowsAffected, nil
}

// DeleteForbid is a function to delete a single record from forbid table in the cap database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func DeleteForbid(argId uint32) (rowsAffected int64, err error) {

	record := &model.Forbid{}
	db := DB.First(record, argId)
	if db.Error != nil {
		return -1, ErrNotFound
	}

	db = db.Delete(record)
	if err = db.Error; err != nil {
		return -1, ErrDeleteFailed
	}

	return db.RowsAffected, nil
}
