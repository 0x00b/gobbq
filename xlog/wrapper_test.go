package xlog_test

// import (
// 	"context"
// 	"os"
// 	"sync"
// 	"testing"

// 	"go.qcloud.com/ago"
// 	"go.qcloud.com/ago/log"
// )

// func TestWrapper(t *testing.T) {

// 	log.Init("trace", true, true, os.Stdout, "", log.DefaultLogTag{})

// 	c := context.Background()

// 	log.Wrapper(func(ctx entity.Context) {

// 		// organize log
// 		log.Infoln(ctx, "test in ")
// 		log.Infoln(ctx, "in info?")

// 		// not organize log
// 		log.Infoln(c, "test out")
// 		log.Infoln(c, "info?")

// 	})(c)

// 	log.Infoln(c, "test wrapper")

// }

// func TestGoWrapper(t *testing.T) {

// 	log.Init("trace", true, true, os.Stdout, "", log.DefaultLogTag{})

// 	c := context.Background()
// 	wg := sync.WaitGroup{}
// 	wg.Add(1)
// 	go log.Wrapper(func(ctx entity.Context) {

// 		// organize log
// 		log.Infoln(ctx, "test in ")
// 		log.Infoln(ctx, "in info?")

// 		// not organize log
// 		log.Infoln(c, "test out")
// 		log.Infoln(c, "info?")

// 		wg.Done()

// 	})(ago.Copy(c))

// 	log.Infoln(c, "test wrapper")

// 	wg.Wait()
// }
