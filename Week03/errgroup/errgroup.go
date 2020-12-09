package errgroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

type Group struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

func (g *Group) Go(task func() error) {
	g.wg.Add(1)

	go func() {
		defer func() {
			g.wg.Done()

			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				g.err = fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
			}
		}()

		if err := task(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func (g *Group) Cancel() {
	g.cancel()
}
