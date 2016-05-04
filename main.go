/**
 * main.go
 * Short URL service based on Redis
 * Author: Brian Powell
 * GitHub: github.com/brianpowell/go-shorturl
 */

package main

import (
	"net/http"
	"os"

	"github.com/brianpowell/go-shorturl/routes"
	"github.com/swhite24/go-debug"
	"gopkg.in/redis.v2"
)

var (
	debug = debugger.NewDebugger("go-shorturl:main")
	// API header Token X-Auth-Token
	TOKEN = "some-sort-of-token"

	// Port to listen on (override with environment PORT var)
	PORT = "80"

	// Domain to use (override with environment DOMAIN var)
	DOMAIN = "https://example.com"

	// Configure for your Redis settings
	r_opts = &redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
)

func main() {

	// Redis Connection
	client := redis.NewTCPClient(r_opts)

	// Token Env override
	if os.Getenv("TOKEN") != "" {
		TOKEN = os.Getenv("TOKEN")
	}

	// Port Env override
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	// DOMAIN Env override
	if os.Getenv("DOMAIN") != "" {
		DOMAIN = os.Getenv("DOMAIN")
	}

	debug.Log("Domain being used", DOMAIN)

	// Build the Routes
	mux := routes.BuildRoutes(client, TOKEN, DOMAIN)

	// Start serving stuff
	// Modify, as necessary to
	debug.Log("Listening on port", PORT)
	http.ListenAndServe(":"+PORT, mux)

}
