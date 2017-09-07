gluttony
---------

[![Build Status](https://travis-ci.org/cryptounicorns/gluttony.svg?branch=master)](https://travis-ci.org/cryptounicorns/gluttony)

Receives data from user configured message queue and writes it into user configured persistent storage.

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

## Name origin

> Gluttony can eat virtually __any substance__ in __any amount at super-speed__. He can consume an unlimited amount of matter in any form - solid, liquid or gas without experiencing harmful effects regardless of what he consumes or how much he consumes.

http://fma.wikia.com/wiki/Gluttony
