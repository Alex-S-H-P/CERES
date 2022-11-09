package ceres

import (
    "CERES/src/utils"
    re "regexp"

    "fmt"
)

type RecognizedEntity struct {
    entity      Entity
    possessive  bool
    attribute   bool
    proposer    proposer
    s string
}

/*
RE implements Entity
*/
func (re *RecognizedEntity) Initialize()  {}

/*
RE implements Entity
*/
func (re *RecognizedEntity) directTypeOf() *EntityType{
    return re.entity.directTypeOf()
}

/*
Re implements Entity
*/
func (re *RecognizedEntity) Equal(other utils.Equalable) bool {
    return re.entity.Equal(other)
}

func (re *RecognizedEntity) GetNumber() int8 {return re.entity.GetNumber()}

func (re *RecognizedEntity) GetGender() int8 {return re.entity.GetGender()}

type tokenT uint8

const (
    TOKEN_TYPE_WORD tokenT = 0
    TOKEN_TYPE_CURR tokenT = 2
    TOKEN_TYPE_NUMB tokenT = 4
    TOKEN_TYPE_UNKN tokenT = 6
    TOKEN_TYPE_INTN tokenT = 8
    TOKEN_TYPE_EOS  tokenT = 10
    TOKEN_TYPE_PRIC tokenT = 12
)

func recognizeType(token string) tokenT{
    if ok, _ := re.MatchString(utils.WordPattern, token); ok{
        return TOKEN_TYPE_WORD
    } else if  ok, _ := re.MatchString(utils.PricePattern, token); ok {
        return TOKEN_TYPE_PRIC
    } else if ok, _ := re.MatchString(utils.CurrencyPattern, token); ok {
        return TOKEN_TYPE_CURR
    } else if ok, _ := re.MatchString(utils.NumberPattern, token); ok {
        return TOKEN_TYPE_NUMB
    }
    return TOKEN_TYPE_CURR
}

func MakeRecognizedEntity(e Entity, p bool, a bool, pr proposer, s string) RecognizedEntity{
    return RecognizedEntity{entity:e, possessive:p, attribute:a, proposer:pr, s:s}
}

func (re*RecognizedEntity) surroundings() *surroundingList {
    switch re.entity.(type) {
    case *EntityType:
        return &re.entity.(*EntityType).surroundingList
    default:
        return &re.entity.directTypeOf().surroundingList
    }
}

type proposer interface {
    proposeOptions(Word, *CTX) []RecognizedEntity
    computeP(RecognizedEntity, *CTX, ...RecognizedEntity) float64
}

type RecognitionNode struct {
    Parent *RecognitionNode
    Content *RecognizedEntity

    LeftChildren  []*RecognitionNode // these one should be for the possibility to be considered.
    RightChildren []*RecognitionNode // these one shouldn't
}

func (rn*RecognitionNode) NbChildren() int {
    s := 1
    for _, child := range rn.LeftChildren {
        s += child.NbChildren()
    }
    for _, child := range rn.RightChildren {
        s += child.NbChildren()
    }
    return s
}

func (rn*RecognitionNode) NbOnTheRight() int {
    var s int
    for _, child := range rn.RightChildren {
        s += child.NbOnTheRight() + 1
    }
    return s
}

type RecognitionTree struct {
    Root []*RecognitionNode

    contents []*RecognizedEntity
    curCoherence float64
}


