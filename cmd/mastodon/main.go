package main

import (
    "encoding/json"
    "fmt"
    "github.com/coleifer/mastodon"
    "os"
    //"os/exec"
    "path/filepath"
    "strconv"
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
}

func getDefaultConfig() mastodon.Config {
    var config mastodon.Config
    config.Data = map[string]string{
        "interval": "1",
        "order": "cpu,memory,disk,battery,ip,loadavg,clock",
        "bar_size": "10",
        "color_good": "#00d000",
        "color_normal": "#cccccc",
        "color_bad": "#d00000",
        "color0": "#1e2320",
        "color1": "#705050",
        "color2": "#60b48a",
        "color3": "#dfaf8f",
        "color4": "#506070",
        "color5": "#dc8cc3",
        "color6": "#8cd0d3",
        "color7": "#dcdccc",
        "color8": "#709080",
        "color9": "#dca3a3",
        "color10": "#c3bf9f",
        "color11": "#f0dfaf",
        "color12": "#94bff3",
        "color13": "#ec93d3",
        "color14": "#93e0e3",
        "color15": "#ffffff",
        "cpu": "color2",
        "memory": "color4",
        "battery": "color3",
        "disk": "color5",
        "loadavg": "color6",
    }
    config.BarSize, _ = strconv.Atoi(config.Data["bar_size"])
    return config
}

func ReadConfig(c mastodon.Config) {
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
            if _, ok := c.Data[key]; ok {
                c.Data[key] = value
            }
            return true
        }
        mastodon.ReadLines(configFile, LineHandler)
    }
}

func ReadInterval(c mastodon.Config) time.Duration {
    interval, err := strconv.Atoi(c.Data["interval"])
    if err != nil {
        interval = 1
    }
    return time.Duration(interval) * time.Second
}

func PrintHeader() {
    fmt.Println("{\"version\":1}")
    fmt.Println("[")
}

func main() {
    config := getDefaultConfig()
    ReadConfig(config)
    duration := ReadInterval(config)

    module_names := strings.Split(config.Data["order"], ",")
    jsonArray := make([]map[string]string, len(module_names))
    PrintHeader()
    for {
        for idx, module_name := range(module_names) {
            si := Modules[module_name](&config)
            color := config.Data["color_normal"]
            if si.IsBad() {
                color = config.Data["color_bad"]
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
