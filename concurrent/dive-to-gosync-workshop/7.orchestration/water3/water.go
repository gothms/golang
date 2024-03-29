package water

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type H2O struct {
	semaH *semaphore.Weighted
	semaO *semaphore.Weighted
}

func New() *H2O {
	semaO := semaphore.NewWeighted(2)
	semaO.Acquire(context.Background(), 2)

	return &H2O{
		semaH: semaphore.NewWeighted(2),
		semaO: semaO,
	}
}

func (h2o *H2O) hydrogen(releaseHydrogen func()) {
	h2o.semaH.Acquire(context.Background(), 1)

	// releaseHydrogen() outputs "H". Do not change or remove this line.
	releaseHydrogen()

	h2o.semaO.Release(1)
}

func (h2o *H2O) oxygen(releaseOxygen func()) {
	h2o.semaO.Acquire(context.Background(), 2)

	// releaseOxygen() outputs "O". Do not change or remove this line.
	releaseOxygen()

	h2o.semaH.Release(2)
}
