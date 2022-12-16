package ceres

import (
    "testing"
    "fmt"
)

func TestCERES_Initialization(t*testing.T){
    ceres := new(CERES)
    ceres.Initialize(1)

    if !ceres.initialized {
        t.Errorf("CERES should have been marked as initialized")
    }

    if !ceres.ics.initialized {
        t.Errorf("ICS should have been initialized")
    }

    if !ceres.pcs.initialized {
        t.Errorf("PCS should have been initialized")
    }

    if ceres.root == nil {
        t.Errorf("CERES should have  ")
    }
}

func TestCeres_TypeAdding(t *testing.T){
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    if len(*c.root.children) != 1 {
        t.Errorf("Created a child for the root, but the root did not care")
    } else if (*c.root.children)[0].directTypeOf() == nil {
        t.Errorf("Created a child that did not hook to its parent correctly")
    }

}

func TestICS_LevelsOfAbstractions(t *testing.T) {
    // setup
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    c.createEntityType("action")
    c.createEntityType("thing")
    thing := (*c.root.children)[2].(*EntityType)
    living_being := c.ics.createEntityType("living being")
    thing.addChild(living_being)
    box := c.ics.createEntityType("box")
    thing.addChild(box)
    shmoop := c.ics.createEntityInstance("shmoop", thing)
    John := c.ics.createEntityInstance("john", living_being)

    if i, e := LevelsOfAbstractionBetween(thing, c.root); i != -1 || e != nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, -1, e, nil)
    }

    if i, e := LevelsOfAbstractionBetween(c.root, thing); i != 1 || e != nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, 1, e, nil)
    }

    if i, e := LevelsOfAbstractionBetween(living_being, c.root); i != -2 || e != nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, -2, e, nil)
    }

    if i, e := LevelsOfAbstractionBetween(c.root, living_being); i != 2 || e != nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, 2, e, nil)
    }

    if i, e := LevelsOfAbstractionBetween(living_being, box);  e == nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, 4, e, "non-nil error")
    }

    if i, e := LevelsOfAbstractionBetween(shmoop, John); e == nil {
        t.Errorf("Unexpected result : (%v instead of %v, %v instead of %v)",
            i, 0, e, "non-nil error")
    }
}

func TestIsTypeOf(t *testing.T) {
    // setup
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    c.createEntityType("action")
    c.createEntityType("thing")
    thing := (*c.root.children)[2].(*EntityType)
    living_being := c.ics.createEntityType("living being")
    thing.addChild(living_being)
    box := c.ics.createEntityType("box")
    thing.addChild(box)
    John := c.ics.createEntityInstance("John", living_being)

    if !IsTypeOf(John, living_being) {
        t.Errorf("IsTypeOf failed")
    }
    if IsTypeOf(living_being, living_being) {
        t.Errorf("Entities should not match with themselves")
    }
}

func TestSentenceParser(t *testing.T){
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    c.createEntityType("action")
    c.createEntityType("thing")

    c.ParseSentence("hello, I am alexandre")
}

const fileSpace = "test.ceres"

