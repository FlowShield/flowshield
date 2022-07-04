package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

/*
DB Table Details
-------------------------------------


CREATE TABLE `self_keypair` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL,
  `private_key` text,
  `certificate` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  UNIQUE KEY `id` (`id`),
  KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

JSON Sample
-------------------------------------
{    "certificate": "WvVWjyOFZykWjUNNiChBewKis",    "created_at": "2235-11-29T17:34:03.437242436+08:00",    "updated_at": "2081-08-11T08:12:20.820976918+08:00",    "id": 42,    "name": "sdagyhnWyjVlIZVposUXUiwHI",    "private_key": "BRAIisbANbBPGCpqNcgUbjnIV"}


Comments
-------------------------------------
[ 0] column is set for unsignedWarning table: self_keypair does not have a primary key defined, setting col position 1 id as primary key




*/

// SelfKeypair struct is a row record of the self_keypair table in the cap database
type SelfKeypair struct {
	//[ 0] id                                             uint                 null: false  primary: true   isArray: false  auto: true   col: uint            len: -1      default: []
	ID uint32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:uint;" json:"id" db:"id"`
	//[ 1] name                                           varchar(40)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 40      default: []
	Name string `gorm:"column:name;type:varchar;size:40;" json:"name" db:"name"`
	//[ 2] private_key                                    text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	PrivateKey sql.NullString `gorm:"column:private_key;type:text;size:65535;" json:"private_key" db:"private_key"`
	//[ 3] certificate                                    text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Certificate sql.NullString `gorm:"column:certificate;type:text;size:65535;" json:"certificate" db:"certificate"`
	//[ 4] created_at                                     timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;" json:"created_at" db:"created_at"`
	//[ 5] updated_at                                     timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;" json:"updated_at" db:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (s *SelfKeypair) TableName() string {
	return "self_keypair"
}
