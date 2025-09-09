package domain

import "time"

type RawSensorData struct {
	ID         int64     `gorm:"primaryKey"`
	DeviceID   int64     `gorm:"index;not null"`
	SensorType string    `gorm:"size:64;not null;index"`
	Value      float64   `gorm:"not null"`
	Timestamp  time.Time `gorm:"not null;index"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type SensorData struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	SensorValue float64    `json:"sensor_value" gorm:"not null"`
	SensorType  string     `json:"sensor_type" gorm:"size:64;not null"`
	ID1         string     `json:"id1" gorm:"size:128;not null"`
	ID2         int        `json:"id2" gorm:"not null"`
	TS          time.Time  `json:"timestamp" gorm:"not null;index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}
