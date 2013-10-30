GoReq
=======

Simple and sane HTTP request library for Go language.

Why GoReq?
==========

Go has very nice native libraries that enables you to do lots of cool things. But sometimes those libraries are too low level, which means that to do a simple thing, like an HTTP Request, takes some time. And if you want to do something as simple as adding a timeout to the request, means that you need to write several lines of code.
This is why we thing GoReq is useful. Because you can do all your HTTP requests in a very simple and comprehensive way, while enabling you to do more advanced stuff by giving you access to the native API.

How do I install it?
====================

```bash
go get github.com/franela/goreq
```

What can I do with it?
======================

## Making requests with different methods
```go
res, err := goreq.Request{ Uri: "http://www.google.com" }.Do()
```

GoReq default method is GET.

```go
res, err := goreq.Request{ Method: "POST", Uri: "http://www.google.com" }.Do()
```

## Sending payloads in the Body

You can send ```string```, ```Reader``` or ```interface{}``` in the body. The first two will be sent as text. The last one will be marshalled to JSON, if possible.

```go
type Item struct {
    Id int
    Name string
}

item := Item{ Id: 1111, Name: "foobar" }

res, err := goreq.Request{ 
    Method: "POST", 
    Uri: "http://www.google.com", 
    Body: item,
}.Do()
```

## Specifiying request headers

## Setting timeouts

GoReq supports 2 kind of timeouts. A general connection timeout and a request specific one. By default the connection timeout is of 1 second. There is no default for request timeout, which means it will wait forever.

You can change the connection timeout doing:

```go
goreq.SetConnectionTimeout(100 * time.Millisecond)
```

And specify the request timeout doing:

```go
res, err := goreq.Request{ 
    Uri: "http://www.google.com",
    Timeout: 500 * time.Millisecond, 
}.Do()
```

## Using the Response and Error

GoReq will always return 2 values: a ```Response``` and an ```Error```.
If ```Error``` is not ```nil``` it means that an error happened while doing the request and you shouldn't use the ```Response``` in any way.
You can check what happened by getting the error message:

```go
fmt.Printlm(err.Error())
```
And you do different things depending on the error type by using:

```go
err.Timeout()
err.Connection()
err.Url()
```
Depending on what happened, those functions will return ```true``` or ```false```.

If you don't get an error, you can use safely the ```Response```.

```go
res.StatusCode //return the status code of the response
res.Body // gives you access to the body
res.Body.AsString() // will return the body as a string
res.Header.Get("Content-Type") // gives you access to all the response headers
```

## Receiving JSON

GoReq will help you to receive JSON.

```go
type Item struct {
    Id int
    Name string
}

var item Item

res.Body.FromJsonTo(item)
```

TODO:
-----

We do have a couple of [issues](https://github.com/franela/goreq/issues) pending we'll be addressing soon. But feel free to
contribute and send us PRs (with tests please :smile:).
