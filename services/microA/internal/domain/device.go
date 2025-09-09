package domain

import "time"

type Device struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"size:128;not null;unique"`
	Location  string    `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

