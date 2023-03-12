[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=dnitsch_uistrategy&metric=bugs)](https://sonarcloud.io/summary/new_code?id=dnitsch_uistrategy)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=dnitsch_uistrategy&metric=sqale_index)](https://sonarcloud.io/summary/new_code?id=dnitsch_uistrategy)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=dnitsch_uistrategy&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=dnitsch_uistrategy)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=dnitsch_uistrategy&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=dnitsch_uistrategy)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=dnitsch_uistrategy&metric=coverage)](https://sonarcloud.io/summary/new_code?id=dnitsch_uistrategy)

# UI Strategy - Beta

[![Go Report Card](https://goreportcard.com/badge/github.com/dnitsch/uistrategy)](https://goreportcard.com/report/github.com/dnitsch/uistrategy)

Config driven UI driver for front end testing, target agnostic, stored in declarative configuration files.

*Can be user for data seeding in cases where there isn't an REST API available.*

> Disclaimer:
Part of strategy series :D - see [reststrategy](https://github.com/dnitsch/reststrategy) :wink: - there is a [module](https://github.com/dnitsch/reststrategy/tree/main/seeder) and a [published CLI](https://github.com/dnitsch/reststrategy/releases) which functions in a similar way for REST calls. Always prefer to use that for any kind of configuration/data seeding where possible!

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

As this is still in *beta* expect bugs and the interface to change. 

>Improvement/Feature:

- accept multiple yaml docs and run them in parallel.
- allow composing of complete strategies from multiple YAML documents - AVOID 1k+ YAML lines

## Configuration

### `setup`

Top level config item to initiate a webBrowser session

#### `baseUrl`

must be provided acts as a baseUrl for all navigations, including login.
  
#### `continueOnError`

Default: false

Will stop execution if an error occurs, useful to set this to 
  true if a large execution sequence 

### `auth`

... 

### `actions`

Is a list of actions to execute - the order in which they are provided. 

Single `action` block has the below structure, at this level the action is only a navigation action/view action. it can contain 0 or more `elementActions`

#### `name` (required)

Name of the view action - will be used in reports

#### `navigate` (required)

The path to append to the baseUrl to navigate to perform actions against elements on that page.

>navigate string is appended to the baseUrl without any slashes - ensure you either specify the baseUrl with a trailing slash or all your `navigate`s should include a preceeding slash.

#### `iframe` (optional)

object with following properties. when the actions you want to perform on that page/view are within an iframe it must be specifed here

##### `selector`

the selector for the iframe - using either CSS or Xpath e.g.: `(//*/iframe)[1]` - i.e. give me the first iframe on the page

##### `waitEval`

oftentimes older apps (e.g. timesheet portals :wink: ) include iframes and they are often loaded slower to avoid losing the context, specify the eval to wait for contents inside the iframe. 

this could be a `myVar !== null` - more info in the godoc or below.

```go
// IframeAction 
type IframeAction struct {
  Selector string `yaml:"selector,omitempty" json:"selector,omitempty"`
  // WaitEval has to be in the form of a boolean return
  // e.g. `myVar !== null` or `(myVar !== null || document.title == "ready")`
  // the supplied value will be appended to an existing
  // `return document.readyState === 'complete' && ${WaitEval};`
  WaitEval string `yaml:"waitEval,omitempty" json:"waitEval,omitempty"`
}
```

#### `elementActions`

list of actions to perform within the page/view, each `elementAction` has the following structure

#### `name` (required)

Name of the action on the element - will be used in reports

#### `element` (required)

the element to identifier - an object with below attrs

##### `selector` (required)

CSS or XPath style selector to attempt to locate the element on the page.

if not found and running `continueOnError` mode the execution will move on to the next element in the sequence.

##### `value` (optional)

if value is not provided it will be a click type action, if value is provided it will be an input type action

##### `assert` (bool)

Defaults to false.

When running in UI test mode only this should be set to true...

When set to true elements presence is only asserted and any input/click actions will be skipped.

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