package main

import (
    "fmt"
    "net/http"
    "os"

    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    e.GET("/health", func(c echo.Context) error {
        return c.String(http.StatusOK, "microA ok")
    })

    e.GET("/api/v1/config", func(c echo.Context) error {
        // Placeholder: return current frequency config
        return c.JSON(http.StatusOK, map[string]interface{}{"frequency_hz": 1})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8090"
    }
    addr := fmt.Sprintf(":%s", port)
    e.Logger.Fatal(e.Start(addr))
}
