package ceres

type entangledRecognitionNode struct {
    possibilities map[*surrounding][]*recognitionNode
    Content  *RecognizedEntity
    parent   *recognitionNode
}

func (ern*entangledRecognitionNode)copy() *entangledRecognitionNode {
    copy := new(entangledRecognitionNode)
    copy.parent = ern.parent
    copy.Content = ern.Content
    copy.possibilities = make(map[*surrounding][]*recognitionNode)
    for s, slice := range ern.possibilities {
        copy.possibilities[s] = make([]*recognitionNode, len(slice))
        for i, ptr := range copy.possibilities[s] {
            copy.possibilities[s][i] = ptr.copy()
        }
    }
    return copy
}

/*
    Removes the rn from the possibilities.
    If the rn is not in it, does nothing

    If this causes one surrounding to have no corresponding possibility,
    then that surrounding is no longer considered

    If this causes the ern to be empty (no surrounding could be considered),
        it removes itself.
*/
func (ern*entangledRecognitionNode)remove(rn*recognitionNode) {
    if possibility, ok := ern.possibilities[rn.Surround]; ok {
        for i, rncurr := range possibility {
            if rncurr == rn {
                if i == len(ern.possibilities) - 1 {
                    ern.possibilities[rn.Surround] = possibility[:i]
                } else {
                    ern.possibilities[rn.Surround] = append(possibility[:i], possibility[i+1:]...)
                }

                if len(ern.possibilities[rn.Surround]) == 0 {
                    // the surrounding is not possible. Remove it.
                    delete(ern.possibilities, rn.Surround)
                    if len(ern.possibilities) == 0 {
                        // this node is not possible. Remove it.
                        ern.parent.remove(ern)
                    }
                }

                return
            }
        }
    }
}


func (ern*entangledRecognitionNode)nodes()[]*recognitionNode {
    var buffer = make([]*recognitionNode, 0, 2*(len(ern.possibilities)+ 1 ))

    for _, array := range ern.possibilities {
        for _, element := range array {
            buffer = append(buffer, element)
        }
    }
    return buffer
}

type recognitionNode struct {
    tangler  *entangledRecognitionNode
    Surround *surrounding
    ChildMap map[int]*entangledRecognitionNode
}

func (rn *recognitionNode) copy() *recognitionNode{
    if rn == nil {
        return nil
    }
    copy := new(recognitionNode)
    copy.tangler = rn.tangler
    copy.Surround = rn.Surround
    for k, v := range rn.ChildMap {
        copy.ChildMap[k] = v.copy()
        v.parent = copy // we fix the links.
    }
    return copy
}


/*
Finds if the content represented by rn is already in m.
If not, adds it and returns true.
If so, returns whether that expression of the content is rn.
*/
func (rn *recognitionNode) is_incompatible(m map[*RecognizedEntity]*recognitionNode) bool {
    if RN, ok := m[rn.tangler.Content]; ok {
        return RN == rn
    } else {
        m[rn.tangler.Content] = rn
        return true
    }
}

/*
    Adds the r.e. as allowed by the surround. r.e. is considered to be a child
    If multiple are allowed, makes the change bubble up to be stored into a
    If none are allowed, returns false, otherwise, returns true

    - rn may be removed, if there are multiple possibilities, but not if there is none.
*/
func (rn*recognitionNode) add (re *RecognizedEntity, left bool) ([]*entangledRecognitionNode, bool) {
    var possible []int

    if left {
        possible = rn.Surround.MatchLeft(re.entity)
    } else {
        possible = rn.Surround.MatchRight(re.entity)
    }

    var answer []*entangledRecognitionNode = make([]*entangledRecognitionNode, 0, len(possible))

    if rn.ChildMap == nil {rn.ChildMap = make(map[int]*entangledRecognitionNode)}

    switch len(possible) {
    case 0:
        // this token does not match this RecognitionNode
        return nil, false
    case 1:
        if _, ok := rn.ChildMap[possible[0]]; !ok {
            rn.ChildMap[possible[0]] = &entangledRecognitionNode{Content:re,
                                            parent:rn,
                                            possibilities:make(map[*surrounding][]*recognitionNode)}
        } else {
            // we can't fit two recognized entities for the same role
            rn.tangler.remove(rn)
            return nil, false
        }
    default:
        for _, pos := range possible {
            copy := rn.copy()
            if _, ok := copy.ChildMap[pos]; ok {
                // this copy can fit
                copy.ChildMap[pos] = &entangledRecognitionNode{Content:re,
                    parent:rn,
                    possibilities:make(map[*surrounding][]*recognitionNode)}
                    rn.tangler.possibilities[rn.Surround] = append(rn.tangler.possibilities[rn.Surround], copy)
                answer = append(answer, copy.ChildMap[pos])
            } else {
                // the pre-existing child makes it impossible to consider the copy viable
                rn.tangler.remove(copy)
            }
        }
        rn.tangler.remove(rn)
    }
    return answer, len(answer) != 0
}

