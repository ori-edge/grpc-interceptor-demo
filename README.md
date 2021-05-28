## Intro 

This blog assumes a general understanding and competency when working with protocol buffers already - but if you're new to this concept, there are several fantastic articles on the web that can serve as a preface to this one. In this blog, we're going to cover what gRPC interceptors are, alongside an open-source example that you can make your own changes to, and see how the different components come together to form more complex systems. 

## What are gRPC Interceptors?

Interceptors are neat little components of a gRPC application that allow us to interact with a proto message either before - or after it is sent or received by the client or server. For example, we may want to modify each request before it is sent, perhaps by adding some information or metadata about the environment the client is running in - and we can do that! On the server side - we can also intercept that message before the actual function call is executed - perhaps running some validation or a check before our main business logic is ran. 

There are limitless reasons why this is useful as you can probably already infer, but to list a few common use cases:
- Tracing (i.e. viewing the flow of data through an application from request to response, time spent processing on the server and more. ([TODO: add niki's article here])
- Authorisation (making sure that an authenticated user is actually allowed to perform the action)
- Adding metadata to a request such as environment configuration of the client

## Types of Interceptor

In order to not confuse things - there are four different kinds of interceptor - 2 unary interceptors, and 2 stream interceptors. One for the client and server!

### Client Unary Interceptor
- These are executed on the client before a singular request is made to a function on the server side. This is where we might enrich a message with client side metadata, such as some information about the hardware or OS the client is running on, or potentially start our tracing flow.
- 
### Server Unary Interceptor
- These are executed server side when a singular request is _received_ from a client. It's at this point we may want to perform some checks on the authenticity of the request, such as authorising it, or checking that certain fields are present / validating the request.

### Client Stream Interceptor
- These serve the same purpose as the unary interceptor, only for when we are _streaming_ data from a client to a server. For example, if we were streaming a list of 100 objects to the server, such as chunks of a file or video, we could intercept before sending each chunk, and validating things like the checksum are valid, adding metadata to a frame and the likes.

### Server Stream Interceptor
- These serve the same purpose as the server interceptor, only for when we are receiving streamed data from a client. For example, if we're receiving the aforementioned chunks of a file, maybe we want to establish that nothing has been lost in transit and verify the checksum again before storing.

[TODO: Diagram Goes Here]

## Real World Example - Edge Locations
Here at Ori, we specialise in making the future of cloud computing, and right now - that's edge. Lets say, hypothetically, that we are running a series of microservice based servers in our cloud. All of these servers allow a customer, whether that be a local supermarket, warehouse, or user with some spare compute in his basement to register with Ori using a corresponding client on their hardware. Once registered, we take some of their information and evaluate whether we can on-board our infrastructure onto their compute, and allow public end users to to run their containerised workloads on these nodes.

Users that wish to enrol, must install a client that takes one simple piece of information, for this example, and that is a region. A region in this sense is the general locality of the compute they want to onboard, such as "London" or "Amsterdam". However - we also want to grab some metadata about the hardware that the layman user may not know, such as the operating system the client is running on, and the IP address the client is connecting from.

This is where our interceptors come in!

In this case, we're going to grab the operating system of the *client* and get the IP adress on the *server* side.

[TODO: INSERT DIAGRAM HERE]

## Lets Build!

### Prerequisites
- Go >=1.15
- protoc https://grpc.io/docs/protoc-installation/

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
- Get
  - Returns an individual client that is part of Ori's network based on a provided ID flag.
- List
  - List a number of registered clients that are a part of Ori's network. The default value is 10, but can be changed with a flag.

### Usage
#### Register
To register a client:
```bash
./bin/edge-client register --region=london
```
Should create the following output in the server:
```
2021/05/28 16:25:24 registering client success
```

If region is not set, it will default to `undefined`.
#### List
To list clients that are registered with Ori's network:
```bash
./bin/edge-client list --limit=10
```
Should create the following output in the server:
```
2021/05/28 16:28:10 streaming edge locations...
```
And return ten (or less) locations to a user:
```
2021/05/28 16:28:10 id:"160b7671-4570-4655-84ca-910e6e6aa529" region:"london" updated_at:{seconds:1622215524 nanos:854851258}
```
#### Get
To get an individual client that is registered with Ori's network:
```bash
./bin/edge-client get --id=160b7671-4570-4655-84ca-910e6e6aa529
```
Should create the following output in the server:
```
2021/05/28 16:29:55 retrieving edge location...
```
And return ten (or less) locations to a user:
```
2021/05/28 16:29:55 id:"160b7671-4570-4655-84ca-910e6e6aa529" region:"london" updated_at:{seconds:1622215524 nanos:854851258}
```
### Building our Interceptors

TODO