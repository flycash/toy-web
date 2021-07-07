package main

import "fmt"

func main() {
	ForLoop()
	ForI()
	ForR()
}

func ForLoop()  {
	arr := []int {9, 8, 7, 6}
	index := 0
	for {
		if index == 3{
			// break 跳出循环
			break
		}
		fmt.Printf("%d => %d ", index, arr[index])
		index ++
	}
	fmt.Println("\n for loop end ")
}

func ForI()  {
	arr := []int {9, 8, 7, 6}
	for i := 0; i < len(arr); i++ {
		fmt.Printf("%d => %d", i, arr[i])
	}
	fmt.Println("\n for i loop end ")
}

func ForR()  {
	arr := []int {9, 8, 7, 6}
	// 如果只是需要 value, 可以用 _ 代替 index
	// 如果只需要 index 也可以去掉 写成 for index := range arr
	for index, value := range arr {
		fmt.Printf("%d => %d", index, value)
	}
	fmt.Println("for r loop end ")
}









