package main

import (
	"context"
	"log"
	"os"
	"time"

	sensorpb "github.com/thomasdarmawan9/datastream-backend/proto/sensorpb"
	"github.com/thomasdarmawan9/datastream-backend/services/microA/internal/usecase"

	"google.golang.org/grpc"
)

func main() {
	// --- Load ENV ---
	microBAddr := os.Getenv("MICROB_GRPC_ADDR") // contoh: "localhost:50051"
	if microBAddr == "" {
		microBAddr = "localhost:50051"
	}

	freqStr := os.Getenv("GEN_FREQ_MS") // default 1000 ms
	freq := time.Millisecond * 1000
	if freqStr != "" {
		if v, err := time.ParseDuration(freqStr + "ms"); err == nil {
			freq = v
		}
	}

	// --- gRPC Dial ke MicroB ---
	conn, err := grpc.Dial(microBAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect MicroB: %v", err)
	}
	defer conn.Close()

	client := sensorpb.NewSensorServiceClient(conn)

	// Open stream
	stream, err := client.StreamData(context.Background())
	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}

	log.Printf("MicroA started. Sending smart-building sensor data every %v → %s", freq, microBAddr)

	// Usecase generator (multi-sensor)
	gen := usecase.NewSensorGenerator(freq)

	// Loop generate & kirim
	for data := range gen.Generate() {
		req := &sensorpb.StreamRequest{
			Data: &sensorpb.SensorData{
				SensorValue: data.SensorValue,
				SensorType:  data.SensorType,
				Id1:         data.ID1,
				Id2:         int32(data.ID2),
				Timestamp:   data.TS.Format(time.RFC3339Nano),
			},
		}

		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send: %v", err)
		}
		log.Printf("Sent: %+v", data)
	}
}
