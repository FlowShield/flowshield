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


CREATE TABLE `ocsp_responses` (
  `serial_number` varchar(128) NOT NULL,
  `authority_key_identifier` varchar(128) NOT NULL,
  `body` text NOT NULL,
  `expiry` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`serial_number`,`authority_key_identifier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

JSON Sample
-------------------------------------
{    "authority_key_identifier": "uccWRqKiBWSiUMLUVFMPDTagi",    "body": "HgsoAwIJdyGtGlyHKIlydlsfZ",    "expiry": "2038-07-05T20:42:34.735852468+08:00",    "serial_number": "BVkGevecZyAogIjcXevDgWVHj"}



*/

// OcspResponses struct is a row record of the ocsp_responses table in the cap database
type OcspResponses struct {
	//[ 0] serial_number                                  varchar(128)         null: false  primary: true   isArray: false  auto: false  col: varchar         len: 128     default: []
	SerialNumber string `gorm:"primary_key;column:serial_number;type:varchar;size:128;" json:"serial_number" db:"serial_number"`
	//[ 1] authority_key_identifier                       varchar(128)         null: false  primary: true   isArray: false  auto: false  col: varchar         len: 128     default: []
	AuthorityKeyIdentifier string `gorm:"primary_key;column:authority_key_identifier;type:varchar;size:128;" json:"authority_key_identifier" db:"authority_key_identifier"`
	//[ 2] body                                           text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Body string `gorm:"column:body;type:text;size:65535;" json:"body" db:"body"`
	//[ 3] expiry                                         timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	Expiry time.Time `gorm:"column:expiry;type:timestamp;" json:"expiry" db:"expiry"`
}

var ocsp_responsesTableInfo = &TableInfo{
	Name: "ocsp_responses",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "serial_number",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(128)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       128,
			GoFieldName:        "SerialNumber",
			GoFieldType:        "string",
			JSONFieldName:      "serial_number",
			ProtobufFieldName:  "serial_number",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "authority_key_identifier",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(128)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       128,
			GoFieldName:        "AuthorityKeyIdentifier",
			GoFieldType:        "string",
			JSONFieldName:      "authority_key_identifier",
			ProtobufFieldName:  "authority_key_identifier",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "body",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(65535)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       65535,
			GoFieldName:        "Body",
			GoFieldType:        "string",
			JSONFieldName:      "body",
			ProtobufFieldName:  "body",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "expiry",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "timestamp",
			DatabaseTypePretty: "timestamp",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "timestamp",
			ColumnLength:       -1,
			GoFieldName:        "Expiry",
			GoFieldType:        "time.Time",
			JSONFieldName:      "expiry",
			ProtobufFieldName:  "expiry",
			ProtobufType:       "uint64",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (o *OcspResponses) TableName() string {
	return "ocsp_responses"
}
