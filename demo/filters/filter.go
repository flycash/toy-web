package filters

import (
	"fmt"
	web "geektime/toy-web/pkg"
)

func init() {
	web.RegisterFilter("my-custom", myFilterBuilder)
}

func myFilterBuilder(next web.Filter) web.Filter {
	return func(c *web.Context) {
		fmt.Println("假装这是我自定义的 filter")
		next(c)
	}
}
