package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	Id         uuid.UUID    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt  time.Time    `json:"created_at" gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `json:"-" gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `json:"-" gorm:"type:TIMESTAMP with time zone;null"`
}
