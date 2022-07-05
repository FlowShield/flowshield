package dao

import (
	"context"
	"time"

	"github.com/guregu/null"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/cloudslit/cloudslit/ca/database/mysql/cfssl-model/model"
)

var (
	_ = time.Second
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// GetAllSelfKeypair is a function to get a slice of record(s) from self_keypair table in the cap database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func GetAllSelfKeypair(db *gorm.DB, page, pagesize int, order string) (results []*model.SelfKeypair, totalRows int64, err error) {

	resultOrm := db.Model(&model.SelfKeypair{})
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
		err = ErrNotFound
		return nil, -1, err
	}

	return results, totalRows, nil
}

// GetSelfKeypair is a function to get a single record from the self_keypair table in the cap database
// error - ErrNotFound, db Find error
func GetSelfKeypair(db *gorm.DB, argId uint32) (record *model.SelfKeypair, err error) {
	record = &model.SelfKeypair{}
	if err = db.Where("id = ?", argId).First(record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return record, nil
}

// AddSelfKeypair is a function to add a single record to self_keypair table in the cap database
// error - ErrInsertFailed, db save call failed
func AddSelfKeypair(ctx context.Context, record *model.SelfKeypair) (result *model.SelfKeypair, RowsAffected int64, err error) {
	db := DB.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}

	return record, db.RowsAffected, nil
}

// UpdateSelfKeypair is a function to update a single record from self_keypair table in the cap database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func UpdateSelfKeypair(ctx context.Context, argId uint32, updated *model.SelfKeypair) (result *model.SelfKeypair, RowsAffected int64, err error) {

	result = &model.SelfKeypair{}
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

// DeleteSelfKeypair is a function to delete a single record from self_keypair table in the cap database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func DeleteSelfKeypair(ctx context.Context, argId uint32) (rowsAffected int64, err error) {

	record := &model.SelfKeypair{}
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
