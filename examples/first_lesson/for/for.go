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
		fmt.Printf("%d => %d\n", index, arr[index])
		index ++
	}
	fmt.Println(" for loop end \n ")
}

func ForI()  {
	arr := []int {9, 8, 7, 6}
	for i := 0; i < len(arr); i++ {
		fmt.Printf("%d => %d \n", i, arr[i])
	}
	fmt.Println("for i loop end \n ")
}

func ForR()  {
	arr := []int {9, 8, 7, 6}

	for index, value := range arr {
		fmt.Printf("%d => %d\n", index, value)
	}

	// 如果只是需要 value, 可以用 _ 代替 index
	for _, value := range arr {
		fmt.Printf("only value: %d \n", value)
	}

	// 如果只需要 index 也可以去掉 写成 for index := range arr
	for index := range arr {
		fmt.Printf("only index: %d \n", index)
	}

	fmt.Println("for r loop end \n ")
}









