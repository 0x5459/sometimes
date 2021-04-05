package main

import "fmt"

type Stack struct {
	data []int
	top  int
}

func (s *Stack) Push(v int) {
	s.data[s.top] = v
	s.top++
}

func (s *Stack) Pop() int {
	s.top--
	return s.data[s.top]
}

let a = 10;


func main() {
	s := Stack{data: make([]int, 64), top: 0}

	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)

	fmt.Println(s.data)
	s.top -= 3
	s.Push(1000)
	fmt.Println(s.data)
	s.Push(200)
	fmt.Println(s.data)
	fmt.Println(s.Pop())
	fmt.Println(s.data)
}