/*
tries to find a tree structure that matches having this node as its root.
Since this function is used recursively,
we pass along a beta value to know what to miss
*/
func (rn *RecognitionNode)shape(current*RecognizedEntity,
    before, after []*RecognizedEntity, beta float64) *RecognitionTree{

    var answerTree = new(RecognitionTree)

    var BestRnCpy *RecognitionNode = nil
    var BestRnCpyScore float64 = -1.


    for _, curSurrounding := range current.surroundings().surr {

        fmt.Println(">> there are", len(curSurrounding.prox), "tokens to be parsed")
        var start_index_on_lookup, end_index_on_lookup int = 0, len(before)
        var is_looking_before bool = true

        if curSurrounding.coherence < beta {
            continue
        }
        // copying rn
        var CurRnCpy *RecognitionNode = new(RecognitionNode)
        var CurRnCpyScore float64 = 1.
        CurRnCpy.Content = rn.Content
        CurRnCpy.Parent = rn.Parent
        fmt.Println("current copy :", CurRnCpy.Content.s, CurRnCpy.LeftChildren, CurRnCpy.RightChildren)

        // filling the copy with the analysis of the
        for _, prox := range curSurrounding.prox {
            var prox_beta float64 = prox.pMissing
            var bestSubTree *RecognitionTree
            var offset int

            fmt.Println("SEARCHING a match of", prox, "for surrounding", curSurrounding, " : ", current.s)

            var ents []*RecognizedEntity = before
            if ! is_looking_before {
                ents = after
            } else if prox.pos > 0 {
                is_looking_before = false
                ents = after
                start_index_on_lookup = 0
                end_index_on_lookup = len(after)
            }

            for i, re := range ents[start_index_on_lookup:end_index_on_lookup]{
                fmt.Println("Looking at", re.s)
                if IsTypeOf(re, prox.stype) || re.Equal(prox.stype){
                    fmt.Println(re.entity, "DOES match.", prox.stype, "(", re.s, ")")
                    childNode := new(RecognitionNode)
                    childNode.Parent = rn
                    childNode.Content = re
                    if rn != nil {
                        if rn.Content == nil {
                            fmt.Println("Making child node from \"\" on", re.s)
                        } else {
                            fmt.Println("Making child node from", rn.Content.s, "on", re.s)
                        }
                    } else {
                        fmt.Println("Making child node from", nil, "on", re.s)
                    }

                    var subtree *RecognitionTree

                    if i+start_index_on_lookup < end_index_on_lookup - 1 {
                        subtree = childNode.shape(re,
                            ents[start_index_on_lookup:i+start_index_on_lookup],
                            ents[i+1+start_index_on_lookup:end_index_on_lookup],
                            prox_beta)
                    } else {
                        subtree = childNode.shape(re,
                            ents[start_index_on_lookup:i+start_index_on_lookup],
                            nil,
                            prox_beta)
                    }

                    fmt.Println("beta :", beta, "|  P :", subtree.curCoherence)

                    if prox_beta < subtree.curCoherence {
                        prox_beta = subtree.curCoherence
                        bestSubTree = subtree
                        offset = subtree.Root[0].NbOnTheRight() + 1
                    }
                } else {
                    fmt.Println(re.entity, "does not match.", prox.stype)
                }
            } // best bestSubTree found
            start_index_on_lookup += offset

            if bestSubTree !=  nil {
                if is_looking_before {
                    CurRnCpy.LeftChildren = append(CurRnCpy.LeftChildren,
                        bestSubTree.Root[0])
                } else {
                    CurRnCpy.RightChildren = append(CurRnCpy.RightChildren,
                        bestSubTree.Root[0])
                }
                CurRnCpyScore *= bestSubTree.curCoherence
            } else {
                CurRnCpyScore *= prox.pMissing
            }
            fmt.Println("cur:", len(CurRnCpy.LeftChildren), len(CurRnCpy.RightChildren),
                        cap(CurRnCpy.LeftChildren), cap(CurRnCpy.RightChildren))
        }

        // maximise surrounding coherence
        if BestRnCpyScore < CurRnCpyScore {
            BestRnCpy = CurRnCpy
            BestRnCpyScore = CurRnCpyScore
        }
        fmt.Println("best:", len(BestRnCpy.LeftChildren), len(BestRnCpy.RightChildren),
                    cap(BestRnCpy.LeftChildren), cap(BestRnCpy.RightChildren))
    }

    if BestRnCpy != nil {
        fmt.Println("best copy :", BestRnCpy.Content.s, BestRnCpy.LeftChildren, BestRnCpy.RightChildren)
        fmt.Println(len(BestRnCpy.LeftChildren), len(BestRnCpy.RightChildren),
                    cap(BestRnCpy.LeftChildren), cap(BestRnCpy.RightChildren))
        rn.LeftChildren = make([]*RecognitionNode, len(BestRnCpy.LeftChildren))
        rn.RightChildren = make([]*RecognitionNode, len(BestRnCpy.RightChildren))//*/
        copy(rn.LeftChildren, BestRnCpy.LeftChildren)
        copy(rn.RightChildren, BestRnCpy.RightChildren)
        fmt.Println("copied into :", rn.Content.s, rn.LeftChildren, rn.RightChildren)
        answerTree.curCoherence = BestRnCpyScore
    } else {
        rn.LeftChildren = nil
        rn.RightChildren = nil
        answerTree.curCoherence = 1.
    }

    answerTree.Root = []*RecognitionNode{rn}
    fmt.Println("node from left/right", rn.LeftChildren, rn.RightChildren)
    return answerTree
}
