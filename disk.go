package mastodon

import "syscall"

func diskUsage(path string) (free float64, total float64) {
    // Return bytes free and total bytes.
    buf := new(syscall.Statfs_t)
    syscall.Statfs("/", buf)
    free = float64(buf.Bsize) * float64(buf.Bfree)
    total = float64(buf.Bsize) * float64(buf.Blocks)
    return
}

func Disk(c *Config) *StatusInfo {
    data := make(map[string]string)
    free, total := diskUsage("/")
    freePercent := 100 * (free / total)
    data["bar"] = MakeBar(freePercent, c.BarSize)
    data["free"] = HumanFileSize(free)
    data["total"] = HumanFileSize(free)
    data["used"] = HumanFileSize(total - free)
    si := NewStatus(c.Templates["disk"], data)
    if freePercent < 10 {
        si.Status = STATUS_BAD
    }
    return si
}
