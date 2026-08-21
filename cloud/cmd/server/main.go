package main

import (
	"fmt"

	"github.com/boqrs/OpenIndustrial/cloud/bootstrap"
	"github.com/boqrs/OpenIndustrial/cloud/internal/pkg/code"
	"github.com/boqrs/zeus"
	"github.com/boqrs/zeus/cmd"
	"github.com/boqrs/zeus/ginx"
)

func main() {
	d := zeus.NewZeus()
	gcmd := cmd.NewGinCommand()
	// code must bigger than 100000
	if err := ginx.SetDefaultErrorCode(code.InternalError.Code); err != nil {
		fmt.Printf("set default error code failed: %v", err)
		return
	}

	// add global gin middleware
	finishFuc, err := bootstrap.InitInfra(gcmd.ZeroGinRouter)
	if err != nil {
		fmt.Printf("init infra failed: %v", err)
		return
	}

	gcmd.Flags().String("host", "0.0.0.0", "http server host")
	if err = d.ZeusStart("forge", gcmd); err != nil {
		fmt.Printf("zeus start error %v\n", err)
		return
	}

	finishFuc()
}