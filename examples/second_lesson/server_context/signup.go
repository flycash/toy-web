package server_context

import (
	"encoding/json"
	"fmt"
	"geektime/toy-web/pkg/v2"
	"io"
	"net/http"
)

// 在没有 context 抽象的情况下，是长这样的
func SignUpWithoutContext(w http.ResponseWriter, r *http.Request) {
	req := &signUpReq{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed: %v", err)
		// 要返回掉，不然就会继续执行后面的代码
		return
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		fmt.Fprintf(w, "deserialized failed: %v", err)
		// 要返回掉，不然就会继续执行后面的代码
		return
	}

	// 返回一个虚拟的 user id 表示注册成功了
	fmt.Fprintf(w, "%d", err)
}

func SignUpWithoutWrite(w http.ResponseWriter, r *http.Request) {
	c := webv2.NewContext(w, r)
	req := &signUpReq{}
	err := c.ReadJson(req)
	if err != nil {
		resp := &commonResponse{
			BizCode: 4, // 假如说我们这个代表输入参数错误
			Msg: fmt.Sprintf("invalid request: %v", err),
		}
		respBytes, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respBytes))
		return
	}
	// 这里又得来一遍 resp 转json
	fmt.Fprintf(w, "invalid request: %v", err)
}

type signUpReq struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int `json:"biz_code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
