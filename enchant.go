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
	"sync"
	"unsafe"
)

// Enchant is the type that contains the package internals.
type Enchant struct {
	sync.Mutex
	broker *C.EnchantBroker
	dict   *C.EnchantDict
}

// New allocates a new Enchant instance.
func New() (e *Enchant) {
	return &Enchant{broker: C.enchant_broker_init()}
}

// Free frees the allocated memory related to the Enchant instance.
func (e *Enchant) Free() {
	e.DictFree() // make sure dictionary is freed to prevent memory leaks.

	if e.broker != nil {
		C.enchant_broker_free(e.broker)
		e.broker = nil
	}
}

// DictFree frees the current dictionary allowing another to be loaded.
func (e *Enchant) DictFree() {
	if e.broker != nil && e.dict != nil {
		C.enchant_broker_free_dict(e.broker, e.dict)
		e.dict = nil
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
		return fmt.Errorf("a dictionary is already loaded")
	}

	cTag := C.CString(tag)
	defer C.free(unsafe.Pointer(cTag))

	e.dict = C.enchant_broker_request_dict(e.broker, cTag)
	if e.dict == nil {
		return fmt.Errorf("failed to load dictionary by tag: %s", tag)
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

	e.Lock()
	status := C.enchant_dict_check(e.dict, cWord, *s)
	e.Unlock()

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

	e.Lock()
	defer e.Unlock()

	response := C.enchant_dict_suggest(e.dict, cWord, *s, ns)
	if response == nil {
		return nil, nil
	}
	defer C.enchant_dict_free_string_list(e.dict, response)

	suggestions = make([]string, 0, suggs)
	for i := uintptr(0); i < suggs; i++ {
		suggestions = append(suggestions, C.GoString(C.getString(response, C.int(i))))
	}

	return suggestions, nil
}
