package utils

import "testing"

func TestMergeMaps(t *testing.T) {
	var m1, m2 map[int]int
	m1 = make(map[int]int)
	m2 = make(map[int]int)

	m1[1] = 1
	m2[2] = 2

	m := Merge(m1, m2)

	if i, ok := m[1]; !ok || i != 1 {
		t.Fail()
	}
	if i, ok := m[2]; !ok || i != 2 {
		t.Fail()
	}
	if _, ok := m[3]; ok {
		t.Fail()
	}
	if _, ok := m1[2]; ok {
		t.Fail()
	}
	if _, ok := m2[1]; ok {
		t.Fail()
	}
}

func TestMinKey(t *testing.T) {
	var m = make(map[int]int)

	m[1] = -1
	m[2] = -2

	key := MinKey(m, func(i, j int) bool { return i < j })
	if key != 1 {
		t.Error("cannot handle map. Got", key, "instead of", 1)
	}

	key = MinKey[int, int](nil, func(i, j int) bool { return i < j })
	if key != 0 {
		t.Fail()
		t.Error("cannot handle <nil> map")
	}
}

func TestReverseMap(t *testing.T) {
	var m = make(map[int]string)
	m[1] = "1"
	m[2] = "2"
	m[3] = "3"

	w, err := ReverseMap(m)

	if err != nil {
		t.Error(err)
	}

	if i, ok := w["1"]; !ok || i != 1 {
		t.Fail()
	}
	if i, ok := w["2"]; !ok || i != 2 {
		t.Fail()
	}

	m[4] = "1"
	_, err = ReverseMap(m)
	if err == nil {
		t.Error("Failed to detect non-reversible map")
	}
}
