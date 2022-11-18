package ceres

import (
    "CERES/src/utils"
    re "regexp"

    //"fmt"
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
    surr *surrounding
    M map[int]*RecognitionNode
}

func (rn*RecognitionNode) RightChildrenOrdered() []*RecognitionNode {
    s := make([]*RecognitionNode, 0, len(rn.M))

    for i := 0; i < len(rn.surr.prox); i++ {
        if child, ok := rn.M[i]; ok {
            s = append(s, child)
        }
    }

    return s
}

func (rn*RecognitionNode) NbChildren() int {
    s := 1
    for _, child := range rn.M {
        s += child.NbChildren()
    }
    return s
}

func (rn*RecognitionNode) NbOnTheRight() int {
    var s int
    for k, child := range rn.M {
        if k > 0 {
            s += child.NbOnTheRight() + 1
        }
    }
    return s
}

type RecognitionTree struct {
    Root []*RecognitionNode

    contents []*RecognizedEntity
    curCoherence float64
}
