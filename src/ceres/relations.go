package ceres

import (
	"fmt"
	"strings"
)

// A is a subclass of B
type HYPONYMY struct {
	A Entity
	B *EntityType
}

func (h HYPONYMY) reverse() Link {
	return HYPERNYMY{A: h.B, B: h.A}
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

func (h HYPONYMY) typeOfLink() string { return "HYPONYM" }

// A is the superclass of B
type HYPERNYMY struct {
	A *EntityType
	B Entity
}

func (h HYPERNYMY) reverse() Link {
	return HYPONYMY{A: h.B, B: h.A}
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

func (h HYPERNYMY) typeOfLink() string { return "HYPERNYM" }

type Link interface {
	reverse() Link
	GetA() Entity
	GetB() Entity

	// returns a new link that is set.
	set(Entity, Entity) Link

	// Ability to say what type of link you are
	typeOfLink() string
}

// A is a part of B
type MERONYMY struct {
	A Entity
	B Entity
}

func (m MERONYMY) reverse() Link {
	return HOLONYMY{A: m.B, B: m.A}
}

func (m MERONYMY) GetA() Entity {
	return m.A
}

func (m MERONYMY) GetB() Entity {
	return m.B
}

func (m MERONYMY) set(A, B Entity) Link {
	return MERONYMY{A: A, B: B}
}

func (m MERONYMY) typeOfLink() string { return "MERONYM" }

// A contains B
type HOLONYMY struct {
	A Entity
	B Entity
}

func (h HOLONYMY) reverse() Link {
	return MERONYMY{A: h.B, B: h.A}
}

func (h HOLONYMY) GetA() Entity {
	return h.A
}

func (h HOLONYMY) GetB() Entity {
	return h.B
}

func (h HOLONYMY) set(A, B Entity) Link {
	return HOLONYMY{A: A, B: B}
}

func (h HOLONYMY) typeOfLink() string { return "HOLONYM" }

// A == -B
type ANTONYMY struct {
	A, B Entity
}

func (a ANTONYMY) reverse() Link {
	return ANTONYMY{A: a.B, B: a.A}
}

func (a ANTONYMY) GetA() Entity {
	return a.A
}

func (a ANTONYMY) GetB() Entity {
	return a.B
}

func (a ANTONYMY) set(A, B Entity) Link {
	return ANTONYMY{A: A, B: B}
}

func (a ANTONYMY) typeOfLink() string { return "ANTONYM" }

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
	case "SUPERLATIVE":
		return SUPERLATIVITY{}
	case "COMPARATIVE":
		return COMPARATIVITY{}
	default:
		if isRadical(lt) {
			return handleRadical(lt)
		} else {
			panic(fmt.Sprintf("Cannot validate %v as a link type (choose between ANTONYM, HOLONYM, MERONYM, HYPERNYM, HYPONYM", lt))
		}
	}
}

type word_derivation interface {
	Link
	SUPERLATIVITY | COMPARATIVITY
}

func handleRadical(lt string) Link {
	stype := strings.TrimPrefix(strings.TrimSuffix(lt, "]"), "RADICAL[")
	switch stype {
	case "SUPERLATIVE":
		return RADICAL[SUPERLATIVITY]{}
	case "COMPARATIVE":
		return RADICAL[COMPARATIVITY]{}
	default:
		panic("Cannot handle unknown word relation " + lt)
	}
}

func isRadical(lt string) bool {
	return strings.HasPrefix(lt, "RADICAL[") && strings.HasSuffix(lt, "]")
}

type SUPERLATIVITY struct {
	A, B Entity
}

func (s SUPERLATIVITY) reverse() Link {
	return RADICAL[SUPERLATIVITY]{A: s.B, B: s.A}
}

func (s SUPERLATIVITY) GetA() Entity {
	return s.A
}

func (s SUPERLATIVITY) GetB() Entity {
	return s.B
}

func (s SUPERLATIVITY) set(A, B Entity) Link {
	s.A = A
	s.B = B
	return s
}

func (s SUPERLATIVITY) typeOfLink() string {
	return "SUPERLATIVE"
}

type COMPARATIVITY struct {
	A, B Entity
}

func (c COMPARATIVITY) reverse() Link {
	return RADICAL[COMPARATIVITY]{A: c.B, B: c.A}
}

func (c COMPARATIVITY) GetA() Entity {
	return c.A
}

func (c COMPARATIVITY) GetB() Entity {
	return c.B
}

func (c COMPARATIVITY) set(A, B Entity) Link {
	c.A = A
	c.B = B
	return c
}

func (c COMPARATIVITY) typeOfLink() string {
	return "COMPARATIVE"
}

type RADICAL[L word_derivation] struct {
	A, B Entity
}

func (r RADICAL[L]) reverse() Link {
	return L{}
}

func (r RADICAL[L]) GetA() Entity {
	return r.A
}
func (r RADICAL[L]) GetB() Entity {
	return r.B
}

func (r RADICAL[L]) set(A, B Entity) Link {
	r.A = A
	r.B = B
	return r
}

func (r RADICAL[L]) typeOfLink() string {
	return "RADICAL[" + L{}.typeOfLink() + "]"
}
