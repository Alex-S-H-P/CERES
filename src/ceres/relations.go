package ceres


// A is a subclass of B
type HYPNOMY struct{
    A Entity
    B *EntityType
}

func (h HYPNOMY) reverse() Link {
    return HYPERNYMY{A:h.B, B:h.A}
}

func (h HYPNOMY) GetA() Entity {
    return h.A
}

func (h HYPNOMY) GetB() Entity {
    return h.B
}

func (h HYPNOMY) set(A, B Entity) Link {
    h.A = A
    h.B = B.(*EntityType)
    return h
}

// A is the superclass of B
type HYPERNYMY struct {
    A *EntityType
    B Entity
}

func (h HYPERNYMY) reverse() Link {
    return HYPNOMY{A:h.B, B:h.A}
}

func (h HYPERNYMY) GetA() Entity {
    return h.A
}

func (h HYPERNYMY) GetB() Entity {
    return h.B
}

func (h HYPERNYMY) set(A, B Entity) Link {
    h.A = A.(*EntityType)
    h.B = B
    return h
}

type Link interface {
    reverse()Link
    GetA()Entity
    GetB()Entity

    // returns a new link that is set.
    set (Entity, Entity) Link
}
