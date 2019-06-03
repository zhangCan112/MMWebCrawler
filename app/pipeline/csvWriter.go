package pipeline

import (
	"fmt"
)
var csvWriter = HandlerFunc(func(first Item, rest ...Item) error {
	total := make([]Item, len(rest)+1)[0:0]
	total = append(total, first)
	total = append(total, rest...)

	for _, it := range total {
		for key, val := range it.KeyValues() {
			fmt.Printf("%s:%s",key,val)
		}
	}
	return nil
})
