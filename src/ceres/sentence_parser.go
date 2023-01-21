package ceres

import (
	"math"
	"sync"
	"strings"
	re "regexp"
	"CERES/src/utils"
)

var RegexpToken *re.Regexp = re.MustCompile(utils.TokenPattern)

func (c *CERES)ParseSentence(sentence string) ([]RecognizedEntity, float64){
	split_sentence := c.SplitSentence(sentence)
	var possibilities = new([]ceres_possibility_scored)

	for _, word := range split_sentence {
		de := c.ics.entityDictionary[word]

		c.updatePossibilities(possibilities, de, word)
	}

	return getBestPossibility(possibilities)
}


func getBestPossibility(possibilities *[]ceres_possibility_scored) ([]RecognizedEntity, float64) {
	var best_res []RecognizedEntity
	var best_score float64 = math.Inf(-1)

	for _, possibilities := range *possibilities {
		if possibilities.score > best_score {
			best_score = possibilities.score
			best_res   = possibilities.res
		}
	}

	return best_res, best_score
}

type ceres_possibility_scored struct {
	res []RecognizedEntity
	score float64
}

func (c *CERES) updatePossibilities(possibilities *[]ceres_possibility_scored,
	de*DictionaryEntry, w Word) {

	if possibilities == nil {
		panic("Needs the possibilities to be non nil")
	} else if len(*possibilities) == 0 {
		*possibilities = make([]ceres_possibility_scored, 1)
	}


	if de == nil {
		for i, possibility := range (*possibilities) {
			(*possibilities)[i].res = append(possibility.res,
									      	 *c.ucs.proposeOptions(w, c.ctx)[0])

			CYK_PARSE((*possibilities)[i].res, c.grammar)
		}
		return
	}

	var counter *sync.WaitGroup = new(sync.WaitGroup)


	results_getter := make(chan ceres_possibility_scored)
	nbOfFusionsNeeded := len(*possibilities)*len(de.entities)


	var new_possibilities []ceres_possibility_scored = make([]ceres_possibility_scored, nbOfFusionsNeeded)

	counter.Add(nbOfFusionsNeeded)
	for _, possibility := range *possibilities {
		for _, found_entity := range de.entities {
			go merge(possibility, found_entity, results_getter)
		}
	}


	possibilities = &new_possibilities
}

func merge(poss ceres_possibility_scored,
	fe Entity,
	result_getter chan ceres_possibility_scored) {


}

/*
Splits a sentence along the whitespace and the apostrophe characters
*/
func (c *CERES)SplitSentence(sentence string)[]Word {
	var seps=" 󠀧'  -−＇‾ʼ՚ߴߵ\"«»"

	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}

	_strings := strings.FieldsFunc(strings.ToLower(sentence), splitter)
	var words []Word = make([]Word, len(_strings))

	for i, word := range _strings {
		words[i] = Word(word)
	}

	return words
}

func (c *CERES)makeNonWordEntity(token string) RecognizedEntity {
	// TODO handle all cases.
	switch recognizeType(token) {
	case TOKEN_TYPE_PRIC:
		// make a temporary price entity
	case TOKEN_TYPE_CURR:
		// replace with the currency's entity.
	case TOKEN_TYPE_INTN, TOKEN_TYPE_NUMB:
		// make temporary number entity
	case TOKEN_TYPE_UNKN:
		// fetch unknown. If failed make one.
	}

	// FIXME: this code is temporary
	ei := new(EntityInstance)

	return RecognizedEntity{entity:ei,
		possessive:false, attribute:false, proposer:&c.ucs}
}
