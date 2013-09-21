package mastodon

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
    "net/http"
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

func ReadWeather(apiKey string, zipCode string) (*Forecast, error) {
    forecast := new(Forecast)
    url := fmt.Sprintf("http://api.wunderground.com/api/%s/forecast10day/q/%s.json", apiKey, zipCode)
    resp, err := http.Get(url)
    defer resp.Body.Close()
    if err == nil {
        body, err := ioutil.ReadAll(resp.Body)
        if err == nil {
            json.Unmarshal(body, &forecast)
        }
    }
    return forecast, err
}
