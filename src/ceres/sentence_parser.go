package ceres

import (
	"CERES/src/utils"
	"fmt"
	"math"
	re "regexp"
	"strings"
	"sync"
	"time"
)

var RegexpToken *re.Regexp = re.MustCompile(utils.TokenPattern)

/*
Parses the sentence and returns an array of recognized entities as well as the score of that possibility.
*/
func (c *CERES) ParseSentence(sentence string) ([]*RecognizedEntity, float64) {
	split_sentence := c.SplitSentence(sentence)
	var possibilities = new([]ceres_possibility_scored)

	for _, word := range split_sentence {
		fmt.Println("============[OPERATING] on", word, "===================")

		options := c.allOptions(Word(word))
		c.updatePossibilities(possibilities, options)
		res, sc := getBestPossibility(possibilities)
		cps := ceres_possibility_scored{res: res, score: sc}
		fmt.Printf("%s, %p\n", cps.ToString(), cps.res[len(cps.res)-1])
	}

	return getBestPossibility(possibilities)
}

func getBestPossibility(possibilities *[]ceres_possibility_scored) ([]*RecognizedEntity, float64) {
	var best_res []*RecognizedEntity
	var best_score float64 = math.Inf(-1)

	for _, possibilities := range *possibilities {
		if possibilities.score > best_score {
			best_score = possibilities.score
			best_res = possibilities.res
		}
	}

	return best_res, best_score
}

type ceres_possibility_scored struct {
	res   []*RecognizedEntity
	score float64
}

func (cps *ceres_possibility_scored) copy() *ceres_possibility_scored {
	cps2 := new(ceres_possibility_scored)

	cps2.score = cps.score
	cps2.res = make([]*RecognizedEntity,
		len(cps.res),
		len(cps.res)+1)
	copy(cps2.res, cps.res)

	return cps2
}

func (cps *ceres_possibility_scored) ToString() string {

	var s string = "["
	for _, re := range cps.res {
		s += fmt.Sprintf("(\"%s\", \"%s\")", re.s, re.proposer.name())
	}
	return s + "] : " + fmt.Sprintf("%f", cps.score)
}

func (c *CERES) updatePossibilities(possibilities *[]ceres_possibility_scored,
	options []*RecognizedEntity) {

	if possibilities == nil {
		panic("Needs the possibilities to be non nil")
	} else if len(*possibilities) == 0 {
		*possibilities = make([]ceres_possibility_scored, 1)
		(*possibilities)[0].score = 1.
	}

	if len(options) == 0 {
		panic("no option found")
	}

	var counter_rwm *sync.RWMutex = new(sync.RWMutex)

	results_getter := make(chan ceres_possibility_scored)
	var counter *int = new(int)
	*counter = len(*possibilities) * len(options)

	var new_possibilities []ceres_possibility_scored = make([]ceres_possibility_scored, *counter)

	for _, possibility := range *possibilities {
		for _, found_entity := range options {
			go c.merge(*possibility.copy(), found_entity,
				results_getter, counter, counter_rwm)
		}
	}

	getCounter := func() int {
		counter_rwm.RLock()
		defer counter_rwm.RUnlock()
		return *counter
	}

	nposs_counter := 0
	for getCounter() > 0 {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Waited a bit")
			continue
		case poss := <-results_getter:
			new_possibilities[nposs_counter] = poss
			nposs_counter++
		}
	}
	*possibilities = new_possibilities
}

func (c *CERES) merge(poss ceres_possibility_scored,
	fe *RecognizedEntity,
	result_getter chan ceres_possibility_scored,
	counter *int, counter_rwm *sync.RWMutex) {

	fmt.Printf("%s * %f\n", poss.ToString(),
		fe.proposer.computeP(fe, c.ctx))
	poss.res = append(poss.res, fe)
	poss.score *= fe.proposer.computeP(fe, c.ctx)
	fmt.Printf("[%s < %s | %3f | %s ] on %s @%p, %p \n", fe.s, fe.proposer.name(),
		poss.score, poss.res[len(poss.res)-1].proposer.name(),
		poss.ToString(), &poss, fe)
	// this operation is actually the time sensitive one, the one we want parallelize.

	//println("waiting on lock")

	defer func() {
		counter_rwm.Lock()
		(*counter)--
		counter_rwm.Unlock()
	}()

	select {
	case result_getter <- poss:
		fmt.Printf("SENT %s %p\n", poss.ToString(), &poss)
	case <-time.After(3 * time.Second):
	}
}

/*
Splits a sentence along the whitespace and the apostrophe characters
*/
func (c *CERES) SplitSentence(sentence string) []Word {
	var seps = " \U000e0027'  -−＇‾ʼ՚ߴߵ\"«»,"

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
