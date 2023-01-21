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
func (re *RecognizedEntity) directTypeOf() []*EntityType{
    return re.entity.directTypeOf()
}

/*
Re implements Entity
*/
func (re *RecognizedEntity) Equal(other utils.Equalable) bool {
    other_re, ok := other.(*RecognizedEntity)
    if !ok {return false}

    return re.entity.Equal(other_re.entity)
}

func (re *RecognizedEntity) GetNumber() int8 {return re.entity.GetNumber()}

func (re *RecognizedEntity) GetGender() int8 {return re.entity.GetGender()}

func (re *RecognizedEntity) GetGrammarGroup() group {
    switch  re.entity.(type) {
    case *EntityType:
        return re.entity.(*EntityType).grammar_group
    default:
        return re.entity.directTypeOf()[0].grammar_group
    }
}

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

type proposer interface {
    proposeOptions(Word, *CTX) []*RecognizedEntity
    computeP(RecognizedEntity, *CTX, ...RecognizedEntity) float64
}

func (re*RecognizedEntity) addLink(el Link, destination Entity) (int, error){
    return re.entity.addLink(el, destination)
}

func (re*RecognizedEntity) removeLink(index int) {
    re.entity.removeLink(index)
}

func (re*RecognizedEntity) hasLink (lt Link, destination Entity) bool {
    return re.entity.hasLink(lt, destination)
}

func (re*RecognizedEntity) load (s[]string,
                                 m1 map[string]group, m2 map[int]Entity) {
    re.entity.load(s, m1, m2)
}

func (re*RecognizedEntity) store(i int, m map[Entity]int, b *[]byte) string {
    return re.entity.store(i, m, b)
}
