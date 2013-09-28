// Status-bar
package mastodon

import (
    "fmt"
    "os"
    "runtime"
    "syscall"
    "time"
)


func LoadAvg(c *Config) *StatusInfo {
    one, five, fifteen := ReadLoadAvg()
    data := make(map[string]string)
    data["one"] = fmt.Sprintf("%.2f", one)
    data["five"] = fmt.Sprintf("%.2f", five)
    data["fifteen"] = fmt.Sprintf("%.2f", fifteen)
    si := NewStatus(c.Templates["loadavg"], data)
    cpu := float64(runtime.NumCPU())
    if one > cpu {
        si.Status = STATUS_BAD
    }
    return si
}

func Clock(c *Config) *StatusInfo {
    data := make(map[string]string)
    data["time"] = time.Now().Format(c.Data["date_format"])
    si := NewStatus(c.Templates["clock"], data)
    return si
}

func IPAddress(c *Config) *StatusInfo {
    data := make(map[string]string)
    data["ip"] = IfaceAddr(c.Data["network_interface"])
    si := NewStatus(c.Templates["ip"], data)
    return si
}

func Hostname(c *Config) *StatusInfo {
    data := make(map[string]string)
    data["hostname"], _ = os.Hostname()
    si := NewStatus(c.Templates["hostname"], data)
    return si
}

func Uptime(c *Config) *StatusInfo {
    data := make(map[string]string)
    buf := new(syscall.Sysinfo_t)
    syscall.Sysinfo(buf)
    data["uptime"] = HumanDuration(buf.Uptime)
    si := NewStatus(c.Templates["uptime"], data)
    return si
}
