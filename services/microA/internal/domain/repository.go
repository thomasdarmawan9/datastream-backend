package domain

type DeviceRepository interface {
	Create(device *Device) error
	FindByID(id int64) (*Device, error)
	FindAll() ([]Device, error)
}

type RawSensorDataRepository interface {
	Insert(data *RawSensorData) error
	FindByDevice(deviceID int64) ([]RawSensorData, error)
	FindLatest(deviceID int64, sensorType string) (*RawSensorData, error)
}
