package main

import (
    "encoding/json"
    "fmt"
    "github.com/coleifer/mastodon"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

type StatusSource func() *mastodon.StatusInfo

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
}

var config = map[string]string{
    "interval": "1",
    "order": "cpu,memory,disk,ip,battery,loadavg,clock",
    "color_good": "#00d000",
    "color_normal": "#cccccc",
    "color_bad": "#d00000",
}

func ReadConfig() {
    configHome := os.Getenv("XDG_CONFIG_HOME")
    if configHome == "" {
        configHome = filepath.Join(os.Getenv("HOME"), ".config")
    }
    configFile := filepath.Join(configHome, "mastodon.conf")
    if mastodon.FileExists(configFile) {
        LineHandler := func(line string) bool {
            pieces := strings.Split(line, "=")
            key := strings.Trim(pieces[0], " \t\r")
            value := strings.Trim(pieces[1], " \t\r")
            if _, ok := config[key]; ok {
                config[key] = value
            }
            return true
        }
        mastodon.ReadLines(configFile, LineHandler)
    }
}

func PrintHeader() {
    fmt.Println("{\"version\":1}")
    fmt.Println("[")
}

func main() {
    ReadConfig()
    interval, err := strconv.Atoi(config["interval"])
    if err != nil {
        interval = 1
    }
    duration := time.Duration(interval) * time.Second
    module_names := strings.Split(config["order"], ",")
    jsonArray := make([]map[string]string, len(module_names))
    PrintHeader()
    for {
        for idx, module_name := range(module_names) {
            si := Modules[module_name]()
            color := config["color_normal"]
            if si.IsGood() {
                color = config["color_good"]
            } else if si.IsBad() {
                color = config["color_bad"]
            }
            jsonArray[idx] = map[string]string{
                "full_text": si.FullText,
                "color": color,
            }
        }
        jsonData, _ := json.Marshal(jsonArray)
        fmt.Printf(string(jsonData))
        fmt.Printf(",\n")
        time.Sleep(duration)
    }
}
