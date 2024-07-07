package subhuman

import (
	"fmt"
	"strings"
)

func Strappend(p, a string) string {
	return strings.Join([]string{p, a}, "")
}

func HumanTimeParse(s int64) map[string]int64 {
	smap := map[string]int64{
		"sec": 1,
		"min": 60,
		"hr":  (60 * 60),
	}
	keys := []string{"hr", "min", "sec"}
	retv := map[string]int64{}
	for _, value := range keys {
		if s >= smap[value] {
			val := (s / smap[value])
			tgt := (val * smap[value])
			rmn := (s - tgt)
			retv[value] = val
			s = rmn
		}
	}
	return retv
}

func HumanTimeColon(s int64) string {
	data := HumanTimeParse(s)
	return fmt.Sprintf("%02d:%02d:%02d", data["hr"], data["min"], data["sec"])
}

func HumanTimeConcise(s int64) string {
	retv := ""
	data := HumanTimeParse(s)
	keys := []string{"hr", "min", "sec"}
	for _, val := range keys {
		retv = Strappend(retv, fmt.Sprintf("%02d%s", data[val], string(val[0])))
	}
	return retv
}
