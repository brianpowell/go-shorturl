/**
 * routes/api.go
 * API for controlling shorten'ed URLs
 * Author: Brian Powell
 * GitHub: github.com/brianpowell/go-shorturl
 */

package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/brianpowell/go-shorturl/models"
	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"github.com/swhite24/go-debug"
	"gopkg.in/redis.v2"
)

type (
	// API exposes the API for go-shorturl
	API struct {
		cache  *redis.Client
		debug  debugger.Debugger
		token  string
		domain string
	}
)

// NewAPI delivers an instance of API
func NewAPI(cache *redis.Client, token string, domain string) *API {
	return &API{
		cache:  cache,
		debug:  debugger.NewDebugger("go-shorturl:controllers:api"),
		token:  token,
		domain: domain,
	}
}

func (a *API) List(w http.ResponseWriter, r *http.Request) {

	keys, _ := a.cache.Keys("*").Result()

	set := []models.ShortURL{}
	for _, v := range keys {
		val, _ := a.cache.Get(v).Result()

		temp := models.ShortURL{}
		json.Unmarshal([]byte(val), &temp)

		set = append(set, temp)
	}

	data, _ := json.Marshal(set)

	// Send out the data
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(string(data)))
}

// Handle API requests
func (a *API) HandleAPI(w http.ResponseWriter, r *http.Request) {

	// Check the Token Value
	tok := r.Header.Get("X-Auth-Token")

	if tok != a.token {
		w.Write([]byte("Invalid Permissions"))
		return
	}

	vars := mux.Vars(r)
	slug := vars["slug"]

	var err string
	var out string

	// What are we doing?
	switch r.Method {
	case "GET":
		err, out = a.get(slug)
	case "POST", "PUT":
		err, out = a.postPut(slug, r)
	case "DELETE":
		err, out = a.delete(slug)
	}

	// Errors?
	if err != "" {
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	// Send out the data
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(out))

}

// GET a redirect data object
func (a *API) get(slug string) (string, string) {

	val, err := a.cache.Get(slug).Result()
	if err != nil {
		return "Could not find resources '" + slug + "'", ""
	}

	return "", val
}

// POST/PUT a redirect data object
func (a *API) postPut(slug string, r *http.Request) (string, string) {

	rType := r.Header.Get("Content-Type")

	// Shoult not happen with the route configurations, but worth the check
	if r.Method == "PUT" && slug == "" {
		return "Missing 'slug' URL parameter", ""
	}

	// What Content-Type are we handling?
	sURL := &models.ShortURL{}
	e := ""
	if rType == "application/json" {
		e, sURL = a.handleJSON(slug, r)
	} else {
		e, sURL = a.handleText(slug, r)
	}

	if e != "" {
		return e, ""
	}

	// Pull out storage key
	key := sURL.Slug

	// JSON back to a string
	data, _ := json.Marshal(sURL)

	// Set the data
	a.cache.Set(key, string(data))

	return "", string(data)
}

// DELETE the redirect from the cache
func (a *API) delete(slug string) (string, string) {

	_, err := a.cache.Get(slug).Result()
	if err != nil {
		return "Could not find resources '" + slug + "'", ""
	}

	// Delete the data
	a.cache.Del(slug)

	return "", "{\"result\": \"" + slug + " has been deleted\"}"
}

// Handle JSON-based POST/PUT requests
func (a *API) handleJSON(slug string, r *http.Request) (string, *models.ShortURL) {

	// JSON decoder
	decoder := json.NewDecoder(r.Body)

	val := &models.ShortURL{}
	err := decoder.Decode(&val)

	// Set the Slug information
	if val.Slug == "" && slug != "" {
		val.Slug = slug
	} else if val.Slug == "" {
		val.Slug = a.genHash(6)
	}

	val.Redirect = a.domain + "/" + val.Slug

	if err != nil {
		return "Invalid JSON", nil
	}

	return "", val

}

// Handle text-based POST/PUT requests
func (a *API) handleText(slug string, r *http.Request) (string, *models.ShortURL) {

	// Generate the Slug
	if slug == "" {
		slug = a.genHash(6)
	}

	url, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return "Invalid data.", nil
	}

	sUrl := string(url)

	if sUrl == "" {
		return "Invalid data.", nil
	}

	out := &models.ShortURL{
		Slug:     slug,
		URL:      sUrl,
		Redirect: a.domain + "/" + slug,
		Count:    0,
	}

	return "", out
}

// GenID produces a random alphanumeric hash of the provided length
func (a *API) genHash(l int) string {
	return strings.ToLower(uniuri.NewLen(l))
}
