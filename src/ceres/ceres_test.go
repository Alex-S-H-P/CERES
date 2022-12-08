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

func setSurrounding(parent, child *EntityType, pos int) {
    proxes := parent.surroundingList.surr[0].prox
    token := surroundingToken{stype:child, pos:pos, pMissing:0.}
    proxes = append(proxes, token)
    parent.surroundingList.surr[0].prox = proxes
    if pos > parent.surroundingList.surr[0].maxPos {
        parent.surroundingList.surr[0].maxPos = pos
    }
    if pos < parent.surroundingList.surr[0].minPos {
        parent.surroundingList.surr[0].minPos = pos
    }
}

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

            fmt.Println("----------------Conf : ", root_id, ":",
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

            proxes := ""
            initial := true
            for _, p := range root.surroundingList.surr[0].prox{
                if !initial {
                    proxes += ", "
                } else {initial = false}
                proxes += trmap[p.stype].s
            }
            fmt.Println("\nprox for root :", proxes, "after set up")

            // testing
            var forest = EntangledRecognitionForest(make([]*RecognitionTree, 0, 4))
            var tree = new(RecognitionTree)
            forest.append(tree)
            forest.Add(&Ar)
            forest.Add(&Br)
            forest.Add(&Cr)
            forest.Display()
            fmt.Println("___________________________________________________")
            return
        }
    }
}
