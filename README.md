gluttony
---------

[![Build Status](https://travis-ci.org/cryptounicorns/gluttony.svg?branch=master)](https://travis-ci.org/cryptounicorns/gluttony)

## Development

All development process accompanied by containers. Docker containers used for development, Rkt containers used for production.

> I am a big fan of Rkt, but it could be comfortable for other developers to use Docker for development and testing.

## Requirements

- [docker](https://github.com/moby/moby)
- [docker-compose](https://github.com/docker/compose)
- [jq](https://github.com/stedolan/jq)
- [rkt](https://github.com/coreos/rkt)
- [acbuild](https://github.com/containers/build)

### Running gluttony

Build a binary release:

``` console
$ GOOS=linux make
# This will put a binary into ./build/gluttony
```

#### Docker

``` console
$ docker-compose up gluttony
```

#### Rkt

There is no rkt container for this service at this time.

#### No isolation

``` console
$ go run ./gluttony/gluttony.go --debug
```
