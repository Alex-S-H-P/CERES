package ceres

import (
    "fmt"
)

type EntityType struct {
    parent *EntityType
    attributes *AttributeTypeList
    children *[]Entity
}

func (et *EntityType)Initialize() {
    et.attributes = new(AttributeTypeList)
    et.children = new([]Entity)
}

type EntityInstance struct {
    typeOf *EntityType // ensure that this is not null !
    values *AttributeInstanceList
}

func (ei *EntityInstance) Initialize() {
    ei.values = new(AttributeInstanceList)
}

func (ei *EntityInstance) directTypeOf() *EntityType{
    return ei.typeOf
}

func (et *EntityType) directTypeOf() *EntityType{
    return et.parent
}

// interface entity is both EntityType and EntityInstance
type Entity interface {
    directTypeOf() *EntityType
    Initialize()
}

/* returns true if the entity e is of type t.
There can be intermediate types between the two */
func IsTypeOf(e Entity, t *EntityType) bool {
    if e.directTypeOf() == nil {
        return false
    } else if *e.directTypeOf() == *t {
        return true
    } else {
        return IsTypeOf(e.directTypeOf(), t)
    }
}

func lAB_internal(e *EntityType, f Entity) (int, error) {
    if F, ok := f.(*EntityType); ok && *e == *F{
        return 0, nil
    } else if f.directTypeOf() != nil {
        i, e := lAB_internal(e, f.directTypeOf())
        return i+1, e
    } else {
        return 0, fmt.Errorf("e is not type of f")
    }
}

func lAB_type_checked(e, f Entity) (int, error){
    if E, ok := e.(*EntityType); ok {
        return lAB_internal(E, f)
    } else {
        return -1, fmt.Errorf("Could not convert to entity type")
    }
}

func LevelsOfAbstractionBetween(e, f Entity) (int, error) {
    E, ok1 := e.(*EntityInstance)
    F, ok2 := f.(*EntityInstance)
    if ok1 && ok2 {
        if *E == *F {
            return 0, nil
        } else {
            return 1, fmt.Errorf("Both entities checked to compute levels of abstractions are instances. They are not equal.")
        }
    } else {
        a, err := lAB_type_checked(e, f)
        if err != nil {
            b, err2 := lAB_type_checked(f, e)
            if err2 != nil {
                return a+b, fmt.Errorf("Could not find a level of abstraction. Got %s one way and %s the other",
                    err.Error(), err2.Error())
            }
            return -b, err2
        }
        return a, err
    }
}

func (et *EntityType)addChild(e Entity) {
    switch e.(type) {
    case *EntityType:
        et2 := e.(*EntityType)
        et2.parent = et
        *et.children = append(*(et.children), et2)
        et.attributes.parentType(et2.attributes)
    case *EntityInstance:
        ei := e.(*EntityInstance)
        ei.typeOf = et
        *et.children = append(*(et.children), ei)
        et.attributes.parentInstance(ei.values)
    }
}
