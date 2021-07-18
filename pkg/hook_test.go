package web

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildCloseServerHook(t *testing.T) {
	svr := NewSdkHttpServer("test-sever")
	h := BuildCloseServerHook(svr, svr, svr, svr, svr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	err := h(ctx)
	assert.Nil(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond * 10)
	defer cancel()
	err = h(ctx)
	assert.Equal(t, ErrorHookTimeout, err)
}
