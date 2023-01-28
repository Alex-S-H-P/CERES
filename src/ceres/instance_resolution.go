package ceres

import (
	"CERES/src/utils"
	"CERES/src/utils/errors"
)

/*
Given an entity type, and the attributes that can be found in the sentence
or in the context, returns the Instance.
*/
func (c *CERES) solveInstance(et *EntityType,
	attrs map[*EntityType]Word) (*EntityInstance,
	error) {
	var eis []*EntityInstance = et.allDescendantsInstance()
	if len(eis) == 0 {
		return nil, errors.NoMatchingEntityInstanceFound
	}

	var bestP float64 = 0.
	var bestEI *EntityInstance = nil

	for _, ei := range eis {

		P := p_joined_ei_resolution(ei, attrs, &c.ics, c.ctx)

		if P > bestP {
			bestP = P
			bestEI = ei
		}
	}

	if bestEI == nil {
		return nil, errors.NoMatchingEntityInstanceFound
	}
	return bestEI, nil
}

func (et *EntityType) allDescendantsInstance() []*EntityInstance {
	if et == nil {
		return nil
	}

	arr := make([]*EntityInstance, 0, len(et.links))
	for _, link := range et.links {
		if link.typeOfLink() != "HYPERNYM" {
			continue
		}

		child := link.GetB()
		switch c := child.(type) {
		case *EntityType:
			desc := c.allDescendantsInstance()
			arr = append(arr, desc...)
		case *EntityInstance:
			arr = append(arr, c)
		}

	}
	return arr
}

func count_occurences_attribute_value(eis []*EntityInstance) map[*EntityType]map[Word]int {
	var m = make(map[*EntityType]map[Word]int)

	for _, ei := range eis {
		for attr, value := range ei.values.values {
			if subM, ok := m[attr]; ok {
				if count, ok := subM[value]; ok {
					subM[value] = count + 1
				} else {
					m[attr][value] = 1
				}

			} else {
				m[attr] = make(map[Word]int)
				m[attr][value] = 1
			}
		}
	}

	return m
}

func p_attrs_ei_resolution(m map[*EntityType]map[Word]int,
	attrs map[*EntityType]Word) float64 {
	var result float64 = .0

	for attr, w := range attrs {
		if submap, ok := m[attr]; ok {
			count_neighborhood := submap[w]
			N := utils.SumMap(submap)
			result *= float64(count_neighborhood) / float64(N)
		} else {
			return 0.
		}

	}

	return result
}

func p_joined_ei_resolution(ei *EntityInstance,
	attrs map[*EntityType]Word,
	ics *ICS,
	ctx *CTX) float64 {

	var P float64 = 1.

	if pece, ok := ctx.Contains(ei); ok {
		return pece.s
	}

	for attr, value := range ei.values.values {
		if val, ok := attrs[attr]; ok {
			if val == value {
				continue
			} else {
				P = 0.
			}
		} else {
			P *= .9
		}
	}

	return P
}
