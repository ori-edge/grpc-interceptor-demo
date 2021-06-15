### Prerequisites
- Go >=1.15
- [protoc](https://grpc.io/docs/protoc-installation/) 
- Optional: Read [our guide](https://github.com/ori-edge/grpc-interceptor-demo/blob/master/docs/blog.md) with a good introduction to gRPC interceptors

### Installation
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
