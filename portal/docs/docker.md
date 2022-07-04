# Docker

## Docker build

```bash
docker build -t rovast/za-portal .
```

## Run

```bash
docker run --name zaportal -p 9080:9080 -d rovast/za-portal:latest
```

You can visit http://127.0.0.1:9080

## Settings

### Expose ports

- 9080: frontend port
- 80: backend port

> Frontend will request `/api`, and will proxy by nginx, you can config it in `nginx.conf`

### Nginx configuration

- conf file: `/etc/nginx/nginx.conf`
- frontend dist path(www root): `/usr/share/nginx/html`

