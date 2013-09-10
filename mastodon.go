// Status-bar
package mastodon

import (
    "fmt"
    "os"
    "runtime"
    "syscall"
    "time"
)

type Config map[string]string

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

func Battery(c *Config) *StatusInfo {
    si := NewStatus()
    bi := ReadBatteryInfo(0)
    prefix := "BAT"
    if bi.IsCharging() {
        prefix = "CHR"
    }
    if bi.IsFull() {
        prefix = "FULL"
        si.FullText = fmt.Sprintf(
            "%s %.1f%%",
            prefix,
            bi.PercentRemaining)
    } else {
        si.FullText = fmt.Sprintf(
            "%s %.1f%% %s (%.1fW)",
            prefix,
            bi.PercentRemaining,
            HumanDuration(int64(bi.SecondsRemaining)),
            bi.Consumption)
    }
    if bi.PercentRemaining < 15 {
        si.Status = STATUS_BAD
    } else if bi.PercentRemaining < 75 {
        si.Status = STATUS_NORMAL
    } else {
        si.Status = STATUS_GOOD
    }
    return si
}

func CPU(c *Config) *StatusInfo {
    si := NewStatus()
    cpuUsage := CpuUsage()
    si.FullText = fmt.Sprintf("C %.1f%%", cpuUsage)
    if cpuUsage < 15 {
        si.Status = STATUS_GOOD
    } else if cpuUsage < 75 {
        si.Status = STATUS_NORMAL
    } else {
        si.Status = STATUS_BAD
    }
    return si
}

func Disk(c *Config) *StatusInfo {
    si := NewStatus()
    free, total := DiskUsage("/")
    si.FullText = fmt.Sprintf("D %s/%s", HumanFileSize(free), HumanFileSize(total))
    if (free / total) < .1 {
        si.Status = STATUS_BAD
    } else {
        si.Status = STATUS_GOOD
    }
    return si
}

func Memory(c *Config) *StatusInfo {
    si := NewStatus()
    free, total := MemInfo()
    percentUsed := 100 * (total - free) / total
    si.FullText = fmt.Sprintf("R %.1f%%", percentUsed)
    if percentUsed > 75 {
        si.Status = STATUS_BAD
    } else if percentUsed < 25 {
        si.Status = STATUS_GOOD
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
    } else {
        si.Status = STATUS_GOOD
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
