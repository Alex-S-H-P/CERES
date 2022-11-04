package ceres

import (
    "fmt"

    "CERES/src/utils"
)

type EntityType struct {
    parent *EntityType
    attributes *AttributeTypeList
    children *[]Entity
    surroundingList surroundingList
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

func (et *EntityType) GetNumber() int8 {
    // TODO: solve this
    return 0
}

func (et *EntityType) GetGender() int8 {
    // TODO: solve this
    return UNKNOWN
}

func (ei*EntityInstance)GetGender() int8 {
    // FIXME: is this working always ? No. Check if there is a locally defined gender
    return ei.directTypeOf().GetGender()
}

func (ei*EntityInstance)GetNumber() int8 {
    // FIXME: Now create a group entity subtype
    return SINGULAR
}


// interface entity is both EntityType and EntityInstance
type Entity interface {
    directTypeOf() *EntityType
    Initialize()
    Equal(utils.Equalable) bool
    GetNumber() int8
    GetGender() int8
}

/*
Returns true if the entity e is of type t.
There can be intermediate types between the two
*/
func IsTypeOf(e Entity, t *EntityType) bool {
    if e.directTypeOf() == nil {
        return false
    } else if e.directTypeOf().Equal(t) {
        return true
    } else {
        return IsTypeOf(e.directTypeOf(), t)
    }
}

func lAB_internal(e *EntityType, f Entity) (int, error) {
    if F, ok := f.(*EntityType); ok && e.Equal(F){
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
    if et.children == nil {
        et.children = new([]Entity)
    }
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

func (e *EntityType) Equal(other utils.Equalable) bool {
    if o, ok := other.(*EntityType); ok {
        return e == o
    } else {return false}
}

func (e *EntityInstance) Equal(other utils.Equalable) bool {
    if o, ok := other.(*EntityInstance); ok {
        return e == o
    } else {return false}
}
