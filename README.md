# REST application integration test tool

## Use case

When you need to test REST API with dependent calls of different endpoints, collect information from HTTP header and/or body (*supported only json*) of the responses and use in the next calls.
You can simply define the calls in one file with specific json structure (the structure will be described below) and run it with the tool.  

## Simple usage

To use the tool your root folder with test cases has to be in the same directory as the tool for example:
```
-some_dir
- - rest.int.test
- - test_cases
- - - test1.json
- - - test2.json
- - - sub_tests
- - - - sub_test1.json
- - - - sub_test2.json
```
To run all the tests you need to call the command:
```
./rest.int.test -suites test_cases
```
In this case the application will scan all folders and sub folders into test_cases, collect all json files and try to run the suites inside the files if they have appropriate format

## Docker usage

There is a docker image into public repository
```
yrashetska/intest:latest
```
It contains the latest version of the application. To run the test suites on the container you can just specify volume with mapping to  /usr/local/bin/<folder_name> in container.
For example:
```
docker ... -v ~/go/src/project/tests:/usr/local/bin/tests ...
```
If you need to test some application in local docker container and it's run with docker-compose the you need to specify 
```yaml
networks:
  app-network :
    name: app-network-dev
```
in docker compose and set the network withing --net flag when run docker test tool:
```
--net app-network-dev
```
Also don't forget that you need to link your service to the container
```
--link some-service:some-service-alias
```
And the use the `some-service-alias` into test suites instead of localhost

Full example of running docker:
```
docker run --rm -v ~/go/src/project/tests:/usr/local/bin/tests --net app-network-dev --link some-service:some-service-alias yrashetska/intest:latest -suites tests
```

## Test suites

Basically one test suite restricted by file scope. What that means. Each file contains scope of variables and several tests. All tests in one file share the variable scope and can fill it up.
Besides that each test contains extract section (to fill up the variables scope) and asserts section where can be checked some simple conditions within response data (which was expracted before)

Let's take a look at simple suite file:
```json
{
  "description" : "Suite description.",
  "executor" : "CURL",
  "vars" :  {
    "host" : "localhost"
  },
  "tests" : [
    {
      "label" : "Test label",
      "command" : "${host}:8080/some/endpoint",
      "extracts" : [
        {
          "header" : "Status",
          "var" : "header-status"
        }  
      ],
      "assertions" : [
        {
          "var" : "header-status",
          "eq" : "OK"
        }
      ]
    }
  ]
}
```
This simple example contains only one test. It does GET request to localhost:8080/some/endpoint, extracts Status from HTTP header into variable with name `header-status` and checks if the status is 'OK'

NOTE:
- *"executor" : "CURL"* - At the moment the tool has only one implementation of the endpoints caller, it's curl (https://curl.haxx.se/)
- *"description" : "Suite description."* - Suite description - shown on test report
- *"label" : "Test label"* - Test description - shown on test report

### Variables
1. The first thing which we have to know about variables is that it's will be parsed into flat map:
```json
{
...
  "vars" : {
    "a" : "One",
    "b" : {
      "sub-1" : 1,
      "sub-2" : 2
    },
    "c" : ["c-1", "c-2", "c-3"],
    "d" : [
      {
        "one" : 1,
        "two" : 2
      },
      {
        "ten" : 10
      }
    ]
  },
...
}

```   
Above structure is parsed into next map (map[string]string):
```go
map["a"] = "One"
map["b.sub-1"] = "1"
map["b.sub-2"] = "2"
map["c.0"] = "c-1"
map["c.1"] = "c-2"
map["c.2"] = "c-3"
map["d.0.one"] = "1"
map["d.0.two"] = "2"
map["d.1.ten"] = "10"
```
And the values of the map can be used on "command" using `${...}`. For example `${d.0.two}`

2. Second is that we can use other variables to define some variable during initial loading (cyclic dependencies will cause a panic).
For example:
 ```json
 {
 ...
   "vars" : {
     "a" : "one",
     "b" : {
       "sub-1" : "${a} + two",
       "sub-2" : "${b.sub-1} = three"
     },
 ...
 }
 
 ```   
The corresponding variables in scope will be parsed into:
```go
map["a"] = "one"
map["b.sub-1"] = "one + two"
map["b.sub-2"] = "one + two = three"
```

3. When some variable is populated in extracts section of some test previous value will be lost. Even if there is no corresponding field in header or body
4. Variables are replaced only in command and can't be used in assertions (maybe it it's implemented later)

### Extracts
To extract some field from header we just need to specify header name and variable where we need to extract the value:
```json
{
  "header" : "Status",
  "var" : "header-status"
}  
```
The value will be extracted and trimmed, then it will be saved into scope with name `header-status` and can be used into assertions or even in the next request `${header-status}`

To extract some information from body we need to know exact path to the field and the we can use the path as it's done for variables:
```json
{
  "response" : {
    "body" : {
      "with" : {
        "array" : [
          "One",
          "Two",
          "Three"
        ]
      }
    }
  }
}
```  
To extract second element of inner array we can set next extract option:
```json
{
  "body" : "response.body.with.array.1",
  "var" : "var-from-body"
}  
```

### Asserts
The asserts can be done only on current variable scope. Asserts for the test executed directly after command and all extracts (to use up to date variable scope)
Available assertions:
```json
{
...
	{
	  "var" : "header-var",
	  "eq" : "OK"
	},
	{
	  "var" : "header-var",
	  "match" : "match to some reg exp [0-9a-z]+"
	},
	{
	  "var" : "body-var",
	  "lt" : 100
	},
	{
	  "var" : "body-var",
	  "gt" : 100
	},
	{
	  "var" : "body-var",
	  "lte" : 100
	},
	{
	  "var" : "body-var",
	  "gte" : 100
	},
	{
	  "var" : "body-var",
	  "and" : [
		{"gte" : 100},
		{"lte" : 200}
	  ]
	},
	{
	  "var" : "body-var",
	  "or" : [
		{"gte" : 200},
		{"lte" : 100}
	  ]
	},
	{
	  "var" : "body-var",
	  "lte" : 100,
	  "gte" : 100
	}
...
}

``` 
- *eq* - true if equal
- *match* - if string matches to reg exp
- *lt* - true if the value less than
- *gt* - true if grater than
- *lte* - lt or eq
- *gte* - gt or eq
- *or* - logical operator, one of condition has to be true
- *and* - logical operator, all of conditions must be true

# TODO:
- [ ] Add filters to run only few test suites
- [ ] Test report in HTML by flag
- [ ] Not only curl runner (I'm not sure if it's really needed) 
- [ ] Improve unit test coverage