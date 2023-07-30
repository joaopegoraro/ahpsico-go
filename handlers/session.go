package handlers

import "fmt"

const (
	notConfirmed int64 = iota
	confirmed
	canceled
	concluded
)

// Remove this when consts are used
func TempSessionEnumInit() {
	fmt.Print(notConfirmed)
	fmt.Print(confirmed)
	fmt.Print(canceled)
	fmt.Print(concluded)
}
