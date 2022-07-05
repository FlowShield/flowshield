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


CREATE TABLE `certificates` (
  `serial_number` varchar(128) NOT NULL,
  `authority_key_identifier` varchar(128) NOT NULL,
  `ca_label` varchar(128) DEFAULT NULL,
  `status` varchar(128) NOT NULL,
  `reason` int(11) DEFAULT NULL,
  `expiry` timestamp NULL DEFAULT NULL,
  `revoked_at` timestamp NULL DEFAULT NULL,
  `pem` text NOT NULL,
  `issued_at` timestamp NULL DEFAULT NULL,
  `not_before` timestamp NULL DEFAULT NULL,
  `metadata` json DEFAULT NULL,
  `sans` json DEFAULT NULL,
  `common_name` text,
  PRIMARY KEY (`serial_number`,`authority_key_identifier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

JSON Sample
-------------------------------------
{    "common_name": "MScvMbPGahOUZxoOmsDxmNYFD",    "ca_label": "ULPoUZlfadBgOnGmGeTYpYlwr",    "pem": "fsrUtxDJCOpQbRZspCyJcqgNf",    "not_before": "2223-09-05T23:15:27.424784706+08:00",    "metadata": "swThVcEqvHBnUptLkmpCDXCQU",    "sans": "wVWFshOhaxDVLKoepAfWDCIyK",    "revoked_at": "2148-12-30T01:02:53.517625773+08:00",    "issued_at": "2022-06-24T10:47:05.902911131+08:00",    "serial_number": "IpDFkiGNIMlxfMZCDtIxJJdUt",    "authority_key_identifier": "VWhZXqDQGfeWADHfPHZUBcWJD",    "status": "jvpyBkTOqxxojTdhHuywdLWIH",    "reason": 32,    "expiry": "2299-08-23T23:34:03.786531333+08:00"}



*/

// Certificates struct is a row record of the certificates table in the cap database
type Certificates struct {
	//[ 0] serial_number                                  varchar(128)         null: false  primary: true   isArray: false  auto: false  col: varchar         len: 128     default: []
	SerialNumber string `gorm:"primary_key;column:serial_number;type:varchar;size:128;" json:"serial_number" db:"serial_number"`
	//[ 1] authority_key_identifier                       varchar(128)         null: false  primary: true   isArray: false  auto: false  col: varchar         len: 128     default: []
	AuthorityKeyIdentifier string `gorm:"primary_key;column:authority_key_identifier;type:varchar;size:128;" json:"authority_key_identifier" db:"authority_key_identifier"`
	//[ 2] ca_label                                       varchar(128)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 128     default: []
	CaLabel sql.NullString `gorm:"column:ca_label;type:varchar;size:128;" json:"ca_label" db:"ca_label"`
	//[ 3] status                                         varchar(128)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 128     default: []
	Status string `gorm:"column:status;type:varchar;size:128;" json:"status" db:"status"`
	//[ 4] reason                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Reason sql.NullInt64 `gorm:"column:reason;type:int;" json:"reason" db:"reason"`
	//[ 5] expiry                                         timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	Expiry time.Time `gorm:"column:expiry;type:timestamp;" json:"expiry" db:"expiry"`
	//[ 6] revoked_at                                     timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	RevokedAt time.Time `gorm:"column:revoked_at;type:timestamp;" json:"revoked_at" db:"revoked_at"`
	//[ 7] pem                                            text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Pem string `gorm:"column:pem;type:text;size:65535;" json:"pem" db:"pem"`
	//[ 8] issued_at                                      timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	IssuedAt time.Time `gorm:"column:issued_at;type:timestamp;" json:"issued_at" db:"issued_at"`
	//[ 9] not_before                                     timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	NotBefore time.Time `gorm:"column:not_before;type:timestamp;" json:"not_before" db:"not_before"`
	//[10] metadata                                       json                 null: true   primary: false  isArray: false  auto: false  col: json            len: -1      default: []
	Metadata sql.NullString `gorm:"column:metadata;type:json;" json:"metadata" db:"metadata"`
	//[11] sans                                           json                 null: true   primary: false  isArray: false  auto: false  col: json            len: -1      default: []
	Sans sql.NullString `gorm:"column:sans;type:json;" json:"sans" db:"sans"`
	//[12] common_name                                    text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	CommonName sql.NullString `gorm:"column:common_name;type:text;size:65535;" json:"common_name" db:"common_name"`
}

var certificatesTableInfo = &TableInfo{
	Name: "certificates",
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
			Name:               "ca_label",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       128,
			GoFieldName:        "CaLabel",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "ca_label",
			ProtobufFieldName:  "ca_label",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "status",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       128,
			GoFieldName:        "Status",
			GoFieldType:        "string",
			JSONFieldName:      "status",
			ProtobufFieldName:  "status",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "reason",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Reason",
			GoFieldType:        "sql.NullInt64",
			JSONFieldName:      "reason",
			ProtobufFieldName:  "reason",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
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
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "revoked_at",
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
			GoFieldName:        "RevokedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "revoked_at",
			ProtobufFieldName:  "revoked_at",
			ProtobufType:       "uint64",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "pem",
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
			GoFieldName:        "Pem",
			GoFieldType:        "string",
			JSONFieldName:      "pem",
			ProtobufFieldName:  "pem",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "issued_at",
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
			GoFieldName:        "IssuedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "issued_at",
			ProtobufFieldName:  "issued_at",
			ProtobufType:       "uint64",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "not_before",
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
			GoFieldName:        "NotBefore",
			GoFieldType:        "time.Time",
			JSONFieldName:      "not_before",
			ProtobufFieldName:  "not_before",
			ProtobufType:       "uint64",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "metadata",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "json",
			DatabaseTypePretty: "json",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "json",
			ColumnLength:       -1,
			GoFieldName:        "Metadata",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "metadata",
			ProtobufFieldName:  "metadata",
			ProtobufType:       "string",
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "sans",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "json",
			DatabaseTypePretty: "json",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "json",
			ColumnLength:       -1,
			GoFieldName:        "Sans",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "sans",
			ProtobufFieldName:  "sans",
			ProtobufType:       "string",
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "common_name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(65535)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       65535,
			GoFieldName:        "CommonName",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "common_name",
			ProtobufFieldName:  "common_name",
			ProtobufType:       "string",
			ProtobufPos:        13,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *Certificates) TableName() string {
	return "certificates"
}
