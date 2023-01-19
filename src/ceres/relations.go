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

// A is a part of B
type MERONYMY struct {
    A Entity
    B Entity
}

func (m MERONYMY) reverse() Link {
    return HOLONYMY{A:m.B, B:m.A}
}

func (m MERONYMY)GetA() Entity {
    return m.A
}

func (m MERONYMY)GetB() Entity {
    return m.B
}

func (m MERONYMY) set(A, B Entity) Link {
    return MERONYMY{A:A, B:B}
}

// A contains B
type HOLONYMY struct {
    A Entity
    B Entity
}

func (h HOLONYMY) reverse() Link {
    return MERONYMY{A:h.B, B:h.A}
}

func (h HOLONYMY)GetA() Entity {
    return h.A
}

func (h HOLONYMY)GetB() Entity {
    return h.B
}

func (h HOLONYMY) set(A, B Entity) Link {
    return HOLONYMY{A:A, B:B}
}
