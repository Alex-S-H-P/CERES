package ceres

import (
    "fmt"
    "CERES/src/utils"
)

var UNKNOWN_GRAMMAR group = group{name:"[⁇Unknown⁇]", instanceSolver:nil}

type EntityType struct {
    word Word
    links []Link

    attributes *AttributeTypeList
    grammar_group group
}

func (et *EntityType)Initialize() {
    et.attributes = new(AttributeTypeList)
    et.links = make([]Link, 0, 8)
    et.grammar_group = UNKNOWN_GRAMMAR
}

type EntityInstance struct {
    typeOf *EntityType // ensure that this is not null !
    otherLinks []Link
    values *AttributeInstanceList
}

func (ei *EntityInstance) Initialize() {
    ei.values = new(AttributeInstanceList)
    ei.otherLinks = make([]Link, 0, 8)
}

func (ei *EntityInstance) directTypeOf() []*EntityType{
    if ei == nil {return nil}

    return []*EntityType{ei.typeOf}
}

func (et *EntityType) directTypeOf() []*EntityType{
    if et == nil {return nil}

    answ := make([]*EntityType, len(et.links)/2)
    for _, link := range et.links {
        if hypo, ok := link.(HYPONYMY); ok {
            answ = append(answ, hypo.GetB().(*EntityType))
        }
    }
    return answ
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
    // FIXME: Detect lack of parent, panic in this case
    return ei.directTypeOf()[0].GetGender()
}

func (ei*EntityInstance)GetNumber() int8 {
    // FIXME: Now create a group entity subtype
    return SINGULAR
}


// interface entity is both EntityType and EntityInstance
type Entity interface {
    directTypeOf() []*EntityType
    Initialize()
    addLink(Link, Entity)(int, error)
    Equal(utils.Equalable) bool
    GetNumber() int8
    GetGender() int8
    removeLink(int)
    hasLink(Link, Entity) bool

    // loading and saving methods
    store(int, map[Entity]int, *[]byte) string
    load([]string, map[string]group, map[int]Entity)
}

/*
Returns true if the entity e is of type t.
There can be intermediate types between the two
*/
func IsTypeOf(e Entity, t *EntityType) bool {
    for _, parent := range (e.directTypeOf()) {
        if parent == nil {
            continue
        } else if parent.Equal(t) {
            return true
        } else {
            if IsTypeOf(parent, t) {
                return true
            }
        }
    }
    return false
}

func lAB_internal(e *EntityType, f Entity) (int, error) {
    if e.Equal(f) {return 0, nil}

    parents := f.directTypeOf()

    var minLAB int = 1<<30
    var minError error = fmt.Errorf("e is not type of f")

    for _, parent := range parents {
        if parent == nil { continue }

        i, err := lAB_internal(e, parent)
        i ++
        if err == nil && ( i < minLAB || minError != nil) {
            minLAB = i
            minError = err
        }
    }
    return minLAB, minError
}

func lAB_type_checked(e, f Entity) (int, error){
    if E, ok := e.(*EntityType); ok {
        return lAB_internal(E, f)
    } else {
        return 0, fmt.Errorf("Could not convert to entity type")
    }
}


