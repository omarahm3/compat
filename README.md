# Compat

Have you ever wondered why there is no inheritance on docker-compose yaml syntax? and why am i writing the same stuff again and again for some services? well .. i had the same questions, and i couldn't find any solution for it, that is why i created this small tool to allow reusability on docker-compose files.

It is so simple as doing:

```
file: compat.yaml
version: '3.2'
services:
  base_service:
    deploy:
      resources:
        limits:
          cpus: '0.1'
          memory: 128M
        reservations:
          memory: 128M

  mongo:
    inherit: base_service
    container_name: database
    image: mongo
```

that will get converted into:

```
file: docker-compose.yaml
version: '3.2'
services:
  mongo:
    container_name: database
    deploy:
      resources:
        limits:
          cpus: '0.1'
          memory: 128M
        reservations:
          memory: 128M
    image: mongo
```

Or check [./example](./example) to know how to inherit multiple services.

## Why not `extends` though?

It depends on the use case, to use [extends](https://docs.docker.com/compose/extends) you must have a common "service" not "configuration" while for `inherit` you're just sharing configuration on whatever services.

## How?

Compat depends on a new file `compat.yaml` that is in the exact same syntax as `docker-compose.yaml` for only single purpose which is to allow the use of `compat` service property.
Running `compat` will actually parse the file by looking for `inherit` and for any `base_*` services, will reuse whatever under `base_*` and add them to services.

## Install

If you already have golang installed:

```
go install github.com/omarahm3/compat@latest
```

## TODO

- [X] Tests
- [X] Allow multiple inheritance?
- [ ] Add release workflow for automated releases
