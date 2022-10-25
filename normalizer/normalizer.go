package normalizer 


import (
	"fmt"
	"strconv"
)


func NormPrint() {
	fmt.Println("Hello from normalizer")
}

func NormalizeNumber(num string) string {
	new_num := ""
	for _, i := range num {
		if _, err := strconv.Atoi(string(i)); err == nil {
			new_num += string(i)
		}
	}
	return new_num
}