package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/labstack/echo"
	"github.com/sdwolfe32/trumail/api"
	"github.com/sdwolfe32/trumail/config"
	"github.com/sdwolfe32/trumail/verifier"
	"github.com/sirupsen/logrus"
)

func main() {
	// Generate a new logrus logger
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	l := logger.WithField("port", config.Port)

	// Define all required dependencies
	l.Info("Defining all service dependencies")
	e := echo.New()
	v := verifier.NewVerifier(retrievePTR(), config.SourceAddr)
	s := api.NewService(logger, config.HTTPClientTimeout, v)

	// Bind endpoints to router
	l.Info("Binding API endpoints to the router")
	e.GET("/v1/:format/:email", s.Lookup)
	e.GET("/v1/health", s.Health)

	// Listen and Serve
	l.WithField("port", config.Port).Info("Listening and Serving")
	l.Fatal(e.Start(":" + config.Port))
}

// retrievePTR attempts to retrieve the PTR record for the IP
// address retrieved via an API call on api.ipify.org
func retrievePTR() string {
	// Request the IP from ipify
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		log.Fatal("Failed to retrieve IP from api.ipify.org")
	}
	defer resp.Body.Close()

	// Decodes the IP response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read IP response body")
	}

	// Retrieve the PTR record for our IP and return without a trailing dot
	names, err := net.LookupAddr(string(data))
	if err != nil {
		return string(data)
	}
	return strings.TrimSuffix(names[0], ".")
}
