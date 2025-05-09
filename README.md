# ECKO

<div align="center">
  <img src="docs/logo.png" alt="ecko logo" width="200">
</div>

[![Docker Pulls](https://img.shields.io/docker/pulls/junchaw/ecko.svg)](https://hub.docker.com/r/junchaw/ecko/)
[![Release a tag](https://github.com/junchaw/ecko/actions/workflows/release.yml/badge.svg)](https://github.com/junchaw/ecko/actions/workflows/release.yml)

HTTP echo server that returns and logs all request details for debugging.

> Why not [traefik/whoami](https://github.com/traefik/whoami)? Because ecko logs request details, I find it extremely useful when developing webhooks where I want to know what's being sent to me.

## Usage

Run the server:

```shell
docker run -p 8080:8080 junchaw/ecko -h
docker run -p 8080:8080 junchaw/ecko
```

### Endpoints

#### `/`

Print help.

#### `/echo`

Print echo response.

Sample response:

```text
Status: 200
Hostname: 8de1402a4966
Name: ecko
RemoteAddr: 172.17.0.1:33746
IP[0]: 127.0.0.1
IP[1]: ::1
IP[2]: 172.17.0.2
Method: POST
URL: /echo
Host: localhost:8080

Headers:
Accept: */*
Postman-Token: 5e92276c-6950-4838-9f1f-bb9a0ef0186f
Accept-Encoding: gzip, deflate, br
Connection: keep-alive
Content-Length: 13
Content-Type: application/json
User-Agent: PostmanRuntime/7.43.3

Request Body:
{"hello":"world"}
```

#### `/api`

Print echo response as JSON.

Sample response:

```json
{
    "hostname": "8de1402a4966",
    "ip": [
        "127.0.0.1",
        "::1",
        "172.17.0.2"
    ],
    "headers": {
        "Accept": [
            "*/*"
        ],
        "Accept-Encoding": [
            "gzip, deflate, br"
        ],
        "Connection": [
            "keep-alive"
        ],
        "Content-Length": [
            "13"
        ],
        "Content-Type": [
            "application/json"
        ],
        "Postman-Token": [
            "39248df1-7ed4-4a90-887a-68f431aa7223"
        ],
        "User-Agent": [
            "PostmanRuntime/7.43.3"
        ]
    },
    "url": "/api",
    "host": "localhost:8080",
    "method": "POST",
    "name": "ecko",
    "remoteAddr": "172.17.0.1:33746",
    "requestBody": "{\"hello\":\"world\"}"
}
```

#### `/status/{code}`

Return the specified HTTP status code.

## Deploy

Kubernetes:

```yaml
kubectl apply -f deploy/kubernetes-deploy-example.yaml
```
