# UI Strategy - Beta 

Config driven UI driver for front end testing and for data seeding for cases where there isn't an REST API available and automated UI tests stored in declarative configuration files.

> Disclaimer: 
Part of strategy series :D - see [reststrategy](https://github.com/dnitsch/reststrategy) :wink: - there is a module and indpendently published CLI which functions in a similar way for REST calls. Always prefer to use that for any kind of configuration/data seeding where possible! 

> Only use this for last resort programatic interaction with a web app and for Front end testing

## Usage 


## Internals

Currently the entire loop through of Actions is using pointers to allow for an easier report builder output - adding concurency via go routines may be problematic and is not really desired at this point.

### Underlying Web Driver

This module and CLI use the [Go-Rod](https://github.com/go-rod/rod) which uses the CPD protocol.

## Tests

To run integration style tests you must have the sample app running, you can run it docker locally.

`docker run --name=pb-app --detach -p 8090:8090 dnitsch/reststrategy-sample:latest`

Then navigate to this [page](http://127.0.0.1:8090/_/?installer#)

Still TODO lots more tests


