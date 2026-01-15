package network

type RecyclePool[T any] struct {
	Pool []T `json:"pool"`
}

func NewRecyclePool[T any](maxSize int) *RecyclePool[T] {
	if maxSize <= 0 {
		maxSize = 10
	}
	return &RecyclePool[T]{Pool: make([]T, 0, maxSize)}
}

func (pool *RecyclePool[T]) GetResource() (T, bool) {
	var zero T
	if pool.isEmpty() {
		return zero, false
	}
	res := pool.Pool[len(pool.Pool)-1]
	pool.Pool = pool.Pool[:len(pool.Pool)-1]
	return res, true
}

func (pool *RecyclePool[T]) PutResource(res T) error {
	pool.Pool = append(pool.Pool, res)
	return nil
}

func (pool *RecyclePool[T]) isEmpty() bool {
	return len(pool.Pool) == 0
}
