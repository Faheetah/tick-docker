package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	// Always return a static metric
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Mimic status codes by pushing to the server
	status := 200

	e.GET("/status", func(c echo.Context) error {
		return c.String(status, fmt.Sprintf(""))
	})

	e.GET("/status/:status", func(c echo.Context) error {
		status, _ = strconv.Atoi(c.Param("status"))
		return c.String(http.StatusOK, "OK")
	})

	// Random value
	e.GET("/random", func(c echo.Context) error {
		return c.String(200, fmt.Sprintf("{\"value\": %d}", rand.Intn(100)))
	})

	// Trending
	base := time.Now()
	e.GET("/trending", func(c echo.Context) error {
		d := time.Now().Sub(base)
		t := (d.Seconds() + float64(rand.Intn(5))) / 10
		return c.String(200, fmt.Sprintf("{\"value\": %02f}", t))
	})

	e.GET("/trending/reset", func(c echo.Context) error {
		base = time.Now()
		return c.String(200, "OK")
	})

	// Seasonality
	e.GET("/seasonality", func(c echo.Context) error {
		d := time.Now().Second()
		if d > 30 {
			d = 60 - d
		}
		return c.String(200, fmt.Sprintf("{\"value\": %d}", d+rand.Intn(5)))
	})

	e.Logger.Fatal(e.Start(":1323"))
}
