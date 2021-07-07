package main

import "fmt"

func main() {
	ChooseFruit("蓝莓")
	ChooseFruit("苹果")
	ChooseFruit("西瓜")
}

func ChooseFruit(fruit string) {
	switch fruit {
	case "苹果":
		fmt.Println("这是一个苹果")
	case "草莓", "蓝莓":
		fmt.Println("这是霉霉")
	default:
		fmt.Printf("不知道是啥：%s \n", fruit)
	}
}
