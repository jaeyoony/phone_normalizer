package normalizer 


import (
	"regexp"
)

func NormalizeNumber(num string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(num, "")
}