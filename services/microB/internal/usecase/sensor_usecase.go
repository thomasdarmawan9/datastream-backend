package usecase

import (
	"time"

	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/domain"
)

type SensorUsecase interface {
	Store(sensor *domain.SensorData) error
	StoreBatch(sensors []*domain.SensorData) error
	GetByFilter(id1 string, id2 *int, from, to *time.Time, limit, offset int) ([]*domain.SensorData, int, error)
	UpdateByFilter(id1 string, id2 *int, from, to *time.Time, newValue float64) (int64, error)
	DeleteByFilter(id1 string, id2 *int, from, to *time.Time) (int64, error)
}

type sensorUsecase struct {
	repo domain.SensorRepository
}

func NewSensorUsecase(repo domain.SensorRepository) SensorUsecase {
	return &sensorUsecase{repo: repo}
}

func (u *sensorUsecase) Store(sensor *domain.SensorData) error {
	return u.repo.Store(sensor)
}

func (u *sensorUsecase) StoreBatch(sensors []*domain.SensorData) error {
	return u.repo.StoreBatch(sensors)
}

func (u *sensorUsecase) GetByFilter(id1 string, id2 *int, from, to *time.Time, limit, offset int) ([]*domain.SensorData, int, error) {
	return u.repo.FindByFilter(id1, id2, from, to, limit, offset)
}

func (u *sensorUsecase) UpdateByFilter(id1 string, id2 *int, from, to *time.Time, newValue float64) (int64, error) {
	return u.repo.UpdateByFilter(id1, id2, from, to, newValue)
}

func (u *sensorUsecase) DeleteByFilter(id1 string, id2 *int, from, to *time.Time) (int64, error) {
	return u.repo.DeleteByFilter(id1, id2, from, to)
}
