//go:build darwin

package extract

import (
	"context"
	"fmt"

	"github.com/sjzar/chatlog/internal/wechat"
)

// ExtractPlatformKeys on macOS.
func ExtractPlatformKeys() ([]string, error) {
	wechat.Load()
	accounts := wechat.GetAccounts()
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no wechat account found")
	}

	var keys []string
	for _, acc := range accounts {
		key, err := acc.GetKey(context.Background())
		if err != nil {
			// Log error but continue to try other accounts
			fmt.Printf("Could not get key for account %s: %v\n", acc.Name, err)
			continue
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("failed to extract key from any account")
	}

	return keys, nil
}
