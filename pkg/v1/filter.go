package webv1

import (
	"fmt"
	"time"
)

type FilterBuilder func(next Filter) Filter

type Filter func(c *Context)

func MetricFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		// 执行前的时间
		startTime := time.Now().UnixNano()
		next(c)
		// 执行后的时间
		endTime := time.Now().UnixNano()
		fmt.Printf("run time: %d \n", endTime-startTime)
	}
}