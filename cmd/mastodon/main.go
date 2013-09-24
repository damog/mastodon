package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/coleifer/mastodon"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
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
    "weather": mastodon.Weather,
}

func getDefaultConfig() mastodon.Config {
    var config mastodon.Config
    config.Data = map[string]string{
        "interval": "1",
        "order": "weather,cpu,memory,disk,battery,ip,loadavg,clock",
        "bar_size": "10",
        "color_good": "#00d000",
        "color_normal": "#cccccc",
        "color_bad": "#d00000",
    }
    config.BarSize, _ = strconv.Atoi(config.Data["bar_size"])
    return config
}

func ApplyXresources(c *mastodon.Config) {
    out, err := exec.Command("xrdb", "-q").Output()
    if err != nil {
        return
    }
    scanner := bufio.NewScanner(bytes.NewReader(out))
    re := regexp.MustCompile(`.*?(color[\d]+):\s*?(#[A-Za-z0-9]+)`)
    for scanner.Scan() {
        line := scanner.Text()
        for _, match := range(re.FindAllStringSubmatch(line, -1)) {
            if match != nil {
                c.Data[match[1]] = match[2]
            }
        }
    }
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
            if _, ok := c.Data[value]; ok {
                value = c.Data[value]
            }
            c.Data[key] = value
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
    ApplyXresources(&config)
    ReadConfig(config)
    duration := ReadInterval(config)

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
            si := Modules[module_name](&config)
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
