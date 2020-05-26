# golexa

_**golexa is very much an early work in progress. While there's enough
here currently to deploy a basic Alexa skill in Go, there are a ton
of rough edges that need to be shaved off and features that I need to
add before I would not be embarrassed by someone using this.**_

golexa is a framework that helps you ship Alexa skill code using
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
        return golexa.NewResponse(req).Speak("Hello " + name).Ok()
    })
    skill.RouteIntent("GoodbyeIntent", func(ctx context.Context, req golexa.Request) (golexa.Response, error) {
        name := req.Body.Intent.Slots.Resolve("name")
        return golexa.NewResponse(req).Speak("Goodbye " + name).Ok()
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

The `golexa` framework provides a couple of middleware functions out of the
box for logging and making sure that the user has an access token (i.e. set up account linking),
but you can easily create your own by implementing a function with the proper signature.

```go
func main() {
    mw := golexa.Middleware{
        middleware.Logger(),
        middleware.RequireAccount(),
        CustomMiddleware
    }
    
    // Log, authorize, and do custom work on the add/remove intents, but not the status intent.
    service := NewFancyService()
	skill := golexa.Skill{}
    skill.RouteIntent("FancyAddIntent", mw.Then(service.Add))
    skill.RouteIntent("FancyRemoveIntent", mw.Then(service.Remove))
    skill.RouteIntent("FancyStatusIntent", service.Status)
    golexa.Start(skill)
}

func CustomMiddleware(ctx context.Context, request Request, next HandlerFunc) (Response, error) {
    // ... do something awesome before the request ...
    res, err := next(ctx, request)
    // ... do something awesome after the request ...
    return res, err
}
```

## Templates

Chances are that most of your intents have some sort of standard format/template for how you want
to respond to the user. It's super easy to define speech templates using standard Go templating
and evaluate them with whatever runtime data you want. In this example, the `.Value` of the template
context is just a string, but you can make it whatever you want for more complex responses.

```go
greetingHello := speech.NewTemplate("Hello {{.Value}}")
greetingGoodbye := speech.NewTemplate("Goodbye {{.Value}}")

skill := golexa.Skill{}
skill.RouteIntent("GreetingIntent", func(ctx context.Context, req golexa.Request) (golexa.Response, error) {
    name := req.Body.Intent.Slots.Resolve("name")

    switch {
    case shouldSayHello(req):
        return golexa.NewResponse(req).SpeakTemplate(greetingHello, name).Ok()
    default:
        return golexa.NewResponse(req).SpeakTemplate(greetingGoodbye, name).Ok()
    }
})
```

You can also inject custom functions into your templates to perform more processing as you see fit by passing
in any number of `WithFunc()` options into your `NewTemplate()` call.

```go
greetingHello := speech.NewTemplate("Hello {{.Value | jumble}}",
    speech.WithFunc("jumble", jumbleText))

...

func jumbleText(value string) string {
    // ... shuffling logic goes here ...
}
```

## Templates: Multi-Language Support

Why limit yourself to just English? Golexa speech templates provide simple hooks to support ANY of the
locales/languages that Alexa supports ([see full list here](https://developer.amazon.com/en-US/docs/alexa/custom-skills/develop-skills-in-multiple-languages.html#h2-code-changes)).

```go
greetingHello := speech.NewTemplate("Hello {{.Value}}",
    speech.WithTranslation(language.Spanish, `Hola {{.Value}}`),
    speech.WithTranslation(language.Italian, `Ciao {{.Value}}`),
)
```

Now when you respond w/ this template and the `name` slot is "Rob", Alexa
will say "Hola Rob" when the locale is "es", "es-MX", "es-ES", etc. It
will cay "Ciao Rob" when the locale is any locale that starts with "it"
and it will fall back to the English translation of "Hello Rob" for
any other language.

In this example, I use the same translation for all Spanish variants, but
you can just as easily support different language variants, too:

```go
// You could use "language.LatinAmericanSpanish" and "language.EuropeanSpanish"
// instead of parsing, but I wanted to show you the actual locale names.
greetingHello := speech.NewTemplate("Hello {{.Value}}",
    speech.WithTranslation(language.MustParse("es-MX"), `Hola {{.Value}}`),
    speech.WithTranslation(language.MustParse("es-ES"), `Hola from Spain, {{.Value}}`))
)
```

# Back-and-Forth Interactions w/ ElicitSlot

You want your interactions to be as friendly to your users as possible. For instance, you might
have an interaction/intent where users might want to add an item to their list that you're 
maintaining for them. You might want to support both of these phrases:

```
AddItemIntent
  Add {item_name} to my list
  Update my list
```

In the first case the user provides the name of the item they want to add, so you have all of
the information you need to complete the request. In the latter case, you don't, so you want
to have Alexa prompt the user for that information and then trying again. Here's how you can
use `golexa` to fulfill the request when you have everything, and "elicit slot" when you don't.

```go
skill.RouteIntent("AddItemIntent", func(ctx context.Context, req golexa.Request) (golexa.Response, error) {
    itemName := req.Body.Intent.Slots.Resolve("item_name")

    // They said "update my list", so have their device ask them for the
    // item name and send the result back to this intent again.
    if itemName == "" {
        return golexa.NewResponse(req).
            Speak("What would you like me to add to the list?").
        	ElicitSlot("AddItemIntent", "item_name").
        	Ok()
    }

    // They either specified the "item_name" slot in the initial request or
    // we were redirected back here after an ElicitSlot.
    addItemToList(req.UserID(), itemName)
    return golexa.NewResponse(req).
        Speak("Great! I've added that to your list.").
        Ok()
})
```

## Future Enhancements

Here are a couple of the things I plan to bang away at. If you have any
other ideas that could help you in your projects, feel free to add
an issue and I'll take a look.

* Echo Show display template directive support
* Name free interactions through `CanFulfillIntentRequest`

Because this is still very much a work in progress, I can't promise that
I won't make breaking changes to the API while I'm still trying to shake
this stuff out.