package sync2

import (
	"errors"
	"testing"
)

type one struct {
	Counter  int
	InitFail int
}

func (o *one) Inc() error {
	defer func() {
		o.Counter++
	}()
	if o.Counter < o.InitFail {
		return errors.New("error")
	}
	return nil
}

func run(t *testing.T, once *OnceSucc, o *one, c chan error) {
	err := once.Do(func() error { return o.Inc() })
	c <- err
}

func TestDoOnceSucc(t *testing.T) {
	testCases := map[string]struct {
		InitFailTimes int
		Concurrents   int
		Expected      int
	}{
		"all succ": {
			InitFailTimes: 0,
			Concurrents:   100,
			Expected:      1,
		},
		"fail call": {
			InitFailTimes: 7,
			Concurrents:   100,
			Expected:      8,
		},
	}

	for label, testCase := range testCases {
		o := &one{
			InitFail: testCase.InitFailTimes,
		}
		once := &OnceSucc{}
		gotErrNum := 0
		c := make(chan error, testCase.Concurrents)
		for i := 0; i < testCase.Concurrents; i++ {
			go run(t, once, o, c)
		}
		for i := 0; i < testCase.Concurrents; i++ {
			err := <-c
			if err != nil {
				gotErrNum++
			}
		}
		if gotErrNum != testCase.InitFailTimes {
			t.Errorf("label: %s expected err num: %d; got: %d", label, testCase.InitFailTimes, gotErrNum)
		}

		if o.Counter != testCase.Expected {
			t.Errorf("label: %s expected counter: %d; got: %d", label, testCase.Expected, o.Counter)
		}
	}
}
