package mastodon

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "time"
)

type Forecast struct {
    Forecast struct {
        SimpleForecast struct {
            ForecastDay []struct {
                Conditions string
                High struct {
                    Fahrenheit string
                }
                Low struct {
                    Fahrenheit string
                }
            }
        }
    }
}

func readWeather(apiKey string, zipCode string) (*Forecast, error) {
    forecast := new(Forecast)
    url := fmt.Sprintf("http://api.wunderground.com/api/%s/forecast10day/q/%s.json", apiKey, zipCode)
    resp, err := http.Get(url)
    if err == nil {
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err == nil {
            json.Unmarshal(body, &forecast)
        }
    }
    return forecast, err
}

// Cache weather data.
const maxDelay = 1800
var latestWeatherStatus StatusInfo
var latestWeatherCheck time.Time
var cacheDelay = 1

func Weather(c *Config) *StatusInfo {
    if time.Since(latestWeatherCheck) < (time.Duration(cacheDelay) * time.Second) {
        return &latestWeatherStatus
    }

    data := make(map[string]string)
    data["weather_zip"] = c.Data["weather_zip"]

    forecast, err := readWeather(c.Data["weather_key"], c.Data["weather_zip"])
    latestWeatherCheck = time.Now()
    if err != nil || len(forecast.Forecast.SimpleForecast.ForecastDay) == 0 {
        data["error"] = "Error fetching weather"
        cacheDelay *= 2
        if cacheDelay > maxDelay {
            cacheDelay = maxDelay
        }
    } else {
        today := forecast.Forecast.SimpleForecast.ForecastDay[0]
        next := forecast.Forecast.SimpleForecast.ForecastDay[1]
        data["today"] = today.Conditions
        data["high"] = string(today.High.Fahrenheit)
        data["low"] = string(today.Low.Fahrenheit)
        data["next"] = string(next.Conditions)
        cacheDelay = maxDelay
    }
    si := NewStatus(c.Templates["weather"], data)
    if err != nil {
        si.Status = STATUS_BAD
    }
    latestWeatherStatus = *si
    return si
}
