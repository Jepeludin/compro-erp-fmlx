//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"ganttpro-backend/utils"
)

func main() {
	// Test passwords
	passwords := map[string]string{
		"admin123": "",
		"13032001": "",
		"20010313": "",
		"ppic2024": "",
		"super123": "",
		"op2024":   "",
		"wh2024":   "",
		"test123":  "",
	}

	fmt.Println("=== Generating Password Hashes ===\n")

	for pass, _ := range passwords {
		hash, err := utils.HashPassword(pass)
		if err != nil {
			fmt.Printf("Error hashing %s: %v\n", pass, err)
			continue
		}
		passwords[pass] = hash
		fmt.Printf("Password: %-12s => Hash: %s\n", pass, hash)
	}

	// Test verification
	fmt.Println("\n=== Testing Password Verification ===\n")

	for pass, hash := range passwords {
		isValid := utils.CheckPasswordHash(pass, hash)
		fmt.Printf("Password: %-12s | Valid: %v\n", pass, isValid)
	}
}
