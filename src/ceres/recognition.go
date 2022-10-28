package ceres

import (
    "CERES/src/utils"
    re "regexp"
)

type RecognizedEntity struct {
    entity      Entity
    possessive  bool
    attribute   bool
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

func MakeRecognizedEntity(e Entity, p bool, a bool) RecognizedEntity{
    return RecognizedEntity{entity:e, possessive:p, attribute:a}
}
