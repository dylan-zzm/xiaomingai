//go:build !darwin && !windows

package extract

import "fmt"

// ExtractPlatformKeys is a placeholder for unsupported platforms.
func ExtractPlatformKeys() ([]string, error) {
	return nil, fmt.Errorf("unsupported platform")
}
