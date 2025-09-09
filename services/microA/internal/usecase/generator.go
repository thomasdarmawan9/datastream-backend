package usecase

import (
	"math/rand"
	"time"

	"github.com/thomasdarmawan9/datastream-backend/services/microA/internal/domain"
)

type SensorGenerator struct {
	sensorTypes []string
	freq        time.Duration
}

func NewSensorGenerator(freq time.Duration) *SensorGenerator {
	return &SensorGenerator{
		sensorTypes: []string{"temperature", "humidity", "co2", "motion", "light", "noise"},
		freq:        freq,
	}
}

func (g *SensorGenerator) Generate() <-chan *domain.SensorData {
	ch := make(chan *domain.SensorData)

	go func() {
		defer close(ch)
		rand.Seed(time.Now().UnixNano())

		for {
			now := time.Now()

			// --- Temperature (18–30 °C) ---
			ch <- &domain.SensorData{
				SensorValue: 18 + rand.Float64()*12,
				SensorType:  "temperature",
				ID1:         "room-101",
				ID2:         rand.Intn(3) + 1,
				TS:          now,
				CreatedAt:   now,
			}

			// --- Humidity (30–70 %) ---
			ch <- &domain.SensorData{
				SensorValue: 30 + rand.Float64()*40,
				SensorType:  "humidity",
				ID1:         "room-102",
				ID2:         rand.Intn(2) + 1,
				TS:          now,
				CreatedAt:   now,
			}

			// --- CO2 (400–2000 ppm) ---
			ch <- &domain.SensorData{
				SensorValue: 400 + rand.Float64()*1600,
				SensorType:  "co2",
				ID1:         "meeting-room-A",
				ID2:         1,
				TS:          now,
				CreatedAt:   now,
			}

			// --- Motion (0/1) ---
			ch <- &domain.SensorData{
				SensorValue: float64(rand.Intn(2)),
				SensorType:  "motion",
				ID1:         "corridor-1",
				ID2:         rand.Intn(5) + 1,
				TS:          now,
				CreatedAt:   now,
			}

			// --- Light (100–1000 lux) ---
			ch <- &domain.SensorData{
				SensorValue: 100 + rand.Float64()*900,
				SensorType:  "light",
				ID1:         "room-101",
				ID2:         rand.Intn(2) + 1,
				TS:          now,
				CreatedAt:   now,
			}

			// --- Noise (30–100 dB) ---
			ch <- &domain.SensorData{
				SensorValue: 30 + rand.Float64()*70,
				SensorType:  "noise",
				ID1:         "cafeteria",
				ID2:         1,
				TS:          now,
				CreatedAt:   now,
			}

			// tunggu sesuai interval freq
			time.Sleep(g.freq)
		}
	}()

	return ch
}
