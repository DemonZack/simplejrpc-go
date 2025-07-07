package server

import (
	"github.com/DemonZack/simplejrpc-go/net/gsock"
	"github.com/DemonZack/simplejrpc-go/server"
)

type AppServer struct{}

func NewAppServer() *AppServer {
	return &AppServer{}
}

func (s *AppServer) Run() {
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
