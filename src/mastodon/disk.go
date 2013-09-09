package mastodon

import "syscall"

func DiskUsage(path string) (free float64, total float64) {
    // Return bytes free and total bytes.
    buf := new(syscall.Statfs_t);
	syscall.Statfs("/", buf);
	free = float64(buf.Bsize) * float64(buf.Bfree)
	total = float64(buf.Bsize) * float64(buf.Blocks)
	return
}
