package main

import (
	"fmt"
	"os"
)

// CheckIfRunUnderRoot returns error when application has no root access.
func CheckIfRunUnderRoot() error {
	if os.Getuid() != 0 {
		return fmt.Errorf("%s needs to be run as root", os.Args[0])
	}

	return nil
}
