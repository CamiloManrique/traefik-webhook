# traefik-url-extractor

A [Traefik](https://traefik.io) middleware plugin that extracts URL path parameters into custom request headers using named regex capture groups.

## Usage

### Static configuration

```yaml
experimental:
  plugins:
    url-extractor:
      moduleName: github.com/CamiloManrique/traefik-url-extractor
      version: v1.0.0
```

### Dynamic configuration

```yaml
http:
  routers:
    my-router:
      rule: host(`api.example.com`)
      service: my-service
      entryPoints:
        - web
      middlewares:
        - extract-params

  services:
    my-service:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:8080

  middlewares:
    extract-params:
      plugin:
        url-extractor:
          regex: '/users/(?P<user>[a-f0-9-]+)'
          headers:
            X-User-Id: user
```

## Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `regex` | string | yes | Regular expression with named capture groups (`(?P<name>...)`) to match against the request URL |
| `headers` | map | yes | Mapping of `header-name → capture-group-name`. Each matched group is set as a request header. |

If the URL does not match the regex, the request is passed through unchanged.

## Local Mode

```yaml
# Static configuration
experimental:
  localPlugins:
    url-extractor:
      moduleName: github.com/CamiloManrique/traefik-url-extractor
```

Place the plugin source under:

```
./plugins-local/src/github.com/CamiloManrique/traefik-url-extractor/
```
