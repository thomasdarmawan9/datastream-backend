package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thomasdarmawan9/datastream-backend/services/microA/internal/usecase"
)

type SensorHandler struct {
	sensorUC usecase.RawSensorDataUsecase
}

func NewSensorHandler(e *echo.Echo, uc usecase.RawSensorDataUsecase) {
	handler := &SensorHandler{sensorUC: uc}

	g := e.Group("/sensors")
	g.POST("", handler.InsertSensorData)
	g.GET("/:device_id", handler.ListSensorData)
	g.GET("/:device_id/latest/:sensor_type", handler.GetLatestSensorData)
}

func (h *SensorHandler) InsertSensorData(c echo.Context) error {
	var req struct {
		DeviceID   int64   `json:"device_id"`
		SensorType string  `json:"sensor_type"`
		Value      float64 `json:"value"`
		Timestamp  string  `json:"timestamp"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ts := time.Now()
	if req.Timestamp != "" {
		parsed, err := time.Parse(time.RFC3339, req.Timestamp)
		if err == nil {
			ts = parsed
		}
	}

	data, err := h.sensorUC.InsertSensorData(req.DeviceID, req.SensorType, req.Value, ts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func (h *SensorHandler) ListSensorData(c echo.Context) error {
	deviceID, _ := strconv.ParseInt(c.Param("device_id"), 10, 64)
	data, err := h.sensorUC.GetAllSensorData(deviceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func (h *SensorHandler) GetLatestSensorData(c echo.Context) error {
	deviceID, _ := strconv.ParseInt(c.Param("device_id"), 10, 64)
	sensorType := c.Param("sensor_type")

	data, err := h.sensorUC.GetLatestSensorData(deviceID, sensorType)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, data)
}
