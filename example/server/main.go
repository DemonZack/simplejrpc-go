package main

import (
	"github.com/DemonZack/simplejrpc-go/net/gsock"
	"github.com/DemonZack/simplejrpc-go/server"
)

type CustomHandler struct{}

func (c *CustomHandler) Hello(req *gsock.Request) (any, error) {
	return "Hello World", nil
}

func (c *CustomHandler) ProcessRequest(req *gsock.Request) {
	println("[*] ProcessRequest. ", req)
}

func (c *CustomHandler) ProcessResponse(resp any) (any, error) {
	println("[*] ProcessResponse. ", resp)
	return resp, nil
}

func main() {
	mockSockPath := "zack.sock"

	ds := server.NewDefaultServer(
		gsock.WithJsonRpcSimpleServiceHandler(gsock.NewJsonRpcSimpleServiceHandler()),
		gsock.WithJsonRpcSimpleServiceMiddlewares([]gsock.RPCMiddleware{
			&CustomMiddleware{},
		}),
	)

	hand := &CustomHandler{}
	ds.RegisterHandle("hello", hand.Hello, []gsock.RPCMiddleware{hand}...)
	err := ds.StartServer(mockSockPath)
	if err != nil {
		panic(err)
	}
}

type CustomMiddleware struct{}

func (c *CustomMiddleware) ProcessRequest(req *gsock.Request) {
	println("[*] ProcessRequest. ", req)
}

func (c *CustomMiddleware) ProcessResponse(resp any) (any, error) {
	println("[*] ProcessResponse. ", resp)
	return resp, nil
}