func TestCERES_Saving(t *testing.T) {
    c := new(CERES)
    c.Initialize(1)
    c.createEntityType("caracteristic")
    c.createEntityType("action")
    c.createEntityType("thing")
    c.save(fileSpace)

    c2 := new(CERES)
    c2.Initialize(1)
    e := c2.load(fileSpace)
    if e != nil {
        t.Error(e)
    }
    if len(*c.root.children) != len(*c2.root.children) {
        t.Fail()
    }
}
/*
func TestShaper(t *testing.T) {
    var At, Bt, Ct *EntityType = new(EntityType),
        new(EntityType), new(EntityType)
    var Ar, Br, Cr RecognizedEntity = MakeRecognizedEntity(At, false, false, nil, "A"),
        MakeRecognizedEntity(Bt, false, false, nil, "B"),
        MakeRecognizedEntity(Ct, false, false, nil, "C")

    var recog = [3]*EntityType{At, Bt, Ct}
    var recognizedEntities = [...]RecognizedEntity{Ar, Br, Cr}
    var trmap = map[*EntityType]RecognizedEntity{At:Ar, Bt:Br, Ct:Cr}

    var idTransformation = [3][3]int {[3]int{0, 2, 1}, [3]int{2, 1, 0}, [3]int{1, 0, 2}}

    for root_id, root := range recog {
        for desc_id, desc := range recog {

            fmt.Println("\n----------------Conf : ", root_id, ":",
                desc_id, "----------------")
            fmt.Println("root is :", recognizedEntities[root_id].s)

            // reboots the surroundingList
            for _, rebootable := range recog {
                rebootable.surroundingList = surroundingList{surr : make([]*surrounding, 1)}
                rebootable.surroundingList.surr[0] = new(surrounding)
                rebootable.surroundingList.surr[0].coherence = 1
            }
            // setting the tree-like structure
            if desc_id == root_id {
                for child_id, child := range recog {
                    if child_id == desc_id {
                        continue
                    }
                    fmt.Println("layer 1 has :", recognizedEntities[child_id].s)
                    setSurrounding(root, child, child_id - desc_id)
                }
            } else {
                middle_id := idTransformation[root_id][desc_id]
                middle := recog[middle_id]
                fmt.Println("layer 1 has :", recognizedEntities[middle_id].s)
                fmt.Println("layer 2 has :", recognizedEntities[desc_id].s)

                setSurrounding(root, middle, middle_id - root_id)
                setSurrounding(middle, desc, desc_id - middle_id)

                if desc_id + middle_id == 2 {
                    continue
                }
            }

            if abs(root_id - desc_id) == 1 {
                //the root cannot be next to the lowest descendant
                fmt.Println("\nConfiguration skipped due to children crossing")
                continue
            }

            proxes := ""
            initial := true
            for _, p := range root.surroundingList.surr[0].prox{
                if !initial {
                    proxes += ", "
                } else {initial = false}
                proxes += trmap[p.stype].s
            }

            // testing
            var forest = EntangledRecognitionForest(make([]*RecognitionTree, 0, 4))
            var tree = new(RecognitionTree)
            forest.append(tree)
            forest.Add(&Ar)
            forest.Add(&Br)
            forest.Add(&Cr)
            tree, score := forest.bestTree()
            tree.Display()
            fmt.Println(score)

            if score != 0.5 || len(tree.roots) != 1 {
                t.Errorf("The tree should only have one root, it has %v (via score) or %v (via counting) roots instead", int(.5 + .5/score), len(tree.roots))
            } else {
                n := tree.roots[0].nodes()[0]
                switch {
                case len(n.ChildMap) == 1 || root_id != desc_id:
                case len(n.ChildMap) == 2 || root_id == desc_id:
                default:
                    var desired int = 1
                    if root_id == desc_id {desired = 2}

                    t.Errorf("The root does not have the right amount of children %v instead of %v",
                        len(n.ChildMap), desired)
                }
            }


        }
    }
}
*/

func TestRuleMatcher(t*testing.T) {
    r := ruleString("A -> A B")
    if g, b := r.matches("A", "B"); b {
        if g != "A" {
            t.Errorf("Did not match to the right group. Matched to \"%s\" instead of \"A\"", string(g))
        }
    } else {
        t.Errorf("Did not match where there should be a match")
    }
}

func TestCYK(t*testing.T){
    var At, Bt, Ct *EntityType = new(EntityType),
        new(EntityType), new(EntityType)
    var Ar, Br, Cr RecognizedEntity = MakeRecognizedEntity(At, false, false, nil, "A"),
        MakeRecognizedEntity(Bt, false, false, nil, "B"),
        MakeRecognizedEntity(Ct, false, false, nil, "C")

    var g grammar = grammar{rules:[]rule{ruleString("E -> B C"),
                                         ruleString("D -> A E"),
                                         ruleString("B -> B"),
                                         ruleString("C -> C"),
                                         ruleString("A -> A")}}

    table := CYK_PARSE([]RecognizedEntity{Ar, Br, Cr}, g)

    fmt.Println("                ")
    fmt.Println(table.rslt())
    for j := 0; j<3; j++ {
        if len(table[j][j].assignments) != 1 {
            line := "["
            sep := ""
            for _, assignment := range table[j][j].assignments{
                line += fmt.Sprintf("%s \"%s\"", sep, assignment)
                sep = ","
            }
            fmt.Println(line, "]")
            t.Errorf("Too many assignments (%v instead of 1) in the case at location (%v, %v)",
                        len(table[j][j].assignments), j, j)
        }
    }
    if len(table.rslt().assignments) < 1 {
        t.Errorf("Could not get an answer")
    } else if len(table.rslt().assignments) > 1 {
        t.Errorf("Got too many answers (%v instead of 1)", len(table.rslt().assignments))
    }

    table.display_tree(0)
}
