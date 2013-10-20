package main

import (
    "encoding/json"
    "fmt"
    "github.com/coleifer/mastodon"
    "strings"
    "time"
)

type StatusSource func(*mastodon.Config) *mastodon.StatusInfo

var Modules = map[string]StatusSource{
    "battery": mastodon.Battery,
    "clock": mastodon.Clock,
    "cpu": mastodon.CPU,
    "disk": mastodon.Disk,
    "hostname": mastodon.Hostname,
    "ip": mastodon.IPAddress,
    "loadavg": mastodon.LoadAvg,
    "memory": mastodon.Memory,
    "uptime": mastodon.Uptime,
    "weather": mastodon.Weather,
}

func PrintHeader() {
    fmt.Println("{\"version\":1}")
    fmt.Println("[")
}

func LoadConfig() *mastodon.Config {
    config := mastodon.NewConfig()
    config.ApplyXresources()
    config.ReadConfig()
    config.ParseTemplates()
    return config
}

func main() {
    config := LoadConfig()
    duration := config.ReadInterval()

    module_names := strings.Split(config.Data["order"], ",")
    for _, module_name := range(module_names) {
        if _, ok := config.Data[module_name]; !ok {
            config.Data[module_name] = "color_normal"
        }
    }

    jsonArray := make([]map[string]string, len(module_names))
    PrintHeader()
    for {
        for idx, module_name := range(module_names) {
            si := Modules[module_name](config)
            color := config.Data[module_name]
            if si.IsBad() {
                color = config.Data["color_bad"]
            }
            if _, ok := config.Data[color]; ok {
                color = config.Data[color]
            }
            jsonArray[idx] = map[string]string{
                "full_text": si.FullText,
                "color": color,
            }
        }
        jsonData, _ := json.Marshal(jsonArray)
        fmt.Print(string(jsonData))
        fmt.Printf(",\n")
        time.Sleep(duration)
    }
}
