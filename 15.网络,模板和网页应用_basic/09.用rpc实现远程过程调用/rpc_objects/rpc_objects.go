package rpc_objects

import "net"

type Args struct {
	N, M int
}

// 旧版本语言要求此方法返回类型error,示例net.Error无法通过编译
func (t *Args) Multiply(args *Args, reply *int) net.Error {
	*reply = args.N * args.M
	return nil
}
