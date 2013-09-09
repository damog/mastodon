package mastodon

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"
)

// Time.
const (
    SECOND = 1
    MINUTE = SECOND * 60
    HOUR = MINUTE * 60
    DAY = HOUR * 24
    WEEK = DAY * 7
    YEAR = WEEK * 52
)

const (
    KB = 1024
    MB = KB * 1024
    GB = MB * 1024
    TB = GB * 1024
)
var diskUsage = []uint64{TB, GB, MB, KB}
var humanUsage = map[uint64]string{
    KB: "KB",
    MB: "MB",
    GB: "GB",
    TB: "TB",
}

func HumanTime(n, resolution int64) string {
    var idx int64
    parts := make([]string, 6)
    addPart := func(part int64, label string) {
        if n > part {
            val := n / part
            n = n % part
            parts[idx] = fmt.Sprintf("%d %s", val, label)
            idx += 1
        }
    }
    addPart(YEAR, "years")
    addPart(WEEK, "weeks")
    addPart(DAY, "days")
    addPart(HOUR, "hours")
    addPart(MINUTE, "minutes")
    addPart(SECOND, "seconds")
    if idx > resolution {
        idx = resolution
    }
    return strings.Join(parts[:idx], ", ")
}

func HumanDuration(n int64) string {
    hours := n / 3600
    minutes := (n % 3600) / 60
    return fmt.Sprintf("%d:%02d", hours, minutes)
}

func HumanFileSize(n float64) (s string) {
    for _, size := range(diskUsage) {
        fsize := float64(size)
        if fsize < n {
            return fmt.Sprintf("%.1f %s", n / fsize, humanUsage[size])
        }
    }
    return fmt.Sprintf("%fb", n)
}

func ReadLines(fileName string, callback func(string) bool) {
    fin, err := os.Open(fileName)
    if err != nil {
        fmt.Fprintf(os.Stderr, "The file %s does not exist!\n", fileName)
        return
    }
    defer fin.Close()

    reader := bufio.NewReader(fin)
    for line, _, err := reader.ReadLine(); err != io.EOF; line, _, err = reader.ReadLine() {
        if !callback(string(line)) {
            break
        }
    }
}

func FileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
