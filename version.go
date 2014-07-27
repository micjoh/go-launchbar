package launchbar

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a version string
type Version string

func parseVersion(s string, width int) int64 {
	strList := strings.Split(s, ".")
	format := fmt.Sprintf("%%s%%0%ds", width)
	v := ""
	for _, value := range strList {
		v = fmt.Sprintf(format, v, value)
	}
	var result int64
	var err error
	if result, err = strconv.ParseInt(v, 10, 64); err != nil {
		return 0
	}
	return result
}

// Cmp compares two versions
func (v Version) Cmp(w Version) int {
	return int(parseVersion(string(v), 4) - parseVersion(string(w), 4))
}
