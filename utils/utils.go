package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Duration time.Duration

var unitsMultipliers = map[string]time.Duration{
	"s": time.Second,
	"m": time.Minute,
	"h": time.Hour,
	"d": time.Hour * 24,
	"w": time.Hour * 24 * 7,
}

func (this *Duration) UnmarshalJSON(b []byte) error {
	match, _ := regexp.MatchString(`^['"][0-9]+[smhdw]['"]$`, string(b))
	if !match {
		return errors.New("Cannot unmarshal value for type utils.Duration")
	}
	v, _ := strconv.Unquote(string(b))
	unit := strings.ToLower(string(v[len(v)-1]))
	val, _ := strconv.Atoi(v[:len(v)-1])
	multiplier := unitsMultipliers[unit]

	*this = Duration(val) * Duration(multiplier)
	return nil
}
