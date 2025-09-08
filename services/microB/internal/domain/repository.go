package domain

import "time"

// Repository untuk SensorData
type SensorRepository interface {
	Store(sensor *SensorData) error
	StoreBatch(sensors []*SensorData) error
	FindByFilter(id1 string, id2 *int, from, to *time.Time, limit, offset int) ([]*SensorData, int, error)
	UpdateByFilter(id1 string, id2 *int, from, to *time.Time, newValue float64) (int64, error)
	DeleteByFilter(id1 string, id2 *int, from, to *time.Time) (int64, error)
}

// Repository untuk User
type UserRepository interface {
	Create(user *User) error
	FindByUsername(username string) (*User, error)
}