/*
Removes the child ern from this node.
If the node ends up being empty(no child), it tries to remove itself from its own entanglement
*/
func (rn*recognitionNode) remove (ern*entangledRecognitionNode) {
    for k, v := range rn.ChildMap {
        if ern == v {
            delete(rn.ChildMap, k)
            break
        }
    }
    if len(rn.ChildMap) == 0 {
        rn.tangler.remove(rn)
    }
}


func (rn*recognitionNode)children_on_the_right() []*entangledRecognitionNode {
    if rn == nil {return nil}

    var children []*entangledRecognitionNode = make([]*entangledRecognitionNode,
        0, rn.Surround.maxPos)
    for i:=1; i<rn.Surround.maxPos; i++ {
        if child, ok := rn.ChildMap[i]; ok {
            children = append(children, child)
        }
    }
    return children
}


/*
Takes all of the descendants of rn, and counts how many could be fathered by re.
*/
func (rn*recognitionNode) try_unspooling_children(ern*entangledRecognitionNode) int {
    var children_to_be_tried = [][]*entangledRecognitionNode{rn.children_on_the_right()}
    var finished bool = len(children_to_be_tried[0]) == 0
    if finished {return 0}

    var re *RecognizedEntity = ern.Content

    var last_children = []*entangledRecognitionNode{children_to_be_tried[0][len(children_to_be_tried) - 1]}
    var count int = 0
    var presumed_new_children = 2 * len(last_children[0].possibilities)



    for !finished {
        // counters
        var last_children_processed int = 0

        // process the children groups
        for _, children_group := range children_to_be_tried {

            for _, s := range re.surroundings().surr {
                var surrounding_specific_rn = &recognitionNode{tangler:ern,
                                                    Surround:s,
                                                    ChildMap:make(map[int]*entangledRecognitionNode)}

                ern.possibilities[s] = append(ern.possibilities[s],
                                                surrounding_specific_rn)
                CHILD_FINDER:
                for i := len(children_group)-1; i >= 0; i-- { // we loop from right to left
                    content := children_group[i].Content
                    /* universes where this is child of re are not
                    compatible with universes where this is compatible with their original parent
                    */
                    if _, ok := surrounding_specific_rn.add(content, true); !ok {
                        // content can't be a child of rn through that surrounding.
                        break CHILD_FINDER
                    }
                }
            }
        }

        // update children groups
        new_children_to_be_tried := make([][]*entangledRecognitionNode, 0, presumed_new_children)
        new_last_children := make([]*entangledRecognitionNode ,0, presumed_new_children)
        presumed_new_children = 0

        for i := range last_children {
            for _, v := range last_children[i].possibilities {
                for _, rnc := range v {
                    new_ones := rnc.children_on_the_right()
                    if len(new_ones) == 0 {
                        continue
                    }
                    new_children_to_be_tried = append(new_children_to_be_tried,
                        append(children_to_be_tried[i], new_ones...))

                    new_last_child := new_ones[len(new_ones)-1]
                    new_last_children = append(new_last_children, new_last_child)

                    presumed_new_children += len(new_last_child.possibilities)

                    last_children_processed ++
                }
            }
        }
        children_to_be_tried = new_children_to_be_tried
        last_children = new_last_children

        // are we finished ?
        finished = last_children_processed == 0
    }
    return count
}


/*
    Considers adding the recognized entity as a child of this node.
    This does not generate the possibility that the RecognizedEntity is not a fitting child

    Returns whether or not the re could be added
*/
func (rn*recognitionNode) Add (re *RecognizedEntity) ([]*entangledRecognitionNode, bool) {
    var answer []*entangledRecognitionNode
    var ok bool

    answer, ok = rn.add(re, false)

    if !ok {
        // there is no matching possibility. Remove the node.
        rn.tangler.remove(rn)
        return nil, false
    }
    return answer, true
}

func (ern*entangledRecognitionNode) Add(re *RecognizedEntity) {
    for _, possibility := range ern.nodes() {
        if erns, ok := possibility.Add(re); ok {
            for _, child := range erns{
                possibility.copy().try_unspooling_children(child) //
            }
        }
    }
}

func (s*surrounding)MatchRight(e Entity)[]int {
    var i = make([]int, 0, len(s.prox))
    for _, p := range s.prox {
        if IsTypeOf(e, p.stype) || e.Equal(p.stype ) {
            if p.pos > 0 {
                i = append(i, p.pos)
            }
        }
    }
    return i
}

func (s*surrounding)MatchLeft(e Entity)[]int {
    var i = make([]int, 0, len(s.prox))
    for _, p := range s.prox {
        if IsTypeOf(e, p.stype) || e.Equal(p.stype ) {
            if p.pos < 0 {
                i = append(i, p.pos)
            }
        }
    }
    return i
}

type RecognitionTree struct {
    roots []*entangledRecognitionNode

}

func (rt*RecognitionTree) Add(re *RecognizedEntity) {

}
