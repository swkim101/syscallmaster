# syscalls.master parser

```go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/swkim101/syscallmaster/syscallmaster"
)

var (
	flagFilename = flag.String("f", "foo", "syscalls.master file to parse")
)

func main() {
	flag.Parse()

	dat, err := os.ReadFile(*flagFilename)
	if err != nil {
		panic(err)
	}

	res := syscallmaster.Parse(dat, *flagFilename)
	js, _ := json.MarshalIndent(res[:2], "", " ")
	fmt.Printf("%v\n", string(js))
}
```

```sh
$ go run ./main.go -f syscalls.master
[
 {
  "Number": 0,
  "Audit": "AUE_NULL",
  "Files": "SYSMUX",
  "Decl": "{\n\t\tint syscall(\n\t\t    int number,\n\t\t    ...\n\t\t);\n\t}",
  "Comments": "{\n\t\tint syscall(\n\t\t    int number,\n\t\t    ...\n\t\t);\n\t}"
 },
 {
  "Number": 1,
  "Audit": "AUE_EXIT",
  "Files": "STD|CAPENABLED",
  "Decl": "{\n\t\tvoid exit(\n\t\t    int rval\n\t\t);\n\t}",
  "Comments": "{\n\t\tvoid exit(\n\t\t    int rval\n\t\t);\n\t}"
 }
]
```