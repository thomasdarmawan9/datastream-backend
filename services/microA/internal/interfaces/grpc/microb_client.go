package grpc

import (
	"context"
	"log"
	"time"

	"github.com/thomasdarmawan9/datastream-backend/proto/sensorpb"
	"google.golang.org/grpc"
)

type MicroBClient struct {
	client sensorpb.SensorServiceClient
	conn   *grpc.ClientConn
}

func NewMicroBClient(address string) (*MicroBClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := sensorpb.NewSensorServiceClient(conn)
	return &MicroBClient{client: client, conn: conn}, nil
}

func (c *MicroBClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *MicroBClient) StreamSensorData(data []*sensorpb.SensorData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := c.client.StreamData(ctx)
	if err != nil {
		return err
	}

	for _, d := range data {
		req := &sensorpb.StreamRequest{Data: d}
		if err := stream.Send(req); err != nil {
			log.Printf("failed to send data: %v", err)
			return err
		}
		log.Printf("sent data: %+v", d)
	}

	// Close stream and get response
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("MicroB response: %s - %s", res.Status, res.Message)
	return nil
}
