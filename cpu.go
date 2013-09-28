package mastodon

import (
    "strconv"
    "strings"
)

var prevTotal, prevIdle uint64

func cpuUsage() (cpuUsage float64) {
    // Return the percent utilization of the CPU.
    var idle, total uint64
    callback := func(line string) bool {
        fields := strings.Fields(line)
        if fields[0] == "cpu" {
            numFields := len(fields)
            for i := 1; i < numFields; i++ {
                val, _ := strconv.ParseUint(fields[i], 10, 64)
                total += val
                if i == 4 {
                    idle = val
                }
            }
            return false
        }
        return true
    }
    ReadLines("/proc/stat", callback)

    if prevIdle > 0 {
        idleTicks := float64(idle - prevIdle)
        totalTicks := float64(total - prevTotal)
        cpuUsage = 100 * (totalTicks - idleTicks) / totalTicks
    }
    prevIdle = idle
    prevTotal = total
    return
}

func CPU(c *Config) *StatusInfo {
    data := make(map[string]string)
    cpuUsage := cpuUsage()
    data["bar"] = MakeBar(cpuUsage, c.BarSize)
    si := NewStatus(c.Templates["cpu"], data)
    if cpuUsage > 80 {
        si.Status = STATUS_BAD
    }
    return si
}
