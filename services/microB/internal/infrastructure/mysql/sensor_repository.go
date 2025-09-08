package mysql

import (
	"database/sql"
	"time"

	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/domain"
)

type sensorRepo struct {
	db *sql.DB
}

func NewSensorRepository(db *sql.DB) domain.SensorRepository {
	return &sensorRepo{db: db}
}

func (r *sensorRepo) Store(sensor *domain.SensorData) error {
	query := `INSERT INTO sensor_data (sensor_value, sensor_type, id1, id2, ts) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, sensor.SensorValue, sensor.SensorType, sensor.ID1, sensor.ID2, sensor.TS)
	return err
}

func (r *sensorRepo) StoreBatch(sensors []*domain.SensorData) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO sensor_data (sensor_value, sensor_type, id1, id2, ts) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range sensors {
		_, err := stmt.Exec(s.SensorValue, s.SensorType, s.ID1, s.ID2, s.TS)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *sensorRepo) FindByFilter(id1 string, id2 *int, from, to *time.Time, limit, offset int) ([]*domain.SensorData, int, error) {
	query := `SELECT id, sensor_value, sensor_type, id1, id2, ts, created_at, updated_at FROM sensor_data WHERE 1=1`
	args := []interface{}{}

	if id1 != "" {
		query += " AND id1 = ?"
		args = append(args, id1)
	}
	if id2 != nil {
		query += " AND id2 = ?"
		args = append(args, *id2)
	}
	if from != nil {
		query += " AND ts >= ?"
		args = append(args, *from)
	}
	if to != nil {
		query += " AND ts <= ?"
		args = append(args, *to)
	}

	// total count
	countQuery := "SELECT COUNT(*) FROM (" + query + ") as tmp"
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// apply pagination
	query += " ORDER BY ts ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*domain.SensorData
	for rows.Next() {
		var s domain.SensorData
		var updatedAt sql.NullTime
		err := rows.Scan(&s.ID, &s.SensorValue, &s.SensorType, &s.ID1, &s.ID2, &s.TS, &s.CreatedAt, &updatedAt)
		if err != nil {
			return nil, 0, err
		}
		if updatedAt.Valid {
			s.UpdatedAt = &updatedAt.Time
		}
		result = append(result, &s)
	}

	return result, total, nil
}

func (r *sensorRepo) UpdateByFilter(id1 string, id2 *int, from, to *time.Time, newValue float64) (int64, error) {
	query := "UPDATE sensor_data SET sensor_value = ?, updated_at = NOW() WHERE 1=1"
	args := []interface{}{newValue}

	if id1 != "" {
		query += " AND id1 = ?"
		args = append(args, id1)
	}
	if id2 != nil {
		query += " AND id2 = ?"
		args = append(args, *id2)
	}
	if from != nil {
		query += " AND ts >= ?"
		args = append(args, *from)
	}
	if to != nil {
		query += " AND ts <= ?"
		args = append(args, *to)
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *sensorRepo) DeleteByFilter(id1 string, id2 *int, from, to *time.Time) (int64, error) {
	query := "DELETE FROM sensor_data WHERE 1=1"
	args := []interface{}{}

	if id1 != "" {
		query += " AND id1 = ?"
		args = append(args, id1)
	}
	if id2 != nil {
		query += " AND id2 = ?"
		args = append(args, *id2)
	}
	if from != nil {
		query += " AND ts >= ?"
		args = append(args, *from)
	}
	if to != nil {
		query += " AND ts <= ?"
		args = append(args, *to)
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
