package configloader

import (
	"fmt"
	"os"
)

type Example struct {
	Str   string
	Int   int64
	Slice []string
	Map   map[string]string
}

func ExampleJsonLoader() {
	os.Args = append(os.Args, "-config=testresources/config.json")

	example := Example{}
	Load(&example)
	fmt.Println(example)

	//Output:
	//{string 123 [a b 1 g] map[a:1 b:2 c:3]}
}
