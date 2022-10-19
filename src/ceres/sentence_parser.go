package ceres

import (
	"fmt"
	"strings"
	"regexp"
	"time"
	"sync"

	"CERES/src/utils"
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
	// token groups that have not yet been analysed.
	tokens_toAnalyse [][]Word
}


// collapses the collapsedAnalysedSentence map into an ordered list of entities
func (idp *ics_DijkstraPossibility)collapse() []Entity {
	var reslt []Entity = make([]Entity, 0, len(idp.collapsedAnalysedSentence))
	var i int
	for i < len(idp.collapsedAnalysedSentence) {
		for k, v := range idp.collapsedAnalysedSentence {
			if k.i == i {
				reslt = append(reslt, v)
				i = k.j + 1
				continue
			}
		}
	}
	return reslt
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

func (idp *ics_DijkstraPossibility) completed() bool { return len(idp.tokens_toAnalyse) == 0 }

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
	candidate := c.ics.dijkstra_main(ls, c.sentence_analyser_workers, wg)
	return candidate.collapse()
}

func (ics*ICS) dijkstra_main(list *organizedDijkstraList, workers int,
	wg *sync.WaitGroup) *ics_DijkstraPossibility{
	// setup of all workers
	handle_chan := make(chan *ics_DijkstraPossibility, 8*workers)
	counter_chans := make([]chan *sync.WaitGroup, 0, workers)
	for w_id := 0; w_id < workers; w_id++ {
		counter_chan := make(chan *sync.WaitGroup)
		counter_chans = append(counter_chans, counter_chan)
		go ics.worker_main(list, wg, handle_chan, counter_chan )
	}

	for {
		list.rwm.Lock()
		candidate := list.pop()
		list.rwm.Unlock()
		if candidate.completed() {
			// this candidate is the best one
			return candidate
		} else  {
			// we inform all workers
			counter := new(sync.WaitGroup)
			counter.Add(workers)
			for _, channel := range counter_chans {
				channel <- counter
			}
			// we evolve the possibility
			ics.evolve(candidate, handle_chan)
			counter.Wait()
		}
	}
}

func (ics *ICS)worker_main(list *organizedDijkstraList, wg *sync.WaitGroup,
	handle_chan chan *ics_DijkstraPossibility, counter_chan chan *sync.WaitGroup){
	if wg != nil {
		defer wg.Done() // we are no longer counting
	}
MAINLOOP:
	for {
		counter := <- counter_chan
		select {
		case poss := <- handle_chan:
			poss.curP = ics.estimateP(poss)
			list.put(poss)
		case <-time.After(5*time.Second):
			counter.Done()
			continue MAINLOOP
		}
	}
}


func (ics *ICS)evolve(cur *ics_DijkstraPossibility,
	handler chan *ics_DijkstraPossibility) []*ics_DijkstraPossibility {
	// generating all possible words
	cur_tokens := cur.tokens_toAnalyse[0]
	var word Word = cur_tokens[0]
	var toAnalyse []Word = make([]Word, 0, len(cur_tokens))
	for i := 0; i <len(cur_tokens)-1; i++ {
		toAnalyse = append(toAnalyse, word)
		word += " " + cur_tokens[i]
	}
	toAnalyse = append(toAnalyse, word)

	var i int
	searcher:
	for i = 0; i<cap(cur.tokens_analysed); i++ {
		for _,j := range cur.tokens_analysed {
			if i == j {
				continue searcher
			}
		}
		break searcher
	}

	// generating all possible entities
	for l, analysis := range toAnalyse {
		cur.tokens_analysed = append(cur.tokens_analysed, i+l)
		cur.tokens_toAnalyse[0] = cur.tokens_toAnalyse[0][i+l+1:]
		if len(cur.tokens_toAnalyse[0]) == 0 {
			cur.tokens_toAnalyse = cur.tokens_toAnalyse[1:]
		}
		if entry, ok := ics.entityDictionary[analysis]; ok {
			entry.DEMutex.RLock()
			for _, possibleEntity := range entry.entities {
				idp := new(ics_DijkstraPossibility)
				idp.tokens_analysed = make([]int, len(cur.tokens_analysed),
					cap(cur.tokens_analysed))
				copy(idp.tokens_analysed, cur.tokens_analysed)
				idp.collapsedAnalysedSentence = utils.DeepCopy[positionTuple,
																Entity](cur.collapsedAnalysedSentence)
				pos := positionTuple{i:i, j:i+len(strings.Split(string(analysis), " "))}
				idp.collapsedAnalysedSentence[pos] = possibleEntity
				idp.tokens_toAnalyse = make([][]Word, len(cur.tokens_toAnalyse))
				copy(idp.tokens_toAnalyse, cur.tokens_toAnalyse)
				handler <- idp
			}
			entry.DEMutex.RUnlock()
		}
	}
	return nil
}
