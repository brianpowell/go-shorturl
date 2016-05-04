/**
 * models/models.go
 * Data structure(s) used in go-shorturl
 * Author: Brian Powell
 * GitHub: github.com/brianpowell/go-shorturl
 */

package models

type (

	// ShortURL
	ShortURL struct {
		Slug     string `json:"slug" bson:"slug"`
		URL      string `json:"url" bson:"url"`
		Redirect string `json:"redirect" bson:"redirect"`
		Count    int    `json:"count" bson:"count"`
	}
)
