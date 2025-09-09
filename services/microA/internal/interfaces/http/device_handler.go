package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/thomasdarmawan9/datastream-backend/services/microA/internal/usecase"
)

type DeviceHandler struct {
	deviceUC usecase.DeviceUsecase
}

func NewDeviceHandler(e *echo.Echo, uc usecase.DeviceUsecase) {
	handler := &DeviceHandler{deviceUC: uc}

	g := e.Group("/devices")
	g.POST("", handler.RegisterDevice)
	g.GET("", handler.ListDevices)
	g.GET("/:id", handler.GetDevice)
}

// RegisterDevice godoc
// @Summary Register new device
// @Tags Devices
// @Accept json
// @Produce json
// @Param body body map[string]string true "Device payload"
// @Success 200 {object} map[string]interface{}
// @Router /devices [post]
func (h *DeviceHandler) RegisterDevice(c echo.Context) error {
	var req struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	device, err := h.deviceUC.RegisterDevice(req.Name, req.Location)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, device)
}

func (h *DeviceHandler) ListDevices(c echo.Context) error {
	devices, err := h.deviceUC.ListDevices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) GetDevice(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	device, err := h.deviceUC.GetDeviceByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "device not found"})
	}
	return c.JSON(http.StatusOK, device)
}
