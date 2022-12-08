package ceres

type entangledRecognitionNode struct {
    possibilities map[*surrounding][]*recognitionNode
    Content  *RecognizedEntity
    parent   entangledRecognitionParent
}

func (ern*entangledRecognitionNode)score() float64 {
    var maxScore float64

    nodes := ern.nodes()
    if len(nodes) == 0 {return 1.}

    for _, node := range nodes {
        var score float64 = node.Surround.coherence
        for _, child := range node.ChildMap {
            score *= child.score()
        }
        for _, prox := range node.Surround.prox {
            if _, ok := node.ChildMap[prox.pos]; !ok {
                score *= prox.pMissing
            }
        }
        //DEBUG_PRINTOUT//fmt.Printf("score %v at node %v[%v] with %v children\n",
            //DEBUG_PRINTOUT//score, ern.Content.s, node.Surround.String(), len(node.ChildMap))
        if score > maxScore {
            maxScore = score
        }
    }
    //DEBUG_PRINTOUT//fmt.Println("maxScore for", ern.Content.s, ": ", maxScore)
    return maxScore
}

func (ern*entangledRecognitionNode)copy() *entangledRecognitionNode {
    copy := new(entangledRecognitionNode)
    copy.parent = ern.parent
    copy.Content = ern.Content
    copy.possibilities = make(map[*surrounding][]*recognitionNode)
    for s, slice := range ern.possibilities {
        copy.possibilities[s] = make([]*recognitionNode, len(slice))
        for i, ptr := range slice {
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
    if ern == nil {panic("cannot tolerate removing node from nil")}
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

    for key, array := range ern.possibilities {
        for i, element := range array {
            if element == nil {
                element = new(recognitionNode)
                element.Surround = key
                element.tangler = ern
                element.ChildMap = make(map[int]*entangledRecognitionNode)
                ern.possibilities[key][i] = element
            }
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
    if rn == nil { return nil }

    copy := new(recognitionNode)
    copy.tangler = rn.tangler
    copy.Surround = rn.Surround
    copy.ChildMap = make(map[int]*entangledRecognitionNode)
    for k, v := range rn.ChildMap {
        v_copy := v.copy()
        copy.ChildMap[k] = v_copy
        v_copy.parent = copy // we fix the links.
    }
    //DEBUG_PRINTOUT//fmt.Printf("Copied %v into %v (nodes of %s). We also copied the %v children of this node (into %v children)\n",
        //DEBUG_PRINTOUT//rn, copy, rn.tangler.Content.s, len(rn.ChildMap), len(copy.ChildMap))
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
    if rn == nil  {panic("should not ever seek an empty recognition node")}

    //DEBUG_PRINTOUT//fmt.Println("We are trying to add", re.s, "onto",
        //DEBUG_PRINTOUT//rn.tangler.Content.s, "which has", len(rn.tangler.nodes()),
        //DEBUG_PRINTOUT//"direct child nodes and", len(rn.ChildMap), "children on this specific node")

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
        //DEBUG_PRINTOUT//fmt.Println("\tWe failed")
        return nil, false
    case 1:
        if _, ok := rn.ChildMap[possible[0]]; !ok {
            //DEBUG_PRINTOUT//fmt.Println("\t", rn.ChildMap)
            rn.ChildMap[possible[0]] = &entangledRecognitionNode{Content:re,
                                            parent:rn,
                                            possibilities:make(map[*surrounding][]*recognitionNode)}
            //DEBUG_PRINTOUT//fmt.Println("\t", rn.ChildMap)
            //DEBUG_PRINTOUT//fmt.Printf("\tWe succeeded @%p. It has now %v direct children\n", rn.ChildMap[possible[0]], len(rn.ChildMap))
            answer = []*entangledRecognitionNode{rn.ChildMap[possible[0]]}
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

/*
Considers adding the recognized entities to all possible interpretations of the ern

Also tries to canibalize descendants directly between ern's content and re (in the sentence) as children of re

Acts recursively. No matter whether we can add to this ern, tries adding to all its descendants.
*/
func (ern*entangledRecognitionNode) Add(re *RecognizedEntity) bool {
    var did_add bool = false
    //DEBUG_PRINTOUT//fmt.Println("Trying to see if", re.s, "can be added to", ern.Content.s)

    ern_copy_nodes := ern.copy().nodes()

    for i, possibility := range ern.nodes() {
        if possibility == nil  {
            //DEBUG_PRINTOUT//fmt.Println(ern.nodes(), ern.Content.s)
            panic("should not get a node that is nil")
        }
        if erns, ok := possibility.Add(re); ok {
            for _, child := range erns{
                ern_copy_nodes[i].try_unspooling_children(child)

            }
            did_add = true
        }
        if right_children := possibility.children_on_the_right(); len(right_children) > 0 {
            rightmost_child := right_children[len(right_children)-1]
            if rightmost_child.Add(re) {
                did_add = true
            }
        }
    }
    //DEBUG_PRINTOUT//if did_add {fmt.Println("We could add", re.s, "to", ern.Content.s)}
    return did_add
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
    master *EntangledRecognitionForest
}

/*
Makes a copy of the tree at a different address and returns it.
Does not attempt to add that tree to this one's forest, but does copy the forest pointer into it.
*/
func (rt*RecognitionTree) copy() *RecognitionTree {
    var roots []*entangledRecognitionNode = make([]*entangledRecognitionNode,
                                                len(rt.roots))
    new := new(RecognitionTree)
    new.roots = roots
    for i, root := range rt.roots {
        new.roots[i] = root.copy()
    }
    new.master = rt.master
    return new
}

func (rt*RecognitionTree) Add(re *RecognizedEntity) {
    // is the tree empty ?
    if len(rt.roots) == 0 {
        // re is the root
        rt.roots = []*entangledRecognitionNode{&entangledRecognitionNode{Content:re, parent:rt,
                                            possibilities:make(map[*surrounding][]*recognitionNode)}}
        //DEBUG_PRINTOUT//fmt.Println("tree now has a root :", re.s)
        for _, s := range re.surroundings().surr {
            rt.roots[0].possibilities[s] = []*recognitionNode{
                &recognitionNode{tangler:rt.roots[0],
                Surround:s,
                ChildMap:make(map[int]*entangledRecognitionNode)}}
        }
        //DEBUG_PRINTOUT//fmt.Println(rt.roots)
        return
    }

    /*
    Could the roots be children of re ?
    */
    re_as_hyper_root_copy := rt.copy()
    re_as_ern := &entangledRecognitionNode{Content:re,
                        parent:re_as_hyper_root_copy,
                        possibilities:make(map[*surrounding][]*recognitionNode)}

    cur_tree := rt.copy()
    var debased_roots int = 0

    for _, s := range re.surroundings().surr {
        rn := new(recognitionNode)
        rn.ChildMap = make(map[int]*entangledRecognitionNode)
        rn.Surround = s
        rn.tangler = re_as_ern
        for i := len(cur_tree.roots) - 1; i>=0; i-- {

            //DEBUG_PRINTOUT//fmt.Println("rn for", re.s, ":", rn)
            //DEBUG_PRINTOUT//fmt.Println("\t", cur_tree.roots[i] == nil, cur_tree.roots)
            if array, ok := rn.add(cur_tree.roots[i].Content, true); ok {
                for _, head := range array {
                    head.possibilities = cur_tree.roots[i].possibilities
                    debased_roots ++
                }
            }
        }
    }
    if debased_roots > 0 {
        rt.master.append(re_as_hyper_root_copy)
    }

    /*
    Could re make for a good root ?
    */
    //DEBUG_PRINTOUT//fmt.Println("Trying to add", re.s, "as a new root in a new tree")
    copy := rt.copy()
    copy.roots = append(copy.roots, &entangledRecognitionNode{Content:re,
                                            parent:copy,
                                            possibilities:make(map[*surrounding][]*recognitionNode)})
    rt.master.append(copy) // added without check because element was added !

    /*
     Can re be a descendant of the rightmost root ?
    */
    copy = rt.copy()
    if copy.roots[len(rt.roots)-1].Add(re) {
        rt.master.append(copy)
        //DEBUG_PRINTOUT//fmt.Println("DBG :", len(copy.roots[0].nodes()[0].ChildMap))
    }

    rt.master.remove(rt) //this tree is no longer up to date. Remove it
}

func (rt*RecognitionTree) remove(ern*entangledRecognitionNode) {
    for i, root := range rt.roots {
        if root == ern {
            if i == len(rt.roots) - 1 {
                rt.roots = rt.roots[:i]
            } else {
                rt.roots = append(rt.roots[:i], rt.roots[i+1:]...)
            }
            if len(rt.roots) == 0 {
                rt.master.remove(rt)
            }
            return
        }
    }
}

func (rt*RecognitionTree) score() float64 {
    var score float64 = 1.

    for _, root := range rt.roots {
        sc := root.score()
        score *= 0.5 * sc
        //DEBUG_PRINTOUT//fmt.Println("score for root", root.Content.s, ":", sc)
    }

    return score
}

type entangledRecognitionParent interface {
    remove(*entangledRecognitionNode)
}

type EntangledRecognitionForest []*RecognitionTree

// FIXME: remove this
const TreeSep string = "=============================================="

func (erf*EntangledRecognitionForest) Add (re*RecognizedEntity){
    if len(*erf) == 0 {
        *erf = EntangledRecognitionForest{new(RecognitionTree)}
    }

    previously_existing_trees := make([]*RecognitionTree, len(*erf))
    //DEBUG_PRINTOUT//previous_length := len(previously_existing_trees)
    copy(previously_existing_trees, *erf)
    for _, tree := range previously_existing_trees {
        //INFO_PRINTOUT//fmt.Println(TreeSep, "Working on tree n°", i+1, "out of", len(previously_existing_trees), TreeSep)
        tree.Add(re)
        //INFO_PRINTOUT//fmt.Println(TreeSep, "Created", len(*erf) - previous_length , "new trees", TreeSep)
        //INFO_PRINTOUT//previous_length = len(*erf)
    }
}

func (erf*EntangledRecognitionForest) append(rt *RecognitionTree) {
    if erf == nil {
        panic("Cannot append tree on non-existant forest")
    }

    rt.master = erf
    //DEBUG_PRINTOUT//fmt.Println("Append a tree to the forest at position", len(*erf))
    *erf = append(*erf, rt)
}

func (erf*EntangledRecognitionForest) remove(rt *RecognitionTree) {
    for i, cur_rt := range *erf {
        if rt == cur_rt {
            //DEBUG_PRINTOUT//fmt.Printf("Removing tree n°%v out of %v\n", i+1, len(*erf))
            if i == len(*erf) - 1 {
                *erf = (*erf)[:i]
            } else {
                *erf = append((*erf)[:i], (*erf)[i+1:]...)
            }
            return
        }
    }
}

func (erf*EntangledRecognitionForest) bestTree() (*RecognitionTree, float64) {
    var max_idx int
    var score float64
    for i, tree := range *erf {
        sc := tree.score()
        if sc > score {
            max_idx = i
            score = sc
        }
    }
    return (*erf)[max_idx], score
}
