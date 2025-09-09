package mysql

import (
	"database/sql"

	"github.com/thomasdarmawan9/datastream-backend/services/microA/internal/domain"
)

type deviceRepo struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) domain.DeviceRepository {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) Create(device *domain.Device) error {
	query := `INSERT INTO devices (name, location) VALUES (?, ?)`
	_, err := r.db.Exec(query, device.Name, device.Location)
	return err
}

func (r *deviceRepo) FindByID(id int64) (*domain.Device, error) {
	query := `SELECT id, name, location, created_at FROM devices WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var d domain.Device
	err := row.Scan(&d.ID, &d.Name, &d.Location, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *deviceRepo) FindAll() ([]domain.Device, error) {
	query := `SELECT id, name, location, created_at FROM devices`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []domain.Device
	for rows.Next() {
		var d domain.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.Location, &d.CreatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, nil
}
