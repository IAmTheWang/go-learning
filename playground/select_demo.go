package main

import "fmt"

func main() {
	ch := make(chan int) // 没有任何人往里面发东西

	select {
	case v := <-ch:
		fmt.Println("收到了:", v)
	default:
		// ch 里没东西也没关系，不阻塞，直接走这里
		fmt.Println("ch 里现在没东西，我不等，直接干别的")
	}
}