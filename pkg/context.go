package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
	PathParams map[string]string
}

func (c *Context) ReadJson(data interface{}) error {
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}
func (c *Context) OkJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.WriteJson(http.StatusOK, data)
}

func (c *Context) SystemErrJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.WriteJson(http.StatusInternalServerError, data)
}

func (c *Context) BadRequestJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.WriteJson(http.StatusBadRequest, data)
}

func (c *Context) WriteJson(status int, data interface{}) error {
	c.W.WriteHeader(status)
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.W.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
		// 一般路径参数都是一个，所以容量1就可以了
		PathParams: make(map[string]string, 1),
	}
}

func newContext() *Context {
	fmt.Println("create new context")
	return &Context{
	}
}

func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.W = w
	c.R = r
	c.PathParams = make(map[string]string, 1)
}