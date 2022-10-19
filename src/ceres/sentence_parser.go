package ceres

import (
	"fmt"
	"regexp"
	"sync"
)

var SPLITTER = regexp.MustCompile("[ ,\t']")

type positionTuple struct {
	i,j int
}

type ics_DijkstraPossibility struct {
	// current probability
	curP float64
	// current list of analysed entities.
	collapsedAnalysedSentence map[positionTuple]Entity
	// current list of remaining words to analyse.
	tokens_analysed[]int // cap tokens_analysed == len(tokens)
}

type organizedDijkstraList struct {
	list []*ics_DijkstraPossibility
	rwm  sync.RWMutex
}

func newDijkstraList() *organizedDijkstraList{
	ls := new(organizedDijkstraList)
	ls.list = make([]*ics_DijkstraPossibility, 0, 1024)
	return ls
}

func (odl *organizedDijkstraList)successed() bool {return odl.list[0].completed()}

func (odl *organizedDijkstraList)pop() *ics_DijkstraPossibility{
	if odl.rwm.TryLock() {
		odl.rwm.Lock()
		defer odl.rwm.Unlock()
	}
	a := odl.list[0]
	if len(odl.list) > 1 {
		odl.list = odl.list[1:]
	}
	return a
}

func (odl *organizedDijkstraList)put(possibility *ics_DijkstraPossibility) {
	if odl.rwm.TryLock() {
		odl.rwm.Lock()
		defer odl.rwm.Unlock()
	}
	i := odl.getRank(possibility.curP)
	tmp := odl.list[i:]
	odl.list = append(odl.list[:i], possibility)
	odl.list = append(odl.list, tmp...)
}

func (odl *organizedDijkstraList)getRank(P float64) int {
	a, b := 0, len(odl.list)
	for a - b > 1 {
		c := (b-a)/2 + a
		if P > odl.list[c].curP {
			// we are in descending order, therefore i_P should be before c
			b = c
		} else {
			a = c
		}
	}
	return b
}

func (idp *ics_DijkstraPossibility) completed() bool {
	return cap(idp.tokens_analysed) == len(idp.tokens_analysed)
}

func (idp *ics_DijkstraPossibility)collapse() []Entity {
	array := make([]Entity, 0, len(idp.collapsedAnalysedSentence))
	var i int

	for {
		for pos, e := range idp.collapsedAnalysedSentence {
			if pos.i == i {
				i = pos.j + 1
				array = append(array, e)
			}
		}
		if i >= len(idp.collapsedAnalysedSentence){
			return array
		}
	}
}

func (c *CERES) ParseSentence(sentence string) []Entity {
	var tokens []Word

	pre_tokens := SPLITTER.Split(sentence, len(sentence))

	tokens = make([]Word, 0, len(pre_tokens))

	for _, token := range pre_tokens {
		if len(token) > 0 {
			tokens = append(tokens, Word(token).Lower())
		}
	}

	// the sentence is split in words
	fmt.Println("Now solving pronouns.")

	var tokens_analysed []int = make([]int, 0, len(tokens))

	// checking for pronouns
	for i, word := range tokens {
		if c.pcs.IsPronoun(word) {
			tokens_analysed = append(tokens_analysed, i)
		}
	}
	var wg *sync.WaitGroup = new(sync.WaitGroup)
	var ls *organizedDijkstraList = newDijkstraList()
	wg.Add(c.sentence_analyser_workers)
	for worker_id := 0; worker_id < c.sentence_analyser_workers; worker_id++ {
		go c.ics.worker_main(ls, wg)
	}

	wg.Done()
	if len(ls.list) <= 0 {
		panic(fmt.Errorf("The dijkstra list is empty ! [%T : %v]", ls, ls))
	}
	return ls.list[0].collapse()
}

func (ics *ICS)worker_main(list *organizedDijkstraList, wg *sync.WaitGroup){
	if wg != nil {
		defer wg.Done() // we are no longer counting
	}

	for {
		list.rwm.Lock()
		if list.successed() || len(list.list) == 0 {
			list.rwm.Unlock()
			return
		}

		// we are now in the list, we can select the first token available
		candidate := list.pop()
		list.rwm.Unlock()
		for _, possibility := range ics.evolve(candidate) {
			list.put(possibility)
		}
	}
}


func (ics *ICS)evolve(cur *ics_DijkstraPossibility) []*ics_DijkstraPossibility {
	// TODO: Do this
	return nil
}
