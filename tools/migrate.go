package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	args := os.Args[1:]
	baseArgs := []string{"-dir", "database/migrations/", "sqlite3", "database/db.sqlite3"}
	cmd := exec.Command("goose", append(baseArgs, args...)...)
	out, _ := cmd.CombinedOutput()

	fmt.Printf("%s", out)
}
