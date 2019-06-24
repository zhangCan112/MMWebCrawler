package pipeline

import (
	"fmt"
)

var CSVWriter = HandlerFunc(func(first Item, rest ...Item) error {
	total := make([]Item, len(rest)+1)[0:0]
	total = append(total, first)
	total = append(total, rest...)

	for _, it := range total {
		for key, val := range it.KeyValues() {
			fmt.Printf("%s:%s", key, val)
		}
		fmt.Print("\n")
	}
	return nil
})
