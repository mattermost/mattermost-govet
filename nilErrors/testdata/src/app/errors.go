package app

import (
	"errors"
	"fmt"
	"math/rand"
)

var check bool
var mlog mlogS

func func1() error {
	err1 := errorGen()
	err2 := errorGen()
	if err1 != nil {
		mlog.Error("error", err2)
	}

	return nil
}

func func2() error {
	err1 := errorGen()
	if err1 != nil {
		return fmt.Errorf("some error: %s", err1.Error())
	}
	err2 := errorGen()
	if err2 != nil && err1 != nil {
		return fmt.Errorf("some error: %s", err1.Error())
	}

	return nil
}

func func3() error {
	_, err1 := errorGen2()
	_, err2 := errorGen2()
	if err1 != nil {
		return fmt.Errorf("some error: %s", err2.Error())
	}

	return nil
}

func func4() error {
	_, err1 := errorGen2()
	_, err2 := errorGen2()
	var errNotExist *ErrorNotExist
	if err1 != nil {
		switch {
		case errors.As(err1, &errNotExist):
			return fmt.Errorf("some error: %s", errNotExist.Error())
		default:
			return fmt.Errorf("some error: %s", err2.Error())
		}
	}

	return nil
}

func errorGen() error {
	if rand.Intn(1) > 0 {
		return errors.New("random error")
	}

	return nil
}

func errorGen2() (int, error) {
	if rand.Intn(1) > 0 {
		return 0, errors.New("random error")
	}

	return 1, nil
}

type ErrorNotExist struct {
	msg string
}

func (e ErrorNotExist) Error() string {
	return e.msg
}

type mlogS struct {
}

func (mlogS) Error(msg string, fields ...interface{}) {
	fmt.Printf("%s, %v\n", msg, fields)
}
