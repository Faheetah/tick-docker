package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

const indexHtml = `
<html>
	<body>
		<h2>Links:</h2>
		<h4><a href="/chronograf/" onclick="javascript:event.target.port=8888">Chronograf</a></h4>
		<h4><a href="/grafana/"    onclick="javascript:event.target.port=3000">Grafana</a></h4>
		<h4><a href="/golerta/"    onclick="javascript:event.target.port=5608">Golerta</a></h4>
		<h2>Modify State:</h3>
		<h4>Current status: %d</h4>
		<h4>Current trend value: %02f</h4>
		<h4><a href="/status/200">Send 200</a></h4>
		<h4><a href="/status/404">Send 404</a></h4>
		<h4><a href="/status/500">Send 500</a></h4>
		<h4><a href="/trending/reset">Reset Trending</a></h4>
    </body>
</html>
`

func weightedRand() (n int) {
	n = 100 + (1000/(1+rand.Intn(100)))
	return
}

var (
	base = time.Now()
	status = 200
)

func main() {
	rand.Seed(time.Now().UnixNano())
	e := echo.New()

	// Always return a static metric
	e.GET("/", func(c echo.Context) error {
		d := time.Now().Sub(base)
		t := d.Seconds() / 10
		return c.HTML(http.StatusOK, fmt.Sprintf(indexHtml, status, t))
	})

	// Mimic status codes by pushing to the server
	e.GET("/status", func(c echo.Context) error {
		r := weightedRand()
		time.Sleep(time.Duration(r) * time.Millisecond)
		return c.String(status, fmt.Sprintf("{\"value\": %d}", r))
	})

	e.GET("/status/:status", func(c echo.Context) error {
		status, _ = strconv.Atoi(c.Param("status"))
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	// Random value
	e.GET("/random", func(c echo.Context) error {
		return c.String(200, fmt.Sprintf("{\"value\": %d}", rand.Intn(100)))
	})

	// Trending
	e.GET("/trending", func(c echo.Context) error {
		d := time.Now().Sub(base)
		t := (d.Seconds() + float64(rand.Intn(30))) / 10
		return c.String(200, fmt.Sprintf("{\"value\": %02f}", t))
	})

	e.GET("/trending/reset", func(c echo.Context) error {
		base = time.Now()
		return c.Redirect(http.StatusTemporaryRedirect, "/")
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
