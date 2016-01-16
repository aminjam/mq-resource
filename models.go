package resource

import "time"

const (
	WaitFor time.Duration = 5
)

type Source struct {
	Queue  string    `json:"queue"`
	Uri    string    `json:"uri"`
	Pub    string    `json:"pub"`
	Sub    string    `json:"sub"`
	Params StringMap `json:"params"`
}

type StringMap map[string]string

func (from *StringMap) IsEqual(to StringMap) bool {
	for k, v := range *from {
		if to[k] != v {
			return false
		}
		delete(to, k)
	}
	if len(to) > 0 {
		return false
	}
	return true
}
