package utils

type Equalable interface {
    Equal(other Equalable) bool
}

type Buffer[T Equalable] struct {
    contents []T
    maxsize int
    start int
}


func NewBuffer[T Equalable](maxSize int) *Buffer[T] {
    b := new(Buffer[T])
    b.maxsize = maxSize
    b.start = 0
    b.contents  = make([]T, 0, maxSize)
    return b
}

func (b *Buffer[T])Put(t T) {
    if len(b.contents) < b.maxsize - 1 {
        b.contents = append(b.contents, t)
    } else {
        b.contents[b.start] = t
        b.start = (b.start + 1) % b.maxsize
    }
}

func (b *Buffer[T])Get(idx int) T {
    return b.contents[(b.start + idx)%b.maxsize]
}

func (b *Buffer[T])Contains(t T) bool {
    for i := 0; i < len(b.contents); i++ {
        if t.Equal(b.contents[i]) {
            return true
        }
    }
    return false
}
