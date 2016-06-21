package service

import (
  "../ifc"
  "v.io/v23/context"
  "v.io/v23/rpc"
)

type impl struct {
}

func Make() ifc.HelloServerMethods {
  return &impl {}
}

func (f *impl) Get(_ *context.T, _ rpc.ServerCall) (
    greeting string, err error) {
  return "Hello World!", nil
}
