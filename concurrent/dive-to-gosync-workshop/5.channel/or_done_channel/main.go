package main

import "fmt"

func orDone(done <-chan struct{}, c <-chan any) <-chan any {
	valStream := make(chan any)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func main() {
	ch := make(chan any)
	go func() {
		defer close(ch)

		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for v := range orDone(nil, ch) {
		fmt.Printf("%v ", v)
	}
}
