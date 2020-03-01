# GogoApps Recruitment Task


## Development

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
| **APP_WEATHER_API_ KEY**      |                                   | Open Weather Map API Key
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
There is one handler, that accepts GET reqeust and requires `city` query parameter:
```curl -v -X GET "http://localhost:8080/?city=Katowice&city=Gliwice&city=DoesNotExist"```
City can be defined many times. Response is a map, where key is a requested city, and value is a response from OpenWeatherMap service.
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
- there can be a different requests for the same city:
```
api.openweathermap.org/data/2.5/weather?q=London
api.openweathermap.org/data/2.5/weather?q=London,uk
```
Thanks to the collection `cityNameToID`, requests for `London` and `London,uk` will use the same cached value.
