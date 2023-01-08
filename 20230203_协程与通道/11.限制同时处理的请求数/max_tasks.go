package main

import "fmt"

const MAXREQS = 50

var sem = make(chan int, MAXREQS)

type Request struct {
	a, b   int
	replyc chan int
}

func process(r *Request) {
	// do sonething
	// ...
	r.replyc <- r.a + r.b
	fmt.Println(r.replyc)

}

func handle(r *Request) {
	// doesn't matter what we put in it
	sem <- 1
	process(r)
	<-sem
}

func server(service chan *Request) {
	for {
		request := <-service
		go handle(request)
	}
}

func main() {
	service := make(chan *Request)
	go server(service)
}
