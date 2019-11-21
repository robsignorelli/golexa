# golexa

_**golexa is very much an early work in progress. While there's enough
here currently to deploy a basic Alexa skill in Go, there are a ton
of rough edges that need to be shaved off and features that I need to
add before I would not be embarrassed by someone using this.**_

This is a framework that helps you ship Alexa skill code using
Go on AWS Lambda. Amazon has official SDKs for Java, Node, and Python, but
even though Go is a great language for FaaS workloads like this, they've
not released a Go SDK yet. This project tries to fill that gap as much
as possible.

Note: golexa does not help you set up your interaction model. You should
still set up your skill/model using Amazon's developer console at
https://developer.amazon.com/alexa/. golexa only takes the sting out
of handling incoming requests and responding w/ as little code as possible.

## Getting Started

```
go get github.com/robsignorelli/golexa
```

## Basic Usage

Let's pretend that you defined an "echo" skill that has these
two interactions; one that has Alexa say hello to someone and one
that has her say goodbye to someone. In the Alexa console, we might
have set up 2 intents with the following sample utterances:

```text
HelloIntent
  Say hello to {name}
  Say hi to {name}

GoodbyeIntent
  Say goodbye to {name}
  Say bye to {name} 
```

This is all it takes to have a fully-baked skill that handles
both interactions.

```go
package main

import (
    "context"
    "github.com/robsignorelli/golexa"
)

func main() {
    skill := golexa.Skill{}

    skill.RouteIntent("HelloIntent", func(ctx context.Context, req golexa.Request) (golexa.Response, error) {
        name := req.Body.Intent.Slots.Resolve("name")
        return golexa.NewResponse().Speak("Hello " + name).Ok()
    })
    skill.RouteIntent("GoodbyeIntent", func(ctx context.Context, req golexa.Request) (golexa.Response, error) {
        name := req.Body.Intent.Slots.Resolve("name")
        return golexa.NewResponse().Speak("Goodbye " + name).Ok()
    })
	
    golexa.Start(skill)
}
```

Obviously, you wouldn't normally write all of your code in `main()`, but
it shows you how little you need to write in order to get a working skill.
For an ever-so-slightly more complex example, you can look in the
[sample/](https://github.com/robsignorelli/golexa/tree/master/sample) directory
for a simple TODO list skill.

## Middleware

There are some units of work you want to execute in most/all of your
intent handlers. For instance you might want to log every incoming
request or validate that a user has linked their Amazon account to 
your system before doing any real work. All of this can be done using
middleware, similar to how you might do this in a REST API.

```
func main() {
    middleware := golexa.Middleware{
        LogRequest,
        ValidateUser,
    }
    
    // Log & authenticate the add/remove intents, but not the status intent.
    service := FancyService{}
    skill.RouteIntent("FancyAddIntent", middleware.Then(service.Add))
    skill.RouteIntent("FancyRemoveIntent", middleware.Then(service.Remove))
    skill.RouteIntent("FancyStatusIntent", service.Status)
    golexa.Start(skill)
}

func LogRequest(ctx context.Context, request Request, next HandlerFunc) (Response, error) {
    fmt.Println("... log something interesting ...")
    return next(ctx, request)
}

func ValidateUser(ctx context.Context, request Request, next HandlerFunc) (Response, error) {
    if request.Session.User.AccessToken == "" {
        return golexa.NewResponse().Speak("No soup for you.").Ok()
    }
    return next(ctx, request)
}
```



## Future Enhancements

Here are a couple of the things I plan to bang away at. If you have any
other ideas that could help you in your projects, feel free to add
an issue and I'll take a look.

* Tests
* Some sort of templating to make generating speech responses easier
* Echo Show display template directive support
* Name free interactions through `CanFulfillIntentRequest`

Because this is still very much a work in progress, I can't promise that
I won't make breaking changes to the API while I'm still trying to shake
this stuff out.