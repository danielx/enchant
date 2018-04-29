package enchant

import (
	"testing"
)

func TestDictExists(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer e.Free()

	expect := map[string]bool{
		"en":        true,
		"gibberish": false,
		"en_US":     true,
	}

	for dict, expectedValue := range expect {
		found, err := e.DictExists(dict)

		if err != nil {
			t.Errorf(err.Error())
			return
		}

		if found != expectedValue {
			t.Errorf("Expected DictExist to return %v for \"%s\"", expectedValue, dict)
		}
	}
}

func TestDictLoad(t *testing.T) {
	expect := map[string]bool{
		"en":        false,
		"gibberish": true,
		"en_US":     false,
	}

	for dict, shouldError := range expect {
		e, err := New()
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		defer e.Free()

		err = e.DictLoad(dict)

		if err != nil && !shouldError {
			t.Errorf("unexpected error returned for dict: %s, %v", dict, err)
		} else if err == nil && shouldError {
			t.Errorf("expected an error to be returned for dict: %s", dict)
		}
	}
}

func TestDictLoadNoBroker(t *testing.T) {
	e := &Enchant{}

	err := e.DictLoad("en")

	if err == nil {
		t.Errorf("expected an error because no broker has been initialized")
	}
}

func TestDictLoadMultiple(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer e.Free()

	err = e.DictLoad("en")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	err = e.DictLoad("en_US")
	if err == nil {
		t.Errorf("expected an error to be returned when loading multiple dictionaries")
	}
}

func TestDictCheck(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer e.Free()

	err = e.DictLoad("en")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	expect := map[string]bool{
		"":         true,
		"helllo":   false,
		"world":    true,
		"bazinga!": false,
	}

	for word, expectedValue := range expect {
		found, err := e.DictCheck(word)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		if found != expectedValue {
			t.Errorf("expected DictCheck to return %v for word: %s", expectedValue, word)
		}
	}
}

func TestDictCheckNoBroker(t *testing.T) {
	e := &Enchant{}

	_, err := e.DictCheck("hi")
	if err == nil {
		t.Errorf("expected an error because no broker has been initialized")
	}
}

func TestDictCheckNoDict(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer e.Free()

	_, err = e.DictCheck("hi")
	if err == nil {
		t.Errorf("expected an error because no dictionary has been loaded")
	}
}

func TestDictSuggestion(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer e.Free()

	err = e.DictLoad("en")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	expect := map[string][]string{
		"":           nil,
		"unexpected": []string{"unexpected", "unaccepted", "expected"},
	}

	for word, expectedValue := range expect {
		suggestions, err := e.DictSuggest(word)

		if err != nil {
			t.Errorf("got unexpected error: %v", err)
			return
		}

		if suggestions == nil && expectedValue != nil {
			t.Errorf("expected %v for word: %s, got nil", word, expectedValue)
		} else if suggestions != nil && expectedValue == nil {
			t.Errorf("expected nil for word: %s, got: %v", word, suggestions)
		} else if len(suggestions) != len(expectedValue) {
			t.Errorf("expected %v for %s, got %v", expectedValue, word, suggestions)
		} else {
			for i := range suggestions {
				if suggestions[i] != expectedValue[i] {
					t.Errorf("expected %v for %s, got %v", expectedValue, word, suggestions)
					return
				}
			}
		}
	}
}

func TestDictSuggestNoBroker(t *testing.T) {
	e := &Enchant{}

	_, err := e.DictSuggest("hi")
	if err == nil {
		t.Errorf("expected an error because no broker has been initialized")
	}
}

func TestDictSuggestNoDict(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer e.Free()

	_, err = e.DictSuggest("hi")
	if err == nil {
		t.Errorf("expected an error because no dictionary has been loaded")
	}
}
