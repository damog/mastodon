package mastodon


import (
    "bytes"
    "text/template"
)


const (
    STATUS_GOOD = iota
    STATUS_BAD
    STATUS_NORMAL
)

type StatusInfo struct {
    FullText string
    Status int64
}

func NewStatus(t *template.Template, data map[string]string) *StatusInfo {
    var buf bytes.Buffer
    si := new(StatusInfo)
    si.Status = STATUS_NORMAL
    t.Execute(&buf, data)
    si.FullText = buf.String()
    return si
}

func (si *StatusInfo) IsGood() bool {
    return si.Status == STATUS_GOOD
}

func (si *StatusInfo) IsBad() bool {
    return si.Status == STATUS_BAD
}
