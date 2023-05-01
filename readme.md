# ReProxied

ReProxied is a middleware plugin for [Traefik](https://github.com/traefik/traefik) to route an incoming request through a proxy.
Be aware that this middleware initiates the call to the proxy and any middlewares after this one will be skipped. If the request to the proxy itself fails the middleware will respond with a 502 bad gateway response.

When set to `true` the parameter `keepHostHeader` allow to keep original Host as HTTP header even if proxied request target any other host.

## Configuration

### Static

```yaml
pilot:
  token: "xxxxx"

experimental:
  plugins:
    reproxied:
      moduleName: "github.com/nilskohrs/reproxied"
      version: "v0.0.5"
      keepHostHeader: true|false # optional, false by default
      logLevel: 1,2,3,4,5  # optional, 2 by default. 1 -> debug, 4 -> error, 5 -> off
```

### Dynamic

```yaml
http:
  middlewares:
    reproxied-foo:
      reproxied:
        proxy: http://proxyHost:3128
        targetHost: https://example.com
        keepHostHeader: true|false # optional, false by default
        logLevel: 1,2,3,4,5  # optional, 2 by default. 1 -> debug, 4 -> error, 5 -> off
```

> log level are defined in [Levels.go](./internal/logging/Levels.go)
