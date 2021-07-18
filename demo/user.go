package demo

import (
	"fmt"
	web "geektime/toy-web/pkg"
	"time"
)

func SignUp(c *web.Context) {
	req := &signUpReq{}
	err := c.ReadJson(req)
	if err != nil {
		_ = c.BadRequestJson(&commonResponse{
			BizCode: 4, // 假如说我们这个代表输入参数错误
			// 注意这里是demo，实际中你应该避免暴露 error
			Msg: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}
	_ = c.OkJson(&commonResponse{
		// 假设这个是新用户的 ID
		Data: 123,
	})
}

func SlowService(c *web.Context) {
	time.Sleep(time.Second * 10)
	_ = c.OkJson(&commonResponse{
		Msg: "Hi, this is msg from slow service",
	})
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