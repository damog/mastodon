// Status-bar
package mastodon

import (
    "bytes"
    "fmt"
    "os"
    "runtime"
    "syscall"
    "time"
)

type Config struct {
    Data map[string]string
    BarSize int
}

const (
    STATUS_GOOD = iota
    STATUS_BAD
    STATUS_NORMAL
)

type StatusInfo struct {
    FullText string
    Status int64
}

func NewStatus() *StatusInfo {
    si := new(StatusInfo)
    si.Status = STATUS_NORMAL
    return si
}

func (si *StatusInfo) IsGood() bool {
    return si.Status == STATUS_GOOD
}

func (si *StatusInfo) IsBad() bool {
    return si.Status == STATUS_BAD
}

func getBarString(percent float64, bar_size int) string {
    var bar bytes.Buffer
    cutoff := int(percent * .01 * float64(bar_size))
    bar.WriteString("[")
    for i := 0; i < bar_size; i += 1 {
        if i <= cutoff {
            bar.WriteString("#")
        } else {
            bar.WriteString(" ")
        }
    }
    bar.WriteString("]")
    return bar.String()
}

func Battery(c *Config) *StatusInfo {
    si := NewStatus()
    bi := ReadBatteryInfo(0)
    barString := getBarString(bi.PercentRemaining, c.BarSize)
    prefix := "BAT"
    if bi.IsCharging() {
        prefix = "CHR"
    }
    if bi.IsFull() {
        prefix = "FULL"
        si.FullText = fmt.Sprintf("%s %s", prefix, barString)
    } else {
        si.FullText = fmt.Sprintf(
            "%s %s (%s %.1fW)",
            prefix,
            barString,
            HumanDuration(int64(bi.SecondsRemaining)),
            bi.Consumption)
    }
    if bi.PercentRemaining < 15 {
        si.Status = STATUS_BAD
    }
    return si
}

func CPU(c *Config) *StatusInfo {
    si := NewStatus()
    cpuUsage := CpuUsage()
    barString := getBarString(cpuUsage, c.BarSize)
    si.FullText = fmt.Sprintf("C %s", barString)
    if cpuUsage > 80 {
        si.Status = STATUS_BAD
    }
    return si
}

func Disk(c *Config) *StatusInfo {
    si := NewStatus()
    free, total := DiskUsage("/")
    freePercent := 100 * (free / total)
    barString := getBarString(freePercent, c.BarSize)
    si.FullText = fmt.Sprintf("D %s", barString)
    if (free / total) < .1 {
        si.Status = STATUS_BAD
    }
    return si
}

func Memory(c *Config) *StatusInfo {
    si := NewStatus()
    free, total := MemInfo()
    percentUsed := 100 * (total - free) / total
    si.FullText = fmt.Sprintf("R %s", getBarString(percentUsed, c.BarSize))
    if percentUsed > 75 {
        si.Status = STATUS_BAD
    }
    return si
}

func LoadAvg(c *Config) *StatusInfo {
    si := NewStatus()
    cpu := float64(runtime.NumCPU())
    one, five, fifteen := ReadLoadAvg()
    si.FullText = fmt.Sprintf("%.2f %.2f %.2f", one, five, fifteen)
    if one > cpu {
        si.Status = STATUS_BAD
    }
    return si
}

func Clock(c *Config) *StatusInfo {
    si := NewStatus()
    si.FullText = time.Now().Format("2006-01-02 15:04:05")
    return si
}

func IPAddress(c *Config) *StatusInfo {
    si := NewStatus()
    si.FullText = IfaceAddr("wlan0")
    return si
}

func Hostname(c *Config) *StatusInfo {
    si := NewStatus()
    host, _ := os.Hostname()
    si.FullText = host
    return si
}

func Uptime(c *Config) *StatusInfo {
    buf := new(syscall.Sysinfo_t)
    syscall.Sysinfo(buf)
    si := NewStatus()
    si.FullText = fmt.Sprintf("U: %s", HumanDuration(buf.Uptime))
    return si
}

// Cache weather data.
var latestWeatherStatus StatusInfo
var latestWeatherCheck time.Time

func Weather(c *Config) *StatusInfo {
    si := NewStatus()
    if time.Since(latestWeatherCheck) < (time.Duration(1800) * time.Second) {
        return &latestWeatherStatus
    }
    forecast, err := ReadWeather(c.Data["weather_key"], c.Data["weather_zip"])
    latestWeatherCheck = time.Now()
    if err != nil || len(forecast.Forecast.SimpleForecast.ForecastDay) == 0 {
        si.FullText = "Error fetching weather"
        si.Status = STATUS_BAD
    } else {
        today := forecast.Forecast.SimpleForecast.ForecastDay[0]
        si.FullText = fmt.Sprintf(
            "%s H %s, L %s",
            today.Conditions,
            today.High.Fahrenheit,
            today.Low.Fahrenheit)
    }
    latestWeatherStatus = *si
    return si
}
