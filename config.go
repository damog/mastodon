package mastodon


import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
    "text/template"
    "time"
)


type StatusSource func(*Config) *StatusInfo


type Config struct {
    Data map[string]string
    BarSize int
    Battery int
    Templates map[string]*template.Template
}


func NewConfig() Config {
    var config Config
    config.Data = map[string]string{
        "bar_size": "10",
        "battery": "0",
        "color_bad": "#d00000",
        "color_good": "#00d000",
        "color_normal": "#cccccc",
        "date_format": "2006-01-02 15:04:05",
        "format_battery": "{{if .battery}}{{.prefix}} {{.bar}} ({{.remaining}} {{.wattage}}W){{else}}No battery{{end}}",
        "format_clock": "{{.time}}",
        "format_cpu": "C {{.bar}}",
        "format_disk": "D {{.bar}}",
        "format_hostname": "{{.hostname}}",
        "format_ip": "{{.ip}}",
        "format_loadavg": "{{.fifteen}} {{.five}} {{.one}}",
        "format_memory": "R {{.bar}}",
        "format_uptime": "{{.uptime}}",
        "format_weather": "{{if .error}}{{.error}}{{ else }}{{.today}} {{.high}}/{{.low}} ({{.next}}){{ end }}",
        "interval": "1",
        "network_interface": "wlan0",
        "order": "weather,cpu,memory,disk,battery,ip,loadavg,clock",
    }
    config.BarSize, _ = strconv.Atoi(config.Data["bar_size"])
    config.Battery, _ = strconv.Atoi(config.Data["battery"])
    return config
}

func (c Config) ApplyXresources() {
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

func (c Config) ReadConfig() {
    configHome := os.Getenv("XDG_CONFIG_HOME")
    if configHome == "" {
        configHome = filepath.Join(os.Getenv("HOME"), ".config")
    }
    configFile := filepath.Join(configHome, "mastodon.conf")
    if FileExists(configFile) {
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
        ReadLines(configFile, LineHandler)
    }
}

func (c Config) ReadInterval() time.Duration {
    interval, err := strconv.Atoi(c.Data["interval"])
    if err != nil {
        interval = 1
    }
    return time.Duration(interval) * time.Second
}

func ParseTemplates(c Config) map[string]*template.Template {
    templates := make(map[string]*template.Template)

    for key, value := range(c.Data) {
        if strings.HasPrefix(key, "format_") {
            name := strings.TrimPrefix(key, "format_")
            t := template.New(name)
            t, err := t.Parse(value)
            if err != nil {
                fmt.Println("Bad template: %s", key)
                panic(err)
            }
            templates[name] = t
        }
    }
    return templates
}
