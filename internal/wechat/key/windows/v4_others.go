//go:build !windows

package windows

import (
	"context"

	"github.com/dylan-zzm/xiaomingai/internal/wechat/model"
)

func (e *V4Extractor) Extract(ctx context.Context, proc *model.Process) (string, error) {
	return "", nil
}
