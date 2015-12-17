package main

import "fmt"
import "unsafe"

import "runtime"
import "reflect"

func main() {
	data := make([]unsafe.Pointer, 100)
	var next unsafe.Pointer

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	next = unsafe.Pointer(hdr.Data)

	for i := 0; i < 100; i++ {
		x := new(int)
		*x = i
		data[i] = unsafe.Pointer(x)
	}

	//	m := data[:100]
	hdr2 := *hdr
	hdr = nil
	data = nil
	hdr2.Data = 0

	runtime.GC()

	//xx := (*reflect.SliceHeader)(unsafe.Pointer(&m))
	//fmt.Println(xx.Data, hdr2.Data)

	hdr2.Data = uintptr(next)
	data = *(*[]unsafe.Pointer)(unsafe.Pointer(&hdr2))
	for i := 0; i < 100; i++ {
		//	fmt.Println(*(*int)(m[i]))
		fmt.Println(*(*int)(data[i]))
	}

}
