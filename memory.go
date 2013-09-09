package mastodon

import (
    "fmt"
    "strings"
)

func MemInfo() (free float64, total float64) {
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
