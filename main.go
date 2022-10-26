package main 

import (
	"fmt"

	// "github.com/jaeyoony/phone_normalizer/normalizer"
	"github.com/jaeyoony/phone_normalizer/sequel"
)

func main() {
	fmt.Println("Phone normalizer main!")
	sequel.SqlMain()
}