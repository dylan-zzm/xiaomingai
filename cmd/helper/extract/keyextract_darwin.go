//go:build darwin

package extract

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/dylan-zzm/xiaomingai/internal/wechat"
)

// GetKey is the main entry point for extracting the WeChat database key on macOS.
// It requires root privileges to run. It returns the raw 32-byte key.
func GetKey() ([]byte, error) {
	if os.Geteuid() != 0 {
		return nil, errors.New("key extraction requires root privileges, please run with sudo")
	}

	// The upstream `wechat` package handles process finding and version detection.
	wechat.Load()
	accounts := wechat.GetAccounts()
	if len(accounts) == 0 {
		return nil, errors.New("no wechat accounts found running")
	}

	var lastErr error
	for _, acc := range accounts {
		// GetKey from the upstream package handles the complex memory reading and key derivation.
		keyHex, err := acc.GetKey(context.Background())
		if err != nil {
			lastErr = fmt.Errorf("failed for account %s: %w", acc.Name, err)
			continue // Try next account
		}

		keyBytes, err := hex.DecodeString(keyHex)
		if err != nil {
			lastErr = fmt.Errorf("failed to decode hex key for account %s: %w", acc.Name, err)
			continue
		}

		if len(keyBytes) == 32 {
			return keyBytes, nil // Success
		}
		lastErr = fmt.Errorf("extracted key for account %s has incorrect length: %d", acc.Name, len(keyBytes))
	}

	// If all accounts failed, return a comprehensive error.
	if lastErr != nil {
		return nil, fmt.Errorf("all attempts failed. last error: %w. if this persists, you may need to disable System Integrity Protection (SIP)", lastErr)
	}

	return nil, errors.New("could not extract a valid key from any account")
}

// ExtractPlatformKeys is the compatibility function called by the main helper command.
// It returns a list of hex-encoded keys.
func ExtractPlatformKeys() ([]string, error) {
	key, err := GetKey()
	if err != nil {
		return nil, err
	}
	return []string{hex.EncodeToString(key)}, nil
}

