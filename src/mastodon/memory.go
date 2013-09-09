package mastodon

import "syscall"

func MemInfo() (free float64, total float64) {
    // Return free and total bytes of RAM.
    buf := new(syscall.Sysinfo_t)
    syscall.Sysinfo(buf)
    return float64(buf.Freeram), float64(buf.Totalram)
}
