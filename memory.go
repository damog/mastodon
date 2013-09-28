package mastodon

import (
    "fmt"
    "strings"
)

func memInfo() (free float64, total float64) {
	mem := map[string]float64{
		"MemTotal": 0,
		"MemFree": 0,
		"Buffers": 0,
		"Cached": 0,
	}
    callback := func(line string) bool {
        fields := strings.Split(line, ":")
        if _, ok := mem[fields[0]]; ok {
            var val float64
            fmt.Sscanf(fields[1], "%f", &val)
            mem[fields[0]] = val * 1024
        }
        return true
    }
    ReadLines("/proc/meminfo", callback)
    return mem["MemFree"] + mem["Buffers"] + mem["Cached"], mem["MemTotal"]
}

func Memory(c *Config) *StatusInfo {
    data := make(map[string]string)
    free, total := memInfo()
    data["free"] = HumanFileSize(free)
    data["used"] = HumanFileSize(total - free)
    data["total"] = HumanFileSize(total)
    percentUsed := 100 * (total - free) / total
    data["bar"] = MakeBar(percentUsed, c.BarSize)
    si := NewStatus(c.Templates["memory"], data)
    if percentUsed > 75 {
        si.Status = STATUS_BAD
    }
    return si
}
