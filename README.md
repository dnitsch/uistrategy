# UI Strategy - Beta

[![Go Report Card](https://goreportcard.com/badge/github.com/dnitsch/uistrategy)](https://goreportcard.com/report/github.com/dnitsch/uistrategy)

Config driven UI driver for front end testing, target agnostic, stored in declarative configuration files.

*Can be user for data seeding in cases where there isn't an REST API available.*

> Disclaimer:
Part of strategy series :D - see [reststrategy](https://github.com/dnitsch/reststrategy) :wink: - there is a [module](https://github.com/dnitsch/reststrategy/tree/main/seeder) and a [published CLI](https://github.com/dnitsch/reststrategy/releases) which functions in a similar way for REST calls. Always prefer to use that for any kind of configuration/data seeding where possible!

> Only use this for last resort programatic interaction with a web app and for **UI Testing**

## Features

The program splits the instructions into a top level slice which includes a navigation to a part of the app where it will include all the actions against the elements present. 

It will go through all of them in sequence as they are defined in the YAML.

Currently all the actions are performed against a single instance of the logged in page. for larger systems where order isn't important separate instances of the CLI can be triggered. 

- Authentication
  - optional authentication
    - supply username/password/submit elements - see [test/integration.yml](./test/integration.yml) for an example.
- Driving UI
  - page Navigation
  - element lookups using a selector either a CSS Style selector or XPath
    - CSSSelector will be tried first and then XPath
  - actions on element currently input + click/swipe
        - would be nice to include double click/right click, etc...
- Report
  - report with all steps
  - screenshots on errors attached to the step.
  - report.json - can be used in an HTML template creation, additionally JUnit or any other kind of format  can be parsed from that base.
- [ConfigManager](https://github.com/dnitsch/configmanager) integrated for easy storage of secrets in YMLs that can be committed - see this [example](./test/integration-with-configmanager.yml) of a password for auth.

As this is still in *beta* expect bugs and the interface to change

>Improvement/Feature:

- accept multiple yaml docs and run them in parallel.
- allow composing of complete strategies from multiple YAML documents - AVOID 1k+ YAML lines

## Usage

Download the correct binary for your architecture - [instructions](./docs/installation.md)

`uiseeder -h` for help

`uiseeder -i path/to/yaml -v`

See test/integration.yml for an [example with auth](#with-auth) you can use with pocketbase.io.

## Internals

Currently the entire loop through of Actions is using pointers to allow for an easier report builder output - adding concurency via go routines may be problematic and is not really desired at this point.

### Underlying Web Driver

This module and CLI use the [Go-Rod](https://github.com/go-rod/rod) which uses the CPD protocol.

## Example

### With Auth

Most scenarios for UI tests will require a login of sorts, for easy simulation of how this 

To run integration style tests you must have the sample app running, you can run it docker locally.

`docker run --name=pb-app --detach -p 8090:8090 dnitsch/reststrategy-sample:latest`

Then navigate to this [page](http://127.0.0.1:8090/_/?installer#)

Add you user name and password and replace it in the YAML with whatever you chose.

Explore any parts of the app and grab elements by XPath or CSS and add new page or element actions within a page.

## Help

always wanted and welcomed

- Still TODO lots more tests
- report formats and outcomes need fixing up
- any current/new features...