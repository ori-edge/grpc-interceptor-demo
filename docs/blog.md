- [Intro](#intro)
- [What are gRPC Interceptors?](#what-are-grpc-interceptors)
- [Types of Interceptor](#types-of-interceptor)
  - [Client Unary Interceptor](#client-unary-interceptor)
  - [Server Unary Interceptor](#server-unary-interceptor)
  - [Client Stream Interceptor](#client-stream-interceptor)
  - [Server Stream Interceptor](#server-stream-interceptor)
- [Real World Example - Edge Locations](#real-world-example---edge-locations)
- [Lets Build!](#lets-build)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Register](#register)
    - [List](#list)
  - [Building our Interceptors](#building-our-interceptors)
    - [Unary Client Interceptor](#unary-client-interceptor)
    - [Unary Server Interceptor](#unary-server-interceptor)
    - [Stream Client Interceptor](#stream-client-interceptor)
    - [Stream Server Interceptor](#stream-server-interceptor)
- [Conclusion](#conclusion)

## Intro 

This blog assumes a general understanding and competency when working with protocol buffers already - but if you're new to this concept, there are several fantastic articles on the web that can serve as a preface to this one, such as the Go basics tutorial from the creators of gRPC, [here](https://grpc.io/docs/languages/go/basics/). 

In this blog, we're going to cover what gRPC interceptors are, alongside an open-source example that you can make your own changes to, and see how the different components come together to form more complex systems. 

## What are gRPC Interceptors?

Interceptors are neat little components of a gRPC application that allow us to interact with a proto message or context either before - or after - it is sent or received by the client or server. For example, we may want to modify each request before it is sent, perhaps by adding some information or metadata about the environment the client is running in - and we can do that! On the server side - we can also intercept that message before the actual function call is executed - perhaps running some validation or a check before our main business logic is run. 

There are limitless reasons why this is useful as you can probably already infer, but to list a few common use cases:
- Tracing (i.e. viewing the flow of data through an application from request to response, time spent processing on the server and more. [Read our tracing blog here!](https://edgehog.blog/tutorial-how-to-implement-jaeger-and-opentracing-as-tracing-middleware-e3e693ee0802)
- Authorisation (making sure that an authenticated user is actually allowed to perform the action)
- Adding metadata to a request such as environment configuration of the client

## Types of Interceptor

In order to not confuse things - there are four different kinds of interceptor - 2 unary interceptors, and 2 stream interceptors. One for the client and server!

### Client Unary Interceptor
- These are executed on the client before a singular request is made to a function on the server side. This is where we might enrich a message with client side metadata, such as some information about the hardware or OS the client is running on, or potentially start our tracing flow.

### Server Unary Interceptor
- These are executed server side when a singular request is _received_ from a client. It's at this point we may want to perform some checks on the authenticity of the request, such as authorising it, or checking that certain fields are present / validating the request.

### Client Stream Interceptor
- These serve the same purpose as the unary interceptor, only for when we are _streaming_ data from a client to a server. For example, if we were streaming a list of 100 objects to the server, such as chunks of a file or video, we could intercept before sending each chunk, and validating things like the checksum are valid, adding metadata to a frame and the likes.

### Server Stream Interceptor
- These serve the same purpose as the server interceptor, only for when we are receiving streamed data from a client. For example, if we're receiving the aforementioned chunks of a file, maybe we want to establish that nothing has been lost in transit and verify the checksum again before storing.

![server-client-interaction](https://i.imgur.com/e9FlK0e.png)

## Real World Example - Edge Locations
Here at Ori, we specialise in making the future of cloud computing, and right now - that's edge. Lets say, hypothetically, that we are running a series of microservices in our cloud. All of these services allow a customer, whether that be a local supermarket, warehouse, or user with some spare compute in his basement, to register with Ori using a corresponding client on their hardware. Once registered, we take some of their information and evaluate whether we can on-board our infrastructure onto their compute, and allow public end users to to run their containerised workloads on these nodes.

Users that wish to enroll, must install a client that allows a user to register their hardware. However - we also want to grab some metadata about the hardware that the layman user may not know, such as the operating system the client is running on, and the IP address the client is connecting from.

This is where our interceptors come in!

In this case, we're going to grab the operating system of the *client* and get the IP adress on the *server* side.

## Lets Build!

### Prerequisites
- Go >=1.15
- [protoc](https://grpc.io/docs/protoc-installation/) 

### Installation
You can download the source code for both the client and server as they exist from this repository: https://github.com/ori-edge/grpc-interceptor-demo

Open a terminal and simply run `make protoc && make build` to ensure that the generated files are up to date, and see the client and server binaries located in your bin directory within the repo.

Now, open up 2 terminals. One terminal is going to be used to run your server, whilst the other is going to be where you are tinkering with your client. To run the server, navigate to your terminal and `cd` to the repo directory. Run the server with `./bin/edge-server`. You should see the following in your terminal, indicating that your server is running correctly:

```
2021/05/28 16:17:20 starting server...
```

The server currently exposes three bits of functionality to an end-user:
- Register
  - This allows a client to register with Ori's platform, indicating they want to provide edge compute to end users.
- List
  - List a number of registered clients that are a part of Ori's network. Provide a comma seperated list of ids you wish to retrieve with the `--list` flag. 

### Usage
#### Register
To register a client:
```bash
./bin/edge-client register
```
Should create the following output in the server:
```
2021/05/28 16:25:24 registering client success
```

#### List
To list clients that are registered with Ori's network:
```bash
./bin/edge-client list --list=c45d87da-95c8-4fe7-8403-cb798e9d2805,dc441b38-851a-4c28-b82c-f1adfa271f10
```
Should create the following output in the server:
```
2021/05/28 16:28:10 streaming edge locations...
```
And return those locations to a user:
```
2021/05/28 16:28:10 id:"c45d87da-95c8-4fe7-8403-cb798e9d2805" updated_at:{seconds:1622215524 nanos:854851258}
2021/05/28 16:28:10 id:"dc441b38-851a-4c28-b82c-f1adfa271f10" updated_at:{seconds:1622215524 nanos:854851258}
```
### Building our Interceptors
Now that we have our local environment up and running, open up a code editor of your choosing - let's build those interceptors!

The first thing we are going to do is to create a new package within our `pkg/` directory, lets call this directory `interceptor` and add an `interceptor.go` file here, this is where all of the business logic for our interceptors will live.

#### Unary Client Interceptor
We are then going to define a function that will handle our unary client side interceptor. For that, lets create a function that has the following method signature:

```go
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
}
```

As you can see, when this is called, we return the gRPC internal type of unary client interceptor. 

Let's flesh this out! 

Within this function we are going to add the following functionality:

```go
return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
    // Get the operating system the client is running on
    cos := runtime.GOOS

    // Append the OS info to the outgoing request
    ctx = metadata.AppendToOutgoingContext(ctx, "client-os", cos)

    // Invoke the original method call
    err := invoker(ctx, method, req, reply, cc, opts...)

    log.Printf("client interceptor hit: appending OS: '%v' to metadata", cos)

    return err
}
```
Breaking this down:

- The `return func()`  part of this function adheres to the method signature of a unary client interceptor
- We grab the operating system that the client is running on
- We append a key-value pair to the outgoing context, i.e. the one that is sent to the server, in the format:
    - Key: "client-os"
    - Value: Client operating system in the variable `cos`
- The `invoker` function calls the intended function call before we intercepted it, with the now modified context
- Print some debug and return any errors!

We have the functionality of our interceptor now, but how do we make it _work_? This is really simple! When we dial our gRPC server in the client:

```go
conn, err := grpc.Dial("localhost:5565", grpc.WithInsecure())
```

We just need to add an additional option:
```go
conn, err := grpc.Dial(
  "localhost:5565",
  grpc.WithInsecure(),
  grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor()),
)
```

To test that it works, run your client and server as detailed above, and see if you see the debug log of the interceptor being hit!

#### Unary Server Interceptor
We are then going to define a function that will handle our unary server side interceptor. For that, lets create a function that has the following method signature within `interceptor.go`:

```go
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
}
```

To flesh out this function call, add the following body to the above function:

```go
return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
  // Get the metadata from the incoming context
  md, ok := metadata.FromIncomingContext(ctx)
  if !ok {
    return nil, fmt.Errorf("couldn't parse incoming context metadata")
  }

  // Retrieve the client OS, this will be empty if it does not exist
  os := md.Get("client-os")
  // Get the client IP Address
  ip, err := getClientIP(ctx)
  if err != nil {
    return nil, err
  }

  // Populate the EdgeLocation type with the IP and OS
  req.(*api.EdgeLocation).IpAddress = ip
  req.(*api.EdgeLocation).OperatingSystem = os[0]

  h, err := handler(ctx, req)
  log.Printf("server interceptor hit: hydrating type with OS: '%v' and IP: '%v'", os[0], ip)

  return h, err
}
```

Then, add an additional function to get the client IP address from the incoming context:

```go
// GetClientIP inspects the context to retrieve the ip address of the client
func getClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("couldn't parse client IP address")
	}

	return p.Addr.String(), nil
}
```

Again, let's break down what this is doing:
- As we saw in our client interceptor, we add additional metadata to the outgoing context. Here, we pull that metadata out with the function `metadata.FromIncomingContext(ctx)`. This grants us access to the client-os key value pair we created earlier
- We then grab the IP address by inspecting the metadata of the incoming context using our `getClientIP` function
- When we have both pieces of additional metadata from the client request, we can modify the incoming message type to include these two pieces of information
- We then continue to the intended server side function call through the use of the `handler`, passing any modified contexts or requests on to that function call

In order to register this with the server, add the following line when creating your new gRPC server instance:
```go
s := grpc.NewServer()
```
Should become:
```go
s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()))
```

#### Stream Client Interceptor
This functionality intercepts every streamed object from the client to the server - in this case - listing edge locations.

Create a new function in `interceptor.go`:

```go
// StreamClientInterceptor allows us to log on each client stream opening
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Printf("opening client streaming to the server method: %v", method)

		return streamer(ctx, desc, cc, method)
	}
}
```

As you can see from previous implementations, we follow the exact same workflow. We don't actually perform any business logic here, bar logging the gRPC method that is being called. However, this is the perfect use case for initialising tracing or the likes, and appending it to the outgoing context.

In order to register this stream client interceptor with our client, add to the existing dial options the following:

```go
conn, err := grpc.Dial(
    "localhost:5565",
    grpc.WithInsecure(),
    grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor())
)
```

Should become:
```go
conn, err := grpc.Dial(
    "localhost:5565",
    grpc.WithInsecure(),
    grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor()),
    grpc.WithStreamInterceptor(interceptor.StreamClientInterceptor()),
)
```

#### Stream Server Interceptor
The way a stream server interceptor works is a little more difficult to wrap your head around - we have to embed the `grpc.ServerStream` type in a struct, in order to get access to the `RecvMsg` function. For example, to create our interceptor as we have seen before:

```go
// Set up a wrapper to allow us to access the RecvMsg function
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &EdgeServerStream{
			ServerStream: ss,
		}
		return handler(srv, wrapper)
	}
}
```

This looks fairly familiar, though we are using a wrapper. What does that allow us to do? Lets investigate further.

```go
// Embedded EdgeServerStream to allow us to access the RecvMsg function on
// intercept
type EdgeServerStream struct {
	grpc.ServerStream
}

// RecvMsg receives messages from a stream
func (e *EdgeServerStream) RecvMsg(m interface{}) error {
	// Here we can perform additional logic on the received message, such as
	// validation
	log.Printf("intercepted server stream message, type: %s", reflect.TypeOf(m).String())
	if err := e.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}
```

By wrapping the sever stream in our own type, we allow ourselves access to the `RecvMsg` function, which as above displays, allows us to interact with _each_ message received at a time! Pretty useful, huh?

In order to register this with our server, change the following code to implement the stream server interceptor you just created:
```go
s := grpc.NewServer(
    grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()),
)
```

```go
s := grpc.NewServer(
    grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()),
    grpc.StreamInterceptor(interceptor.StreamServerInterceptor()),
)
```

## Conclusion

The above has hopefully given you an example of each of the interceptors that we have available to us when working with gRPC applications with a real world application and use case that you can interact with. Remember that all of the code is available in our public GitHub repository here: https://github.com/ori-edge/grpc-interceptor-demo/. There is an open pull request on this repo where you can see the changes that are required to implement the above interceptors easily: https://github.com/ori-edge/grpc-interceptor-demo/pull/2.

Thanks for reading!