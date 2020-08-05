
package main

import (
	"context"
	"k8s.io/klog"
	"os"
	"os/signal"
	//"runtime"
	"syscall"

)

type hybridOverLay struct {
	ctx                 *context.Context
	runAsWindowsService bool
	appName             string
}

func (svc *hybridOverLay) Run() error {
	ctx, cancel := context.WithCancel(*svc.ctx)
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer func() {
		signal.Stop(exitCh)
		cancel()
	}()
	go func() {
		select {
		case s := <-exitCh:
			klog.Infof("Received signal %s. Shutting down", s)
			cancel()
		case <-ctx.Done():
		}
	}()
	return nil
}
