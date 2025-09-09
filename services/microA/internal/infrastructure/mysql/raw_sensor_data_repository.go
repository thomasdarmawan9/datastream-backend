package mysql

import (
	"database/sql"

	"github.com/thomasdarmawan9/datastream-backend/services/microA/internal/domain"
)

type rawSensorRepo struct {
	db *sql.DB
}

func NewRawSensorDataRepository(db *sql.DB) domain.RawSensorDataRepository {
	return &rawSensorRepo{db: db}
}

func (r *rawSensorRepo) Insert(data *domain.RawSensorData) error {
	query := `INSERT INTO raw_sensor_data (device_id, sensor_type, value, timestamp) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, data.DeviceID, data.SensorType, data.Value, data.Timestamp)
	return err
}

func (r *rawSensorRepo) FindByDevice(deviceID int64) ([]domain.RawSensorData, error) {
	query := `SELECT id, device_id, sensor_type, value, timestamp, created_at 
	          FROM raw_sensor_data WHERE device_id = ? ORDER BY timestamp DESC`
	rows, err := r.db.Query(query, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.RawSensorData
	for rows.Next() {
		var d domain.RawSensorData
		if err := rows.Scan(&d.ID, &d.DeviceID, &d.SensorType, &d.Value, &d.Timestamp, &d.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	return results, nil
}

func (r *rawSensorRepo) FindLatest(deviceID int64, sensorType string) (*domain.RawSensorData, error) {
	query := `SELECT id, device_id, sensor_type, value, timestamp, created_at 
	          FROM raw_sensor_data WHERE device_id = ? AND sensor_type = ? 
	          ORDER BY timestamp DESC LIMIT 1`
	row := r.db.QueryRow(query, deviceID, sensorType)

	var d domain.RawSensorData
	err := row.Scan(&d.ID, &d.DeviceID, &d.SensorType, &d.Value, &d.Timestamp, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
