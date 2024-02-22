package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r3labs/sse/v2"
)

/*

	This example uses the Echo framework, but the sse package can be used with any
	framework as it uses a default HTTP handler to serve the events.

*/

func main() {
	// create the sseServer
	sseServer := sse.New()
	// create a new stream, in this case we're sending a stringified time
	sseServer.CreateStream("time")

	go publishEvents(sseServer) // run the publisher as a background task

	e := echo.New()
	// add the required middleware to the echo server
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// create a new route for the server sent events
	e.GET("/events", echo.WrapHandler(sseServer)) // use the sseServers default handler

	// create a new route for the index.html file
	e.GET("/", func(c echo.Context) error {
		return c.File("index.html") // send a html file to the client
	})

	// create a new route for the sse.html file
	e.GET("/sse.html", func(c echo.Context) error {
		c.Response().Header().Set("Hx-Push-URL", "/sse") // change url to /sse
		return c.File("sse.html")                        // send a html file to the client
	})

	e.Start(":8080")
}

// publishEvents publishes a new event to the "time" stream every 5 seconds
func publishEvents(s *sse.Server) {
	ticker := time.NewTicker(5 * time.Second)
	for t := time.Now(); true; t = <-ticker.C {
		msg := t.Format("2006-01-02 15:04:05")
		s.Publish("time", &sse.Event{
			Data: []byte(msg),
		})
	}
}
