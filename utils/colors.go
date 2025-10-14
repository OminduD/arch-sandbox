package utils

import (
	"fmt"
	"log"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

// Info prints a formatted informational message in cyan.
func Info(format string, a ...interface{}) {
	fmt.Printf(ColorCyan+"[INFO] "+ColorReset+format+"\n", a...)
}

// Success prints a formatted success message in green.
func Success(format string, a ...interface{}) {
	fmt.Printf(ColorGreen+"[SUCCESS] "+ColorReset+format+"\n", a...)
}

// Warn prints a formatted warning message in yellow.
func Warn(format string, a ...interface{}) {
	fmt.Printf(ColorYellow+"[WARN] "+ColorReset+format+"\n", a...)
}

// Fatal prints a formatted fatal error message in red and exits.
func Fatal(format string, a ...interface{}) {
	// We use log.Fatalf to ensure the timestamp and exit behavior are consistent.
	log.Fatalf(ColorRed+"[FATAL] "+ColorReset+format, a...)
}
