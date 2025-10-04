package app

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/yanking/app-skeleton/pkg/log"
)

// IComponent server interface
type IComponent interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Name() string
}

// App servers
type App struct {
	name       string
	components []IComponent
}

// New create an app
func New(name string, components []IComponent) (*App, error) {
	log.Infof("app: all components initialized and dependencies injected.")
	return &App{
		name:       name,
		components: components,
	}, nil
}

// Run starts the application and manages the lifecycle of all components
func (a *App) Run() error {
	log.Infof("app: starting application: %s", a.name)

	appCtx, appStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appStop()

	// 使用 errgroup 并发启动所有组件
	g, ctx := errgroup.WithContext(appCtx)

	for _, component := range a.components {
		// 在闭包中使用局部变量捕获 component
		c := component
		g.Go(func() error {
			log.Infof("app: starting %s", c.Name())
			if err := c.Start(ctx); err != nil {
				log.Errorf("app: failed to start %s: %v", c.Name(), err)
				return err
			}
			return nil
		})
	}
	log.Infof("app: waiting for initial component start goroutines to complete...")
	// 等待所有组件启动完成或任一组件启动失败
	if err := g.Wait(); err != nil {
		log.Errorf("app: failed to start application: %v", err)
		return err
	}

	log.Infof("app: all components started successfully")

	<-appCtx.Done()

	errCause := context.Cause(appCtx)
	if errCause != nil && !errors.Is(errCause, context.Canceled) && !errors.Is(errCause, context.DeadlineExceeded) { // context.Canceled from signal is normal
		log.Infof("app: shutdown initiated due to: %v", errCause)
	} else {
		log.Infof("app: shutdown signal received or context cancelled normally.")
	}

	// --- Graceful Shutdown Procedure ---
	log.Infof("app: initiating graceful stop of application %s...", a.name)

	// Create a new context for the shutdown procedure itself, with a timeout.
	// This timeout is for the *entire* shutdown sequence of all components.
	shutdownOverallTimeout := 20 * time.Second // Configurable
	stopCtx, cancelStopCtx := context.WithTimeout(context.Background(), shutdownOverallTimeout)
	defer cancelStopCtx()

	// Stop components in reverse order of registration (LIFO).
	// This assumes a simple dependency order; more complex apps might need a dependency graph.
	for i := len(a.components) - 1; i >= 0; i-- {
		comp := a.components[i]
		log.Infof("app: attempting to stop component: %s", comp.Name())

		// We pass stopCtx to each component's Stop method.
		// The component's Stop method should respect this context's deadline.
		if err := comp.Stop(stopCtx); err != nil {
			log.Errorf("app: error stopping component %s: %v", comp.Name(), err)
		} else {
			log.Infof("app: component %s stopped successfully.", comp.Name())
		}
	}

	log.Infof("app: application %s stopped gracefully.", a.name)

	if errCause != nil && !errors.Is(errCause, context.Canceled) {
		return fmt.Errorf("application shutdown due to error: %w", errCause)
	}
	return nil // Graceful shutdown completed

}
