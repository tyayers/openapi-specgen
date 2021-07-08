# openapigenlib

This internal library does the actual generating of an OpenAPI spec based on a RESTful endpoint.

## Using

To use in your Go program or library, simply get and import the openapigenlib module like this:

```
go get github.com/tyayers/openapigen/openapigenlib
```

Then generate an OpenAPI spec in your code by calling:

```
spec := openapigenlib.GenerateSpec("https://jsonplaceholder.typicode.com/todos")
```

The spec variable above will contain a string of the complete OpenAPI spec that fits the payload returned by the given URL.

## Testing

To test the module just run 

```
go test
```

This will test the generation with some common REST endpoints.