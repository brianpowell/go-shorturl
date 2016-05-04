/**
 * routes/redirect.go
 * Entry point for URLs to be redirected
 * Author: Brian Powell
 * GitHub: github.com/brianpowell/go-shorturl
 */

package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/brianpowell/go-shorturl/models"
	"github.com/gorilla/mux"
	"github.com/swhite24/go-debug"
	"gopkg.in/redis.v2"
)

type (
	// API exposes the API for go-shorturl
	Redirect struct {
		cache *redis.Client
		debug debugger.Debugger
	}
)

// NewRedirect delivers an instance of API
func NewRedirect(cache *redis.Client) *Redirect {
	return &Redirect{
		cache: cache,
		debug: debugger.NewDebugger("go-shorturl:controllers:redirect"),
	}
}

func (a *Redirect) HandleRedirect(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	slug := vars["slug"]

	// Check that we are doing things right
	a.debug.Log("Request URL value", slug)
	if r.Method != "GET" || slug == "" {
		http.NotFound(w, r)
		return
	}

	// Get the information
	val, err := a.cache.Get(slug).Result()
	if err != nil {
		w.Write([]byte("We could not find the Short URL you were looking for."))
		return
	}

	// Parse the data
	sURL := models.ShortURL{}
	e := json.Unmarshal([]byte(val), &sURL)
	if e != nil {
		w.Write([]byte("Something weird happened and we could not redirect you."))
		return
	}

	// Do the redirect
	a.debug.Log("Redirecting to ", sURL.URL)
	http.Redirect(w, r, sURL.URL, http.StatusFound)

	// Bump the stats
	a.stats(sURL)

}

func (a *Redirect) stats(sURL models.ShortURL) {

	// Bump the count
	sURL.Count = sURL.Count + 1

	// Parse the data
	data, _ := json.Marshal(sURL)

	// Put it back in the cache
	a.cache.Set(sURL.Slug, string(data))

}
