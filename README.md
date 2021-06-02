# openapigen
Go library and CLI for generating an OpenAPI 3.0 spec from a RESTful endpoint (following the https://aip.dev best practices)

## Usage

### Go module usage
To use in your Go program or library, simply get and import the openapigenlib module like this:

```
go get github.com/tyayers/openapigen/openapigenlib
```

Then you can generate OpenAPI specs in your code by calling:

```
spec := openapigenlib.GenerateSpec("https://jsonplaceholder.typicode.com/todos")
```

The spec variable above will contain a string of the complete OpenAPI spec that fits the payload returned by the given URL.

### CLI usage
Check out a simple openapigen CLI in the **cmd** directory, which can be installed like this:

```
go get github.com/tyayers/openapigen/cmd/openapigen
```

And then you can call the CLI like this and direct the output into a local file.

```
./openapigen https://jsonplaceholder.typicode.com/todos > todo.yaml
```

## Projects
The main project that uses this library is a simple web frontend for generating OpenAPI specs here: [https://github.com/tyayers/openapigen-client](https://github.com/tyayers/openapigen-client)
