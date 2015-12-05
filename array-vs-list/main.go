package main

import "time"
import "fmt"

const arrsz = 7

var rand int

type node struct {
	x    int
	next *node
}

type arraylist struct {
	arr  [arrsz]int
	next *arraylist
}

func newArrayList(n int) *arraylist {
	var head *arraylist
	var last *node
	for i := 0; i < n; i++ {
		if i%arrsz == 0 {
			head = &arraylist{next: head}
			last = &node{next: &node{next: &node{}}}
			if last.next == nil {
				rand++
			}
		}

		head.arr[i%arrsz] = i
	}

	return head
}

func newList(n int) *node {
	var head *node
	var last *node
	for i := 0; i < n; i++ {
		head = &node{x: i, next: head}
		last = &node{next: &node{next: &node{next: &node{next: &node{}}}}}
	}

	if last.next == nil {
		rand++
	}

	return head
}

func sumArrayList(r *arraylist) int {
	var s int

	mod := 0
	for r != nil {
		s += r.arr[mod]
		if mod == arrsz-1 {
			r = r.next
			mod = 0
		} else {
			mod++
		}
	}

	return s
}

func sumList(r *node) int {
	var s int

	for r != nil {
		s += r.x
		r = r.next
	}

	return s
}

func main() {
	n := 100000000
	t0 := time.Now()
	l := newList(n)
	fmt.Println("build list", time.Since(t0))
	t0 = time.Now()
	a := newArrayList(n)
	fmt.Println("build alist", time.Since(t0))

	t0 = time.Now()
	s := sumList(l)
	fmt.Println("sum list", time.Since(t0), s)

	t0 = time.Now()
	s = sumArrayList(a)
	fmt.Println("sum alist", time.Since(t0), s)
	fmt.Println(rand)
}
