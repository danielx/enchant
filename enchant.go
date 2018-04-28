// Package enchant provides bindings to the Enchant spell checking library.
package enchant

/*
#cgo LDFLAGS: -lenchant
#include <stdlib.h>
#include <enchant/enchant.h>

static char* getString(char **c, int i) {
	return c[i];
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Enchant is the type that contains the package internals.
type Enchant struct {
	broker *C.EnchantBroker
	dict   *C.EnchantDict
}

// New allocates a new Enchant instance.
func New() (e *Enchant, err error) {
	return &Enchant{C.enchant_broker_init(), nil}, nil
}

// Free releases the allocated memory related to the Enchant instance.
func (e *Enchant) Free() {
	if e.dict != nil {
		C.enchant_broker_free_dict(e.broker, e.dict)
		e.dict = nil
	}

	if e.broker != nil {
		C.enchant_broker_free(e.broker)
		e.broker = nil
	}
}

// DictExists checks if a dictionary exists or not.
func (e *Enchant) DictExists(tag string) (exists bool, err error) {
	if e.broker == nil {
		return false, fmt.Errorf("no broker initialized")
	}

	cTag := C.CString(tag)
	defer C.free(unsafe.Pointer(cTag))
	return C.enchant_broker_dict_exists(e.broker, cTag) == 1, nil
}

// DictLoad loads a dictionary to spell check against.
func (e *Enchant) DictLoad(tag string) error {
	if e.broker == nil {
		return fmt.Errorf("no broker initialized")
	}

	if e.dict != nil {
		return fmt.Errorf("dictionary already loaded")
	}

	cTag := C.CString(tag)
	defer C.free(unsafe.Pointer(cTag))

	e.dict = C.enchant_broker_request_dict(e.broker, cTag)
	if e.dict == nil {
		return fmt.Errorf("no suitable dictionary could be found by name %s", tag)
	}

	return nil
}

// DictCheck checks if a word is found in the loaded dictionary.
func (e *Enchant) DictCheck(word string) (found bool, err error) {
	if e.broker == nil {
		return false, fmt.Errorf("no broker initialized")
	}

	if e.dict == nil {
		return false, fmt.Errorf("no dictionary loaded")
	}

	if len(word) == 0 {
		return true, nil
	}

	cWord := C.CString(word)
	defer C.free(unsafe.Pointer(cWord))

	size := uintptr(len(word))
	s := (*C.ssize_t)(unsafe.Pointer(&size))

	status := C.enchant_dict_check(e.dict, cWord, *s)

	if status < 0 {
		return false, fmt.Errorf("could not check word: %s", word)
	}

	return status == 0, nil
}

// DictSuggest suggests spelling for a word.
func (e *Enchant) DictSuggest(word string) (suggestions []string, err error) {
	if e.broker == nil {
		return nil, fmt.Errorf("no broker initialized")
	}

	if e.dict == nil {
		return nil, fmt.Errorf("no dictionary loaded")
	}

	if len(word) == 0 {
		return nil, nil
	}

	cWord := C.CString(word)
	defer C.free(unsafe.Pointer(cWord))

	size := uintptr(len(word))
	s := (*C.ssize_t)(unsafe.Pointer(&size))

	// number of suggestions returned.
	var suggs uintptr
	ns := (*C.size_t)(unsafe.Pointer(&suggs))

	response := C.enchant_dict_suggest(e.dict, cWord, *s, ns)
	defer C.enchant_dict_free_string_list(e.dict, response)
	if response == nil {
		return nil, nil
	}

	suggestions = make([]string, 0, suggs)
	for i := uintptr(0); i < suggs; i++ {
		suggestions = append(suggestions, C.GoString(C.getString(response, C.int(i))))
	}

	return suggestions, nil
}
