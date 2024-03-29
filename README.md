# goenchant

This package provides bindings for the enchant spell checking library.

See https://github.com/danielx/enchant2 for enchant v2 bindings

## Prerequisites

You need to have `libenchant` version 1.6 and its development files
installed to use this package. Additionally you need to install the
dictionaries (hunspell, aspell etc.) you would like to use.

### Example
```
sudo apt-get install enchant libenchant-dev hunspell-en
```

## Install

```
go get github.com/danielx/enchant
```

## Running tests

Given you have installed the prerequisites locally you can run the tests with:
```
go test
```

Or you can run the tests in a docker container with:
```
make test
```

## Docs

https://godoc.org/github.com/danielx/enchant

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/danielx/enchant"
)

func main() {
	e := enchant.New()
	defer e.Free()

	err := e.DictLoad("en_US")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	found, err := e.DictCheck("hello")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if found {
		fmt.Println("\"hello\" found in dictionary")
	} else {
		fmt.Println("\"hello\" not found in dictionary")
	}
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details


## Acknowledgements

- https://github.com/hermanschaaf/enchant
- https://github.com/taruti/enchant
