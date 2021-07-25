package api

import (
	"net/http"
)

type Context struct {
	Err *AppError
}

type AppError struct{}

func (er *AppError) Error() string {
	return "This is an error"
}

func iR() *AppError {
	return &AppError{}
}

func a(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	if err != nil { // ok
		c.Err = err
		return
	}
}

func b(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	if err != nil { // not ok
		c.Err = err
	}
}

func c(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	if err != nil { // ok
		c.Err = err
		return
	}

	err = iR()
	if err != nil { // ok
		c.Err = err
		return
	}
}

func d(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	if err != nil { // not ok
		c.Err = err
	}

	err = iR()
	if err != nil { // not ok
		c.Err = err
	}
}

func e(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()
	v := iR()
	s := true
	z := true

	if err != nil && t == nil { // ok
		c.Err = err
		return
	}

	if t == nil && err != nil { // ok
		c.Err = err
		return
	}

	if s && t == nil && err != nil && v != nil && z { // ok
		c.Err = err
		return
	}

	if err != nil && s && t == nil && v != nil && z { // ok
		c.Err = err
		return
	}

	if s && t == nil && v != nil && z && err != nil { // ok
		c.Err = err
		return
	}
}

func f(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()
	v := iR()
	s := true
	z := true

	if err != nil && t == nil { // not ok
		c.Err = err
	}

	if t == nil && err != nil { // not ok
		c.Err = err
	}

	if s && t == nil && err != nil && v != nil && z { // not ok
		c.Err = err
	}

	if err != nil && s && t == nil && v != nil && z { // not ok
		c.Err = err
	}

	if s && t == nil && v != nil && z && err != nil { // not ok
		c.Err = err
	}
}

func g(c *Context, w http.ResponseWriter, r *http.Request) {
	t := true
	err := iR()

	if err != nil { // ok
		if t {
			c.Err = err
		}
		return
	}
}

func h(c *Context, w http.ResponseWriter, r *http.Request) {
	t := true
	err := iR()

	if err != nil { // not ok
		if t {
			c.Err = err
		}
	}
}

func i(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()

	if (err != nil) && (t == nil) { // ok
		c.Err = err
		return
	}

}

func j(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()

	if (err != nil) && (t == nil) { // not ok
		c.Err = err
	}
}

func k(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()

	if (err != nil) && t == nil { // ok
		c.Err = err
		return
	}
}

func l(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()

	if (err != nil) && t == nil { // not ok
		c.Err = err
	}
}

func m(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()

	if t == nil && (err != nil) { // ok
		c.Err = err
		return
	}
}

func n(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()

	if t == nil && (err != nil) { // not ok
		c.Err = err
	}
}

func o(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()
	s := true

	if t == nil && (err != nil) && s { // ok
		c.Err = err
		return
	}
}

func p(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()
	s := true

	if t == nil && (err != nil) && s { // not ok
		c.Err = err
	}
}

// skip appErrReturn check
func q(c *Context, w http.ResponseWriter, r *http.Request) {
	err := iR()
	t := iR()
	s := true

	if t == nil && (err != nil) && s { // not ok but skipped
		c.Err = err
	}
}
