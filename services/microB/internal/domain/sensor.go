package domain

import "time"

// SensorData merepresentasikan data dari sensor (entity utama).
type SensorData struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement"`
	SensorValue float64    `gorm:"not null"`
	SensorType  string     `gorm:"type:varchar(64);not null;index"`
	ID1         string     `gorm:"type:char(8);not null;index:idx_ids_ts,priority:1"`
	ID2         int        `gorm:"not null;index:idx_ids_ts,priority:2"`
	TS          time.Time  `gorm:"precision:6;not null;index:idx_ids_ts,priority:3"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime"`
}
