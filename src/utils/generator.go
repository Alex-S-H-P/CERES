package utils

type Generator[T any] interface {
    // Returns the next element of the generator, if the generator is finished, returns true. Does return the last element
    Next() (T, bool)
    // stops the generator, and releases ressources it took
    Stop()
}

type BaseGenerator[T any] struct {
    next func() (T, bool) // the boolean value returns true when the last element is given
    stop func()

    returned chan T
    stopchan chan bool

    nb int
    stopped bool
    el T // currently considered element
}

func (g*BaseGenerator[T]) Start(next func()(T, bool), stop func()) {
    g.next = next
    g.stop = stop
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



// Returns the next element of the generator, if the generator is previously finished, returns true. Does return the last element
func (g*BaseGenerator[T]) Next() (T, bool){
    default_t := *(new(T))
    if g.stopped {
        return default_t, true
    } else {
        el := <- g.returned
        return el, false
    }

}

func (g*BaseGenerator[T])Stop() {
    g.stop()
    g.stopchan <- true
}

// returns a slice of all of the remaining elements.
func Slice[T any](g Generator[T]) []T {
    slice := make([]T, 32)
    for {
        el, finished := g.Next()
        if finished {break}
        slice = append(slice, el)
    }
    return slice
}


type MapGenerator[K comparable, V any] struct {
    stopchan chan bool
    returned chan K
    stopped  bool
}

func (mg*MapGenerator[K, V]) Start(m map[K]V) {
    mg.returned = make(chan K)
    mg.stopchan = make(chan bool)

    go func () {
        for k, _ := range m {
            select {
            case mg.returned<- k:
            case <- mg.stopchan:
                return
            }
        }
        mg.Stop()
    }()
}

// Returns the next element of the generator, if the generator is finished, returns true. Does return the last element
func (mg*MapGenerator[K, V])Next() (K, bool) {
    default_k := *(new(K))
    if mg.stopped {
        return default_k, true
    } else {
        e := <- mg.returned
        return e, false
    }
}

func (mg*MapGenerator[K, V])Stop()  {
    mg.stopchan <- true
}

func (mg*MapGenerator[K,V])Values(m map[K]V) Generator[V] {
    with := func(k K)V{return m[k]}
    return Transform[K,V](mg, with)
}

// Given a transformation from K -> L, transforms the Generator
func Transform[K, L any](from Generator[K], with func(K)L ) Generator[L] {
    var g = new(BaseGenerator[L])

    next := func() (L,bool) {
        n, b := from.Next()
        return with(n),b
    }
    stop := from.Stop

    g.Start(next, stop)

    return g
}

// Given a way to parse children element of T, and a generator of T, returns a generator of all children of all T
func Combine[T, V any](meta Generator[T],
    mesa_parser func(T)Generator[V]) Generator[V]{

    var g = new(BaseGenerator[V])
    var mesa Generator[V]

    next := func()(V, bool){
        if mesa != nil {
            v, bb := mesa.Next()
            if bb {
                mesa_parser = nil
            } else {
                return v, false
            }
        }

        t, b := meta.Next()
        if b {
            return *new(V), true
        }
        mesa = mesa_parser(t)
        return mesa.Next()
    }

    stop := func() {
        meta.Stop()
    }

    g.Start(next, stop)

    return g
}

func SliceGenerator[K any](slice []K)Generator[K] {
    var g = new(BaseGenerator[K])
    var i int
    next := func() (K, bool){
        if i < len(slice) {
            return slice[i], false
        }
        return *(new(K)), true
    }
    stop := func() {}

    g.Start(next, stop)

    return g
}
