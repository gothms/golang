package main

import (
	"errors"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group

	// 启动第一个子任务,它执行成功
	g.Go(func() error {
		time.Sleep(5 * time.Second)
		log.Println("exec #1")
		return nil
	})

	// 启动第二个子任务，它执行失败
	g.Go(func() error {
		time.Sleep(10 * time.Second)
		log.Println("exec #2")
		return errors.New("failed to exec #2")
	})

	// 启动第三个子任务，它执行成功
	g.Go(func() error {
		time.Sleep(15 * time.Second)
		log.Println("exec #3")
		return nil
	})

	if err := g.Wait(); err == nil {
		log.Println("Successfully exec all")
	} else {
		log.Println("failed:", err)
	}
}
