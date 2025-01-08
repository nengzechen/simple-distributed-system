package main

import (
	"context"
	"fmt"
	stlog "log"
	"sds/log"
	"sds/portal"
	"sds/registry"
	"sds/service"
)

func main() {
	err := portal.ImportTemplates()
	if err != nil {
		stlog.Fatal(err)
	}
	host, port := "localhost", "8080"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName: registry.PortalService,
		ServiceURL:  serviceAddress,
		RequiredServices: []registry.ServiceName{
			registry.LogService,
			registry.GradingService,
		},
		ServiceUpdateURL: serviceAddress + "/services",
		HeartbeatURL:     serviceAddress + "/heartbeat",
	}

	ctx, err := service.Start(context.Background(),
		host,
		port,
		r,
		portal.RegisterHandlers)
	if err != nil {
		stlog.Fatal(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		log.SetClientLogger(logProvider, r.ServiceName)
	}
	<-ctx.Done()
	fmt.Println("Shutting down portal")
}
