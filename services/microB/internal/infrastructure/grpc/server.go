package grpc

import (
	"io"
	"log"
	"sync"
	"time"

	sensorpb "github.com/thomasdarmawan9/datastream-backend/proto/sensorpb"

	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/domain"
)

type SensorGRPCServer struct {
	sensorpb.UnimplementedSensorServiceServer
	sensorRepo domain.SensorRepository
}

func NewSensorGRPCServer(repo domain.SensorRepository) *SensorGRPCServer {
	return &SensorGRPCServer{sensorRepo: repo}
}

// StreamData menerima stream dari MicroA
func (s *SensorGRPCServer) StreamData(stream sensorpb.SensorService_StreamDataServer) error {
	var (
		mu      sync.Mutex
		sensors []*domain.SensorData
	)

	// ticker buat auto flush tiap 5 detik
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// goroutine buat auto flush
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				if len(sensors) > 0 {
					log.Printf("Auto flushing %d records...", len(sensors))
					if err := s.sensorRepo.StoreBatch(sensors); err != nil {
						log.Printf("Error storing batch: %v", err)
					} else {
						log.Printf("Successfully flushed %d records", len(sensors))
						sensors = nil // kosongkan buffer
					}
				}
				mu.Unlock()
			case <-done:
				return
			}
		}
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// flush terakhir
			mu.Lock()
			if len(sensors) > 0 {
				if err := s.sensorRepo.StoreBatch(sensors); err != nil {
					log.Printf("Error storing batch on EOF: %v", err)
					mu.Unlock()
					close(done)
					return stream.SendAndClose(&sensorpb.StreamResponse{
						Status:  "error",
						Message: err.Error(),
					})
				}
				log.Printf("Successfully flushed %d records on EOF", len(sensors))
			}
			mu.Unlock()

			close(done)
			return stream.SendAndClose(&sensorpb.StreamResponse{
				Status:  "ok",
				Message: "data stored",
			})
		}
		if err != nil {
			close(done)
			return err
		}

		data := req.GetData()

		// parse timestamp
		t, err := time.Parse(time.RFC3339, data.Timestamp)
		if err != nil {
			log.Printf("Invalid timestamp: %v, using now()", err)
			t = time.Now()
		}

		mu.Lock()
		sensors = append(sensors, &domain.SensorData{
			SensorValue: data.SensorValue,
			SensorType:  data.SensorType,
			ID1:         data.Id1,
			ID2:         int(data.Id2),
			TS:          t,
			CreatedAt:   time.Now(),
		})
		mu.Unlock()

		log.Printf("Received data: %+v", data)
	}
}
