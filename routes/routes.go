/**
 * routes/routes.go
 * Entry point for URLs to be redirected
 * Author: Brian Powell
 * GitHub: github.com/brianpowell/go-shorturl
 */

package routes

import (
	"github.com/brianpowell/go-shorturl/controllers"
	"github.com/gorilla/mux"
	"github.com/swhite24/go-debug"
	"gopkg.in/redis.v2"
)

var (
	debug = debugger.NewDebugger("go-shorturl:routes:routes")
)

// BuildRoutes is the entry point redirect URLs
func BuildRoutes(c *redis.Client, t string, d string) *mux.Router {

	mux := mux.NewRouter()

	api := controllers.NewAPI(c, t, d)
	mux.HandleFunc("/api", api.List).Methods("GET")
	mux.HandleFunc("/api", api.HandleAPI).Methods("POST")
	mux.HandleFunc("/api/", api.HandleAPI).Methods("POST")
	mux.HandleFunc("/api/{slug}", api.HandleAPI).Methods("GET", "POST", "PUT", "DELETE")

	redirect := controllers.NewRedirect(c)
	mux.HandleFunc("/", redirect.HandleRedirect).Methods("GET")
	mux.HandleFunc("/{slug}", redirect.HandleRedirect).Methods("GET")

	return mux
}
