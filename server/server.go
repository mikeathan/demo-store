package server

import (
	"demo-store/endpoints"
	"demo-store/store"
	"demo-store/users"
	"demo-store/utils"
	"errors"
	"fmt"
	"net/http"
)

func Listen(port int, depth int) error {

	userDatabase := users.Load("cache")
	kvStore := store.CreateKvStore(utils.ApplicationTracer(), userDatabase, depth)
	shutdownListener := store.CreateShutdownListener()
	kvStore.RegisterShutdownListener(shutdownListener)

	register(kvStore)

	return start(port, *shutdownListener)
}

func start(port int, shutdownListener store.ShutdownListener) error {

	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	go func() {

		utils.ApplicationTracer().LogInfo("HTTP Server Listening ", httpServer.Addr)
		resp := <-shutdownListener.Listener
		if resp {
			shutdown(httpServer)
		}
	}()

	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		utils.ApplicationTracer().LogInfo("HTTP server error: ", err)
		return err
	}

	utils.ApplicationTracer().LogInfo("HTTP server closed")
	return nil
}

func shutdown(httpServer *http.Server) {
	if err := httpServer.Close(); err != nil {
		utils.ApplicationTracer().LogInfo("HTTP close error: ", err)
	}
}

func register(store store.Store) {

	routes := endpoints.APIRoutes(utils.HttpTracer(), store)
	for _, route := range routes.Secure {
		registerRoute(route)
	}
	for _, route := range routes.Insecure {
		registerRoute(route)
	}
}

func registerRoute(route endpoints.Route) {
	utils.ApplicationTracer().LogInfo("Register route: ", route.RootPath())
	http.Handle(route.RootPath(), route)
}
