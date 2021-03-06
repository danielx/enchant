package enchant

import (
	"testing"
)

func TestDictExists(t *testing.T) {
	e := New()
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

func TestDictExistsNoBroker(t *testing.T) {
	e := &Enchant{}

	_, err := e.DictExists("en")
	if err == nil {
		t.Errorf("expected an error because no broker has been initialized")
	}
}

func TestDictLoad(t *testing.T) {
	expect := map[string]bool{
		"en":        false,
		"gibberish": true,
		"en_US":     false,
	}

	for dict, shouldError := range expect {
		e := New()
		defer e.Free()

		err := e.DictLoad(dict)

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

func TestDictLoadMultipleWithoutFree(t *testing.T) {
	e := New()
	defer e.Free()

	err := e.DictLoad("en")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	err = e.DictLoad("en_US")
	if err == nil {
		t.Errorf("expected an error to be returned when loading multiple dictionaries")
	}
}

func TestDictLoadMultipleWithFree(t *testing.T) {
	e := New()
	defer e.Free()

	err := e.DictLoad("en")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	e.DictFree()

	err = e.DictLoad("en_US")
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestDictCheck(t *testing.T) {
	e := New()
	defer e.Free()

	err := e.DictLoad("en")
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
	e := New()
	defer e.Free()

	_, err := e.DictCheck("hi")
	if err == nil {
		t.Errorf("expected an error because no dictionary has been loaded")
	}
}

func TestDictSuggestion(t *testing.T) {
	e := New()
	defer e.Free()

	err := e.DictLoad("en")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	expect := map[string]bool{
		"":           false,
		"unexpected": true,
	}

	for word, expectedValue := range expect {
		suggestions, err := e.DictSuggest(word)

		if err != nil {
			t.Errorf("got unexpected error: %v", err)
			return
		}

		if suggestions == nil && expectedValue {
			t.Errorf("expected suggestions for word: %s, got none", word)
		} else if suggestions != nil && !expectedValue {
			t.Errorf("expected no suggestions for word: %s, got: %v", word, suggestions)
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
	e := New()
	defer e.Free()

	_, err := e.DictSuggest("hi")
	if err == nil {
		t.Errorf("expected an error because no dictionary has been loaded")
	}
}

func BenchmarkFreeBrokerAndDictAfterUse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := New()

		err := e.DictLoad("en")
		if err != nil {
			b.Errorf(err.Error())
			return
		}

		found, err := e.DictCheck("hello")
		if err != nil {
			b.Errorf(err.Error())
			return
		}

		if found != true {
			b.Errorf("expected hello to be found in the dictionary")
			return
		}

		e.Free()
	}
}

func BenchmarkFreeDictAfterUse(b *testing.B) {
	e := New()
	defer e.Free()

	for i := 0; i < b.N; i++ {
		err := e.DictLoad("en")
		if err != nil {
			b.Errorf(err.Error())
			return
		}

		found, err := e.DictCheck("hello")
		if err != nil {
			b.Errorf(err.Error())
			return
		}

		if found != true {
			b.Errorf("expected hello to be found in the dictionary")
			return
		}

		e.DictFree()
	}
}

func BenchmarkReuseDictAndBroker(b *testing.B) {
	e := New()
	defer e.Free()

	err := e.DictLoad("en")
	if err != nil {
		b.Errorf(err.Error())
		return
	}

	for i := 0; i < b.N; i++ {
		found, err := e.DictCheck("hello")
		if err != nil {
			b.Errorf(err.Error())
			return
		}

		if found != true {
			b.Errorf("expected hello to be found in the dictionary")
			return
		}
	}
}