/*
Returns how many ancestors you need to go up to find e going from f.
The result can be negative if f is an ancestor of e instead.

May return a non-nil error:
  - if the data is not set correctly,
  - if there is no link between e and f

In which case the integer returned is the number of ancestors checked.
*/
func LevelsOfAbstractionBetween(e, f Entity) (int, error) {
    if e.Equal(f) {
        return 0, nil
    }
    _, ok1 := e.(*EntityInstance)
    _, ok2 := f.(*EntityInstance)
    if ok1 && ok2 {
        return 1, fmt.Errorf("Both entities checked to compute levels of abstractions are instances. They are not equal.")
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

func (et*EntityType) hasLink(linkType Link, destination Entity) bool {
    for _, l := range et.links {
        if l.typeOfLink() == linkType.typeOfLink() && l.GetB().Equal(destination){
            return true
        }
    }
    return false
}

func (ei*EntityInstance) hasLink(linkType Link, destination Entity) bool {
    if linkType.typeOfLink() == "HYPONYM" {
        return ei.typeOf == destination
    }

    for _, l := range ei.otherLinks {
        if l.typeOfLink() == linkType.typeOfLink() && l.GetB().Equal(destination){
            return true
        }
    }
    return false
}

func (et*EntityType) addChild(childEntity Entity) error {
    return AddLink(HYPERNYMY{}, et, childEntity)
}

func (et *EntityType) addLink(emptyLink Link, destination Entity) (int, error) {
    link := emptyLink.set(et, destination)
    et.links = append(et.links, link)

    return len(et.links)-1, nil
}

func (ei*EntityInstance) addLink(emptyLink Link, destination Entity) (int, error) {
    _, is_hypnomy := emptyLink.(HYPONYMY)
    destination_as_class, ok := destination.(*EntityType)
    if is_hypnomy && ok && ei.typeOf == nil {
        ei.typeOf = destination_as_class
        return -1, nil
    } else if ei.typeOf != nil{
        return -1, fmt.Errorf("Cannot set typeOf of instance with \"%s\", as it already is of class \"%s\"",
                            ei.typeOf.word, destination_as_class.word)
    } else if !ok {
        return -1, fmt.Errorf("Cannot set typeOf of instance with non class %v", destination)
    }


    if _, ok := emptyLink.(HYPERNYMY); ok {
        return -1, fmt.Errorf("Nothing can be subclass of an entity")
    }


    link := emptyLink.set(ei, destination)
    ei.otherLinks = append(ei.otherLinks, link)
    return len(ei.otherLinks)-1, nil
}

// adds a link between source and destination..
// You can pass an empty link to specify the type of the link
func AddLink(emptyLink Link, source, destination Entity) error {
    var i int
    var err error
    var canRemove bool

    if !source.hasLink(emptyLink, destination) {
        i, err = source.addLink(emptyLink, destination)
        if err != nil {
            return err
        }
        canRemove = true
    }

    if !destination.hasLink(emptyLink.reverse(), source) {
        _, err = destination.addLink(emptyLink.reverse(), source)
        if err != nil {
            if canRemove {
                source.removeLink(i)
            }
            return err
        }
    }

    return nil
}


func (et*EntityType) removeLink(index int) {
    if index < len(et.links) - 1 {
        et.links = append(et.links[:index], et.links[index+1:]...)
    } else {
        et.links = et.links[:index]
    }
}

func (ei*EntityInstance) removeLink(index int) {
    if index == -1 {
        ei.typeOf = nil
    } else if index < len(ei.otherLinks) - 1 {
        ei.otherLinks = append(ei.otherLinks[:index], ei.otherLinks[index+1:]...)
    } else {
        ei.otherLinks = ei.otherLinks[:index]
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

// returns the lowest entityType such that e1 and e2 descend from it
// if none are found, returns nil
// if e1 and e2 are entityType and equal, they return themselves
func ClosestAncestor(e1, e2 Entity)*EntityType {
    var head1, head2 *EntityType
    if et1, ok := e1.(*EntityType); ok {
        head1 = et1
    } else {
        head1 = e1.directTypeOf()[0]
    }
    if et2, ok := e2.(*EntityType); ok {
        head2 = et2
    } else {
        head2 = e2.directTypeOf()[0]
    }

    return _closestAncestor(head1, head2, nil)
}

func _closestAncestor(head1, head2 *EntityType,
                                m map[*EntityType]byte) *EntityType {

    if m ==  nil {
        m = make(map[*EntityType]byte)
        m[head1] = 1
        m[head2] = 2
    }

    if head1.Equal(head2) {return head1}


    var heads1, heads2 []*EntityType = head1.directTypeOf(), head2.directTypeOf()

    for len(heads1) + len(heads2) > 0 {
        var nheads1, nheads2 []*EntityType = make([]*EntityType, 0, len(heads1)),
                                             make([]*EntityType, 0, len(heads2))
        for _, h1 := range heads1 {
            found, continue_execution := _testIfClosestAncestor(h1, m, 1)
            if found {return h1}
            if !continue_execution {continue} // loop found

            nheads1 = append(nheads1, h1.directTypeOf()...)
        }

        for _, h2 := range heads2 {
            found, countinue_execution := _testIfClosestAncestor(h2, m, 2)
            if found {return h2}
            if !countinue_execution {continue}

            nheads2 = append(nheads2, h2.directTypeOf()...)
        }
    }


    return nil
}

func _testIfClosestAncestor(head *EntityType,
                            m map[*EntityType]byte,
                            _side byte) (bool, bool) {

    if side, ok := m[head]; ok {
        if side == _side {
            return false, false // we looped
        } else {
            return true, true // going up from e2 made us go through head1. That's the closest.
        }
    } else {
        m[head] = _side
    }

    return false, true
}
