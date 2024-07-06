mainpkg
=======

Service runner for managing gRPC and HTTP services.

Running Services
----------------

Setting up a service (in the "microservices" sense) is done using `mainpkg.Setup` with a number of configuration options. See options.go for service configuration options.

### gRPC service

Each gRPC service should implement `mainpkg.ServiceRegisterer` interface. Optionally, any service exposed using gRPC gateway for REST requests should also implement `mainpkg.GatewayRegisterer` interface.

A simple gRPC service with gRPC gateway enabled can be spun up using:

```go
func main() {
	// set up viper config, logging logger, and services

	srv, err := mainpkg.Setup(config, logger,
		mainpkg.WithDualService(mainpkg.External, services...),
	)
	if err != nil {
		logging.WithError(err).Fatal("server setup")
	}
	srv.Run()
}
```

### HTTP service

A standard `http.Handler` can be wrapped using `mainpkg.HTTPOnly` to serve all non-CORS requests on the handler. Helper is included for gorilla/mux routers, as they expose matching functionality.

A simple HTTP service can be spun up using:

```go
func main() {
	// set up viper config, logging logger, and HTTP handler

	srv, err := mainpkg.Setup(config, logger,
		mainpkg.WithWebHandler(mainpkg.HTTPOnly(handler)),
	)
	if err != nil {
		logging.WithError(err).Fatal("server setup")
	}
	srv.Run()
}
```

### Multiple servers/ports

Exposing different services on different ports is a common requirement, especially for keeping internal services from being exposed to the world. Loading different interceptors (ex: api-key vs admin-only auth) enables security tailored to the running service.

Using `mainpkg.AdditionalServers` option, any number of serving ports can be managed on a single `Run` call.

```go
func main() {
	// set up viper config, logging logger, gRPC services, and HTTP handler

	// HTTP handler for UI and webhooks.
	httpSrv, err := mainpkg.Setup(config, logger,
		mainpkg.WithPort(config.GetInt("http.port")),
		mainpkg.WithWebHandler(mainpkg.HTTPOnly(handler)),
	)
	if err != nil {
		logging.WithError(err).Fatal("http server setup")
	}

	// Internal gRPC-only service, no gateway.
	adminSrv, err := mainpkg.Setup(config, logger,
		mainpkg.WithPort(config.GetInt("admin.port")),
		mainpkg.WithServices(mainpkg.External, adminServices...),
	)
	if err != nil {
		logging.WithError(err).Fatal("internal server setup")
	}

	// External gRPC service, with gRPC-gateway.
	srv, err := mainpkg.Setup(config, logger,
		// Port configured with environment variable "SERVER_PORT" or key "server.port" in config.

		// Register service with gRPC-gateway REST endpoints.
		mainpkg.WithDualService(mainpkg.External, services...),

		// Serves specific endpoints as HTTP, for doing redirects.
		// All other requests will be matched against gRPC-gateway or gRPC handler.
		mainpkg.WithWebHandler(redirSvc),

		// Add other servers.  All will be started/stopped as one.
		mainpkg.AdditionalServers(httpSrv, adminSrv),
	)
	if err != nil {
		logging.WithError(err).Fatal("internal server setup")
	}

	// All three services start, along with the debug server.
	// All services are
	srv.Run()
}
```

Scheduled Tasks
---------------

Background tasks can be scheduled using the Cron Options. Each pod can schedule background tasks using `WithCron` option.

When multiple replicas are deployed, these scheduled tasks can be de-duplicated using leader election. The leader options `WithLeaderCron` and `WithLeaderFunc` run that scheduled task on only one pod.

### Task Functions

**WARNING** every task function MUST respect context cancelation, exiting the task as soon as possible. Failure to do so will leak resources.

Debugging
---------

Running services can get into inconsistent/hang states and need to be debugged without restart. Along with running the services on their configured ports, a debug port is run (port 12000) for help in these situations.

### `/debug/pprof` Go Profiler

Connect and use go pprof tooling.

### `/log/trigger` Read and Set the service logger level

HTTP GET request will respondtwith the currently set log level.

HTTP POST will set or toggle between Info, Debug, and Trace levels.  
Use query parameter "level" to set a specific level (ex: `/log/trigger?level=warn`\)  
Response will be the new log level.

Useful when more verbose logs are temporarily needed for finding the root cause of an error.

### `/healthz/` General Service Health

Returns a 200(OK) when service is running and returns 503(Unavailable) when service is shutting down. Useful for kubernetes startup/liveliness probes.

### `/healthz/{name}` Detailed Service Health

Init options can be used to signal readiness when the init function(s) completes. Useful for kubernetes liveliness/readiness probes.

#### Signal Readiness

Initializing with:

```go
mainpkg.InitFunc("customname", func(ctx context.Context) error {
	// initialize service - ex: load in-mem cache from disk/database/backend or parse HTML templates
	// non-nil error signals a failure to initialize.  Service will shut down.
	return err
})
```

Matching readiness probe, for signaling:

```yaml
readinessProbe:
  httpGet:
    path: /healthz/customname
    port: 12000
  initialDelaySeconds: 5
  periodSeconds: 5
```

##### Special Case: Parent Service Init Function

**CAUTION**: This setup can result in a crash loop if the init function does not complete before the failed liveliness probe causes the pod to be killed. Use a readiness probe to handle awaiting initialization, when possible.

Adding an init function for the top-level service will cause the service not to show ready on `/healthz/` until that initialization function completes successfully.

Using init function on the parent service example:

```go
mainpkg.InitFunc("", func(ctx context.Context) error {
	// Service is not ready until this function completes.
	return err
})
```

Liveliness probe will fail until the init function completes:

```yaml
livelinessProbe:
  httpGet:
    path: /healthz/
    port: 12000
  initialDelaySeconds: 5
  periodSeconds: 5
```

#### Monitor Readiness

Signaling unhealthy state can be done with `server.Ready("customname", mainpkg.NotReady)` as soon as a failure state is identified.

Within a cron or leader function, the helper function `ServerReady` will set the readiness state on the parent server.

Readiness example - here are two ways to set up a database check every five seconds and signal kubernetes via readiness probe if the database is temporarily unavailable:

```go
// Option 1:  Monitoring using a cron function.
srv, err := mainpkg.Setup(config, logger,
	// name the cron job (tagged logs) - trigger every 5 seconds.
	mainpkg.WithCron("db monitor", mainpkg.NewCronTicker(5*time.Second),
		// Set a context timeout of 1 second for the function.
		mainpkg.WithFuncTimeout(1*time.Second,
			func(ctx context.Context) error {
				// Check if the database is connected.
				err := db.PingContext(ctx)
				if err != nil {
					mainpkg.ServerReady(ctx, "customname", mainpkg.NotReady)
				} else {
					mainpkg.ServerReady(ctx, "customname", mainpkg.Ready)
				}
				return err // returned error logged when non-nil
			})),
)

// Option 2:  Monitoring in a background goroutine.
go func() {
	// every 5 seconds
	t := time.Tick(5 * time.Second)
	for {
		select {
		// Need to create a context that is canceled when srv.Run() returns, to exit properly.
		case <-ctx.Done():
			return
		case <-t:
		}
		// Set a context timeout of 1 second for the ping.
		pctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		// Check if the database is connected.
		if err := db.PingContext(pctx); err != nil {
			srv.Ready("customname", mainpkg.NotReady)
			logging.WithError(logger, err).Error("db monitor")
		} else {
			srv.Ready("customname", mainpkg.Ready)
		}
		cancel()
	}
}()

// Start the service.
srv.Run()
// context cancellation for Option 2 after `Run` exits
```
