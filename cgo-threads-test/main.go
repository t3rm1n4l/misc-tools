package main

//#include <unistd.h>
import "C"
import "fmt"
import "sync"

func sleep(wg *sync.WaitGroup) {
        defer wg.Done()
        C.sleep(C.uint(120))
}

func main() {
        var wg sync.WaitGroup

        for i := 0; i < 1000; i++ {
                wg.Add(1)
                go sleep(&wg)
        }

        wg.Wait()
        fmt.Println("done")
}
