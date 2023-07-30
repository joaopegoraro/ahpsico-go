package handlers

import "fmt"

const (
	pending int64 = iota
	done
	missed
)

// Remove this when consts are used
func TempAssignmentEnumInit() {
	fmt.Print(pending)
	fmt.Print(done)
	fmt.Print(missed)
}
