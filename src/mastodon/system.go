package mastodon

import (
    "fmt"
    "io/ioutil"
    "net"
)

func ReadLoadAvg() (one, five, fifteen float64) {
    buffer, err := ioutil.ReadFile("/proc/loadavg")
    if err != nil {
        panic(err)
    }
    fmt.Sscanf(string(buffer), "%f %f %f", &one, &five, &fifteen)
    return
}

func IfaceAddr(name string) string {
    iface, _ := net.InterfaceByName(name)
    addrs, _ := iface.Addrs()
    if len(addrs) > 0 {
        return addrs[0].String()
    } else {
        return "n/a"
    }
}
