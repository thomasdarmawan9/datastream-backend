package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/usecase"
)

type SensorHandler struct {
	usecase usecase.SensorUsecase
}

func NewSensorHandler(g *echo.Group, uc usecase.SensorUsecase) {
	handler := &SensorHandler{usecase: uc}

	g.GET("/sensors", handler.GetByFilter)       // GET /api/sensors
	g.PUT("/sensors", handler.UpdateByFilter)    // PUT /api/sensors
	g.DELETE("/sensors", handler.DeleteByFilter) // DELETE /api/sensors
}

// GetByFilter godoc
// @Summary Get sensor data by filter
// @Description Retrieve sensor data based on various filters
// @Tags sensors
// @Accept json
// @Produce json
// @Param id1 query string false "ID1 filter"
// @Param id2 query int false "ID2 filter"
// @Param from query string false "Start timestamp (RFC3339)"
// @Param to query string false "End timestamp (RFC3339)"
// @Param limit query int false "Limit number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sensors [get]
func (h *SensorHandler) GetByFilter(c echo.Context) error {
	id1 := c.QueryParam("id1")
	id2Str := c.QueryParam("id2")
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	var id2 *int
	if id2Str != "" {
		val, err := strconv.Atoi(id2Str)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id2"})
		}
		id2 = &val
	}

	var from, to *time.Time
	if fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &t
		}
	}
	if toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}

	limit := 10
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}
	offset := 0
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil {
			offset = v
		}
	}

	data, total, err := h.usecase.GetByFilter(id1, id2, from, to, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total": total,
		"data":  data,
	})
}

// UpdateByFilter godoc
// @Summary Update sensor data by filter
// @Description Update sensor data values based on various filters
// @Tags sensors
// @Accept json
// @Produce json
// @Param id1 query string false "ID1 filter"
// @Param id2 query int false "ID2 filter"
// @Param from query string false "Start timestamp (RFC3339)"
// @Param to query string false "End timestamp (RFC3339)"
// @Param new_value query number true "New sensor value"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sensors [put]
func (h *SensorHandler) UpdateByFilter(c echo.Context) error {
	id1 := c.QueryParam("id1")
	id2Str := c.QueryParam("id2")
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")
	newValueStr := c.QueryParam("new_value")

	if newValueStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "new_value required"})
	}
	newValue, err := strconv.ParseFloat(newValueStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid new_value"})
	}

	var id2 *int
	if id2Str != "" {
		val, err := strconv.Atoi(id2Str)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id2"})
		}
		id2 = &val
	}

	var from, to *time.Time
	if fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &t
		}
	}
	if toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}

	updated, err := h.usecase.UpdateByFilter(id1, id2, from, to, newValue)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"updated": updated,
	})
}

// DeleteByFilter godoc
// @Summary Delete sensor data by filter
// @Description Delete sensor data based on various filters
// @Tags sensors
// @Accept json
// @Produce json
// @Param id1 query string false "ID1 filter"
// @Param id2 query int false "ID2 filter"
// @Param from query string false "Start timestamp (RFC3339)"
// @Param to query string false "End timestamp (RFC3339)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /sensors [delete]
func (h *SensorHandler) DeleteByFilter(c echo.Context) error {
	id1 := c.QueryParam("id1")
	id2Str := c.QueryParam("id2")
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")

	var id2 *int
	if id2Str != "" {
		val, err := strconv.Atoi(id2Str)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id2"})
		}
		id2 = &val
	}

	var from, to *time.Time
	if fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &t
		}
	}
	if toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &t
		}
	}

	deleted, err := h.usecase.DeleteByFilter(id1, id2, from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"deleted": deleted,
	})
}
