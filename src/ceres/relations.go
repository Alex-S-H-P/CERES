package ceres


// A is a subclass of B
type HYPONYMY struct{
    A Entity
    B *EntityType
}

func (h HYPONYMY) reverse() Link {
    return HYPERNYMY{A:h.B, B:h.A}
}

func (h HYPONYMY) GetA() Entity {
    return h.A
}

func (h HYPONYMY) GetB() Entity {
    return h.B
}

func (h HYPONYMY) set(A, B Entity) Link {
    h.A = A
    h.B = B.(*EntityType)
    return h
}

func (h HYPONYMY) typeOfLink() string {return "HYPONYM"}

// A is the superclass of B
type HYPERNYMY struct {
    A *EntityType
    B Entity
}

func (h HYPERNYMY) reverse() Link {
    return HYPONYMY{A:h.B, B:h.A}
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

func (h HYPERNYMY) typeOfLink() string {return "HYPERNYM"}

type Link interface {
    reverse()Link
    GetA()Entity
    GetB()Entity

    // returns a new link that is set.
    set (Entity, Entity) Link

    // Ability to say what type of link you are
    typeOfLink() string
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

func (m MERONYMY) typeOfLink() string {return "MERONYM"}

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

func (h HOLONYMY) typeOfLink() string {return "HOLONYM"}

// A == -B
type ANTONYMY struct {
    A, B Entity
}

func (a ANTONYMY) reverse() Link {
    return ANTONYMY{A:a.B, B:a.A}
}

func (a ANTONYMY) GetA() Entity {
    return a.A
}

func (a ANTONYMY) GetB() Entity {
    return a.B
}

func (a ANTONYMY) set(A, B Entity) Link {
    return ANTONYMY{A:A, B:B}
}

func (a ANTONYMY) typeOfLink() string {return "ANTONYM"}

func SetList(listType string, A, B Entity) Link {
    return FindListType(listType).set(A, B)
}

func FindListType(lt string) Link {
    switch lt {
    case "ANTONYM":
        return ANTONYMY{}
    case "HOLONYM":
        return HOLONYMY{}
    case "MERONYM":
        return MERONYMY{}
    case "HYPERNYM":
        return HYPERNYMY{}
    case "HYPONYM":
        return HYPONYMY{}
    default:
        return nil
    }
}
