package utilities

import "testing"

func TestCreateGenericMap(t *testing.T) {
	toConvert := map[int][]string{
		1: {"a", "b", "c"},
		2: {"d", "e", "f"},
	}

	converted := CreateGenericMap(toConvert)
	if converted[1].([]string)[0] != "a" {
		t.Errorf("expected map key 1 to have a first value of a, got %s instead", converted[1].([]string)[0])
	}

	if converted[2].([]string)[2] != "f" {
		t.Errorf("expected map key 2 to have a first value of f, got %s instead", converted[2].([]string)[2])
	}
}
