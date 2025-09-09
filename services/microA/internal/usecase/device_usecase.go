package usecase

import (
	"errors"

	"github.com/thomasdarmawan9/datastream-backend/services/microA/internal/domain"
)

type DeviceUsecase interface {
	RegisterDevice(name, location string) (*domain.Device, error)
	GetDeviceByID(id int64) (*domain.Device, error)
	ListDevices() ([]domain.Device, error)
}

type deviceUsecase struct {
	repo domain.DeviceRepository
}

func NewDeviceUsecase(repo domain.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{repo: repo}
}

func (uc *deviceUsecase) RegisterDevice(name, location string) (*domain.Device, error) {
	if name == "" {
		return nil, errors.New("device name cannot be empty")
	}
	device := &domain.Device{
		Name:     name,
		Location: location,
	}
	if err := uc.repo.Create(device); err != nil {
		return nil, err
	}
	return device, nil
}

func (uc *deviceUsecase) GetDeviceByID(id int64) (*domain.Device, error) {
	return uc.repo.FindByID(id)
}

func (uc *deviceUsecase) ListDevices() ([]domain.Device, error) {
	return uc.repo.FindAll()
}
