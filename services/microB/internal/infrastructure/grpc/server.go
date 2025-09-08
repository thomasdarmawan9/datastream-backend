package grpc

import (
	"io"
	"log"

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
	var sensors []*domain.SensorData

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Simpan batch ke DB ketika stream selesai
			if len(sensors) > 0 {
				if err := s.sensorRepo.StoreBatch(sensors); err != nil {
					return stream.SendAndClose(&sensorpb.StreamResponse{
						Status:  "error",
						Message: err.Error(),
					})
				}
			}
			return stream.SendAndClose(&sensorpb.StreamResponse{
				Status:  "ok",
				Message: "data stored",
			})
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		sensors = append(sensors, &domain.SensorData{
			SensorValue: data.SensorValue,
			SensorType:  data.SensorType,
			ID1:         data.Id1,
			ID2:         int(data.Id2),
			// di sini timestamp string harus di-parse (misalnya pakai time.Parse RFC3339)
			// sementara anggap valid
			// TS: ...
		})
		log.Printf("Received data: %+v", data)
	}
}
