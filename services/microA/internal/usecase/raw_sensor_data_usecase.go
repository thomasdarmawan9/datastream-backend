package usecase

import (
	"errors"
	"time"

	"github.com/thomasdarmawan9/datastream-backend/services/microA/internal/domain"
)

type RawSensorDataUsecase interface {
	InsertSensorData(deviceID int64, sensorType string, value float64, ts time.Time) (*domain.RawSensorData, error)
	GetLatestSensorData(deviceID int64, sensorType string) (*domain.RawSensorData, error)
	GetAllSensorData(deviceID int64) ([]domain.RawSensorData, error)
}

type rawSensorDataUsecase struct {
	repo domain.RawSensorDataRepository
}

func NewRawSensorDataUsecase(repo domain.RawSensorDataRepository) RawSensorDataUsecase {
	return &rawSensorDataUsecase{repo: repo}
}

func (uc *rawSensorDataUsecase) InsertSensorData(deviceID int64, sensorType string, value float64, ts time.Time) (*domain.RawSensorData, error) {
	if deviceID == 0 || sensorType == "" {
		return nil, errors.New("invalid sensor data")
	}
	data := &domain.RawSensorData{
		DeviceID:   deviceID,
		SensorType: sensorType,
		Value:      value,
		Timestamp:  ts,
	}
	if err := uc.repo.Insert(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (uc *rawSensorDataUsecase) GetLatestSensorData(deviceID int64, sensorType string) (*domain.RawSensorData, error) {
	return uc.repo.FindLatest(deviceID, sensorType)
}

func (uc *rawSensorDataUsecase) GetAllSensorData(deviceID int64) ([]domain.RawSensorData, error) {
	return uc.repo.FindByDevice(deviceID)
}
