package main

import (
	"context"
	"fmt"
	stlog "log"
	"sds/grades"
	"sds/log"
	"sds/registry"
	"sds/service"
)

func main() {
	host, port := "localhost", "6000"

	serviceAddress := fmt.Sprintf("http://%v:%v", host, port) // localhost:6000

	r := registry.Registration{
		ServiceName:      registry.GradingService,
		ServiceURL:       serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}

	// 启动服务
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers,
	)

	if err != nil {
		stlog.Fatal(err)

	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("logging service found at :%s\n", &logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}
	<-ctx.Done()

	fmt.Println("Shutting down grading service.")
}
