/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// GormJSON defined JSON data type, need to implement driver.Valuer, sql.Scanner interface
type GormJSON json.RawMessage

// Value return json value, implement driver.Valuer interface
func (j GormJSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	bytes, err := json.RawMessage(j).MarshalJSON()
	return string(bytes), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *GormJSON) Scan(value any) error {
	if value == nil {
		*j = GormJSON("null")
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = GormJSON(result)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (j GormJSON) MarshalJSON() ([]byte, error) {
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON to deserialize []byte
func (j *GormJSON) UnmarshalJSON(b []byte) error {
	result := json.RawMessage{}
	err := result.UnmarshalJSON(b)
	*j = GormJSON(result)
	return err
}

func (j GormJSON) String() string {
	return string(j)
}

// GormDataType gorm common data type
func (j GormJSON) GormDataType() string {
	return "json"
}

// GormDBDataType gorm dao data type
func (j GormJSON) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (j GormJSON) GormValue(_ context.Context, db *gorm.DB) clause.Expr {
	if len(j) == 0 {
		return gorm.Expr("NULL")
	}

	data, _ := j.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if _, ok := db.Dialector.(*mysql.Dialector); ok {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}
