package subhuman

import (
	"fmt"
	"math"
)

func HumanSize(s int64) string {
	smap := map[string]int64{
		"Byte": int64(1),
		"Kilo": int64(1024),
		"Mega": int64(math.Pow(float64(1024), float64(2))),
		"Giga": int64(math.Pow(float64(1024), float64(3))),
		"Tera": int64(math.Pow(float64(1024), float64(4))),
		"Peta": int64(math.Pow(float64(1024), float64(5))),
		"Exxa": int64(math.Pow(float64(1024), float64(6))),
	}
	keys := []string{"Exxa", "Peta", "Tera", "Giga", "Mega", "Kilo", "Byte"}
	for index, value := range keys {
		if s >= smap[value] {
			bigB := "B"
			if index == len(keys)-1 {
				bigB = ""
				return fmt.Sprintf("%d%s%s",
					s,
					string(value[0]),
					bigB)
			}
			return fmt.Sprintf("%.02f%s%s",
				float64(s)/float64(smap[value]),
				string(value[0]),
				bigB)
		}
	}
	return ""
}
