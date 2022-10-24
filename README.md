# UI Strategy

Config driven UI driver for data seeding for cases where there isn't an REST API available and automated UI tests stored in declarative configuration files.

Part of strategy series :D - see reststrategy :wink:

Run local tests

`docker run --name=pb-app --detach -p 8090:8090 dnitsch/reststrategy-sample:latest`

Then navigate to this [page](http://127.0.0.1:8090/_/?installer#)

## Underlying Web Driver

This module and CLI use the [Go-Rod](https://github.com/go-rod/rod) which uses the CPD protocol.
