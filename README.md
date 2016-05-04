# A Short URL service written in Go

## What is it
This is a simple Short URL service written in Go and backed by Redis that allows for API-based creation of URL redirects that also tracks the number of times the redirect has been activated. 
 
## Config Values (see `main.go` for additional details)
* `TOKEN` - A token that provides some level of authentication on the API (override with environment `TOKEN` var)
* `PORT` - The port the service should run on (override with environment `PORT` var)
* `DOMAIN` - The domain the service should use to build url redirects (override with environment `DOMAIN` var)

## Redis
This Short URL service should be setup with the configuration options for your Redis instance in the `main.go` file.

## Considerations
This service will process `application/json` and `text/plain` POST/PUT bodies for creating and updating Short URLs.

## Redirect
It is pretty simple. Point a browser or issue a `GET` request at an establish Short URL in the service and it will redirect you.
```GET https://example.com/2efdfg```

## API:
*IMPORTANT:* Make sure to set a `X-Auth-Token` header set to the value being used for the `TOKEN` config.

### GET /api > Get all Short URLs
`$ curl -H "X-Auth-Token: some-sort-of-token" "https://example.com/api"`
*Response*
```json
[
    {
        "slug": "fownom",
        "url": "http://google.com",
        "redirect": "https://example.com/fownom",
        "count": 1
    },
    {...}
]
```

### GET /api/{slug} > Get a specific Short URL info
`HTTP GET: https://example.com/api/fownom`
*Response*
```json
{
    "slug": "fownom",
    "url": "http://google.com",
    "redirect": "https://example.com/fownom",
    "count": 1
}
```

### POST/PUT
There are two approachs to posting/putting data.  
* POST/PUT a JSON-based object with the `Content-Type: application/json` header
* POST/PUT a TEXT-based URL string with no header

### POST & PUT - Text-based
/api/{slug} > Create or Modify a specific Short URL
`$ curl -X POST -H "X-Auth-Token: some-sort-of-token" -d 'http://google.com' "https://example.com/api/fownom"`
`$ curl -X PUT -H "X-Auth-Token: some-sort-of-token" -d 'http://google.com' "https://example.com/api/fownom"`
*Response*
```json
{
    "slug": "fownom",
    "url": "http://google.com",
    "redirect": "https://example.com/fownom",
    "count": 1
}
```

### POST - JSON-based
/api/{slug} > Create a specific Short URL
`$ curl -X PUT -H "X-Auth-Token: some-sort-of-token" -d '{"slug":"fownom", "url":"http://google.com"}' "https://example.com/api"`
*Response*
```json
{
    "slug": "fownom",
    "url": "http://google.com",
    "redirect": "https://example.com/fownom",
    "count": 1
}
```

### PUT - JSON-based
/api/{slug} > Modify a specific Short URL
`$ curl -X PUT -H "X-Auth-Token: some-sort-of-token" -d '{"url":"http://google.com"}' "https://example.com/api/fownom"`
*Response*
```json
{
    "slug": "fownom",
    "url": "http://google.com",
    "redirect": "https://example.com/fownom",
    "count": 1
}
```

### DELETE /api/{slug} > Remove a specific Short URL info
`$ curl -X DELETE -H "X-Auth-Token: some-sort-of-token" "https://example.com/api"`
*Response*
```json
{
    "result": "fownom has been deleted"
}
```
