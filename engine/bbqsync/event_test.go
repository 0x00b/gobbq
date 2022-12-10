package bbqsync

// import (
// 	"testing"

// 	"google.golang.org/grpc/internal/grpctest"
// )

// type s struct {
// 	grpctest.Tester
// }

// func Test(t *testing.T) {
// 	grpctest.RunSubTests(t, s{})
// }

// func (s) TestEventHasFired(t *testing.T) {
// 	e := NewEvent()
// 	if e.HasFired() {
// 		t.Fatal("e.HasFired() = true; want false")
// 	}
// 	if !e.Fire() {
// 		t.Fatal("e.Fire() = false; want true")
// 	}
// 	if !e.HasFired() {
// 		t.Fatal("e.HasFired() = false; want true")
// 	}
// }

// func (s) TestEventDoneChannel(t *testing.T) {
// 	e := NewEvent()
// 	select {
// 	case <-e.Done():
// 		t.Fatal("e.HasFired() = true; want false")
// 	default:
// 	}
// 	if !e.Fire() {
// 		t.Fatal("e.Fire() = false; want true")
// 	}
// 	select {
// 	case <-e.Done():
// 	default:
// 		t.Fatal("e.HasFired() = false; want true")
// 	}
// }

// func (s) TestEventMultipleFires(t *testing.T) {
// 	e := NewEvent()
// 	if e.HasFired() {
// 		t.Fatal("e.HasFired() = true; want false")
// 	}
// 	if !e.Fire() {
// 		t.Fatal("e.Fire() = false; want true")
// 	}
// 	for i := 0; i < 3; i++ {
// 		if !e.HasFired() {
// 			t.Fatal("e.HasFired() = false; want true")
// 		}
// 		if e.Fire() {
// 			t.Fatal("e.Fire() = true; want false")
// 		}
// 	}
// }
