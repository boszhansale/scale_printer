package utils

import (
	"fmt"
	"time"
)

func DateToCode() string {

	currentDate := time.Now()
	formattedDate := currentDate.Format("0201")
	fmt.Println(formattedDate)

	return formattedDate
}
