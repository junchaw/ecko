# ecko

<div align="center">
  <img src="docs/logo.png" alt="ecko logo" width="200">
</div>

[![Docker Pulls](https://img.shields.io/docker/pulls/junchaw/ecko.svg)](https://hub.docker.com/r/junchaw/ecko/)
[![Build Status](https://github.com/junchaw/ecko/workflows/Main/badge.svg?branch=main)](https://github.com/junchaw/ecko/actions)

HTTP echo server that returns and logs all request details for debugging.

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

#### `/api`

Print echo response as JSON.

#### `/status/{code}`

Return the specified HTTP status code.

## Deploy

Kubernetes:

```yaml
kubectl apply -f deploy/kubernetes-deploy-example.yaml
```
