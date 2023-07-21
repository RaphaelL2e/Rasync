package worker

import "context"

type Cancel struct {
	ctx context.Context
}

func (c *Cancel) IsCancel() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}
