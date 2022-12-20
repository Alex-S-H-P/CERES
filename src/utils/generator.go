package utils

type Generator[T any] struct {
    next func() (T, bool) // the boolean value returns true when the last element is given
    stop func()

    returned chan T
    stopchan chan bool

    nb int
    stopped bool
    el T // currently considered element
}

func (g*Generator[T]) Start() {
    g.returned = make(chan T)
    g.stopchan = make(chan bool)
    g.nb = 0
    g.stopped = false

    go func(){
        g.nb ++
        for {
            g.el, g.stopped = g.next()
            if g.stopped {
                return
            }
            select {
            case g.returned <- g.el:
            case <- g.stopchan:
                g.stopped = true
                return
            }
        }
    }()
}

// returns a slice of all of the remaining elements.
func (g*Generator[T]) Slice() []T {
    slice := make([]T, 32)
    slice = append(slice, g.el)
    for {
        el, finished := g.next()
        slice = append(slice, el)
        if finished {break}
    }
    return slice
}


// Returns the next element of the generator, if the generator is finished, returns true. Does return the last element
func (g*Generator[T]) Next() (T, bool){
    default_t := *(new(T))
    if g.stopped {
        return default_t, true
    } else {
        el := <- g.returned
        return el, false
    }

}
