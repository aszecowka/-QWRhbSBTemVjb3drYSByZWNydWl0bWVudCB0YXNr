# GogoApps Recruitment Task


## Development
To build, test and execute static checks, execute the following command:
```bash
make
```
To run the program, execute the following command:
```bash
export APP_WEATHER_API_KEY={your key}
make run

```

Program is configurable via environment variables

| Environment variable          | Default value                     | Description                                       |                                                                             
| ------------------------------|-----------------------------------|---------------------------------------------------| 
| **APP_PORT**                  | 8080                              | Port on which application listens on
| **APP_CACHE_ADDRESS**         | localhost:6379                    | Redis host
| **APP_CACHE_PASSWORD**        |                                   | Redis password, if exist
| **APP_CACHE_DB**              |                                   | Redis database number
| **APP_CACHE_TTL**             | 5 minutes                         | How long weather data are kept in cache
| **APP_WEATHER_API_URL**       | https://api.openweathermap.org    | Open Weather Map API URL
| **APP_WEATHER_API_ KEY**      |                                   | Open Weather Map API Key (required)
| **APP_WEATHER_API_TIMEOUT**   | 1 second                          | Timeout for requests to Open Weather Map API

To run a containerized application, execute the following commands:
```
export API_KEY={your key}
make run-with-deps
```

To run integration tests, execute the following commands:
```
export API_KEY={your key}
make run-int-test
```

## Proposed solution
- API
There is one handler, that accepts GET requests and requires `city` query parameter:
```
curl -v -X GET "http://localhost:8080/?city=Katowice&city=Gliwice&city=DoesNotExist"
```

The `city` query parameter can be defined many times. The response is a map, where the key is a requested city, and value is a response from OpenWeatherMap service.
In case, that city does not exist, value is set to null.
```
{
  "Gliwice": {
    "base": "stations",
    "main": {
      "temp": 278.17,
      "feels_like": 274.3,
      "temp_min": 275.93,
      "temp_max": 279.82,
      "pressure": 1001,
      "humidity": 80
    },
    "id": 3099230,
    "name": "Gliwice"
  },
  "Katowice": {
    "base": "stations",
    "main": {
      "temp": 278.04,
      "feels_like": 274.15,
      "temp_min": 275.93,
      "temp_max": 279.82,
      "pressure": 1001,
      "humidity": 80
    },
    "timezone": 3600,
    "id": 3096472,
    "name": "Katowice"
  },
  "DoesNotExist": null
}
```

- Redis is used for caching responses. Environment variable `APP_CACHE_TTL` defines how long responses are cached.
- There can be a different requests for the same city:
```
api.openweathermap.org/data/2.5/weather?q=London
api.openweathermap.org/data/2.5/weather?q=London,uk
```
Thanks to the collection `cityNameToID` stored in the Redis, requests for `London` and `London,uk` will use the same cached value.

## Next steps

Below you can find a list of possible improvements:
- provide Open API Specification
- add metrics to detect how often cache miss happens. This can be useful for adjusting **APP_CACHE_TTL** parameter.
- protect from malicious users who bypass cache by requesting data for non-existing cities
