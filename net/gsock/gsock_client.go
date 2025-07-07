package gsock

import (
	"context"
	"net"

	"github.com/sourcegraph/jsonrpc2"
)

// RPCClient provides JSON-RPC 2.0 client functionality over Unix domain sockets
type rpcClient struct {
	sockPath  string        // Path to the Unix domain socket
	adapter   ClientAdapter // Protocol adapter (defaults to JSON-RPC)
	idCounter int64         // Atomic counter for generating request IDs
}

// NewRPCClient creates a new RPC client instance configured for Unix domain sockets
//
// Parameters:
//   - socketPath: Filesystem path to the Unix domain socket (e.g. "/tmp/rpc.sock")
//
// Returns:
//   - *RPCClient: A new client instance ready for making RPC calls
func NewRpcSimpleClient(socketPath string) *rpcClient {
	return &rpcClient{
		sockPath:  socketPath,
		adapter:   &JsonRpcSimpleClient{},
		idCounter: 0,
	}
}

// Request makes a JSON-RPC 2.0 method call and decodes the response
//
// Parameters:
//   - ctx: Context for request cancellation/timeout
//   - method: RPC method name to call
//   - params: Input parameters (will be JSON-encoded)
//   - result: Pointer to struct for decoding the response
//   - opts: Additional JSON-RPC 2.0 call options
//
// Returns:
//   - error: Any error that occurred during the RPC call
//
// Example:
//
//	var result MyResult
//	err := client.Request(ctx, "MyMethod", MyParams{}, &result)
func (c *rpcClient) Request(ctx context.Context, method string, params, result any, opts ...jsonrpc2.CallOption) error {
	conn, err := net.Dial("unix", c.sockPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 添加自定义ID生成器选项
	c.idCounter++
	idOpt := jsonrpc2.PickID(jsonrpc2.ID{Num: uint64(c.idCounter)})
	opts = append(opts, idOpt)

	client := c.adapter.NewConn(ctx, conn)
	return client.Request(ctx, method, params, result, opts...)
}
