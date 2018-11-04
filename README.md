# Yet another gRPC echo server

YAGES (yet another gRPC echo server) is an educational gRPC server implementation. The goal is to learn gRPC and communicate best practices around its deployment and usage in the context of Kubernetes.

- [As an Kubernetes app](#as-an-kubernetes-app)
  - [From inside the cluster](#from-inside-the-cluster)
  - [From outside the cluster](#from-outside-the-cluster)
- [As a local app](#as-a-local-app)
  - [Install](#install)
  - [Use](#use)
  - [Develop](#develop)

## As an Kubernetes app

You can install YAGES as an app in your Kubernetes cluster (tested with Kubernetes v1.9, v1.10, and v1.11) like so:

```bash
$ kubectl apply -f http://mhausenblas.info/yages/app.yaml
```

Then, in order to invoke the service you've got essentially two options: from inside the cluster or from the outside the cluster, by exposing the service.

### From inside the cluster

To access the gRPC server from inside the cluster, you can for example use the [gump](https://quay.io/repository/mhausenblas/gump) container image that has [grpcurl](https://github.com/fullstorydev/grpcurl) installed:


```bash
$ kubectl run -it --rm grpcurl --restart=Never --image=quay.io/mhausenblas/gump:0.1 -- sh 
If you don't see a command prompt, try pressing enter.
/go $ grpcurl --plaintext yages:9000 yages.Echo.Ping
{
  "text": "pong"
}
```

### From outside the cluster

TBD: Using Ingress as shown in [ingress.yaml](ingress.yaml) or an OpenShift Route object with TLS passthrough set.

## As a local app

### Install

Requires Go 1.9 or above, do:

```bash
$ go get -u github.com/mhausenblas/yages
```


### Use

You can run `go run main.go` in `$GOPATH/src/github.com/mhausenblas/yages` or if you've added `$GOPATH/bin` to your path, directly call the binary:

```bash
$ yages
2018/03/25 16:23:42 YAGES in version dev serving on 0.0.0.0:9000 is ready for gRPC clients …
```

Open up a second terminal session and using [grpcurl](https://github.com/fullstorydev/grpcurl) execute the following:

```bash
# invoke the ping method:
$ grpcurl --plaintext localhost:9000 yages.Echo.Ping
{
  "text": "pong"
}
# invoke the reverse method with parameter:
$ grpcurl --plaintext -d '{ "text" : "some fun here" }' localhost:9000 yages.Echo.Reverse
{
  "text": "ereh nuf emos"
}
# invoke the reverse method with parameter from JSON file:
$ cat echo.json | grpcurl --plaintext -d @ localhost:9000 yages.Echo.Reverse
{
  "text": "ohce"
}
```

Note that you can execute `grpcurl --plaintext localhost:9000 list` and `grpcurl --plaintext localhost:9000 describe` to get further details on the available services and their respective methods.

### Develop

First you want to generate the stubs based on the protobuf schema. Note that this requires the Go gRPC runtime and plug-in installed on your machine, including `protoc` in v3 set up, see [grpc.io](https://grpc.io/blog/installation) for the steps.

Do the following:

```bash
$ protoc \
  --proto_path=$GOPATH/src/github.com/mhausenblas/yages \
  --go_out=plugins=grpc:yages \
  yages-schema.proto
```

Executing above command results in the auto-generated file `yages/yages-schema.pb.go`. **Do not** manually edit this file, or put in other words: if you add a new message or service to the schema defined in `yages-schema.proto` just run above `protoc` command again and you'll get an updated version of `yages-schema.pb.go` in the `yages/` directory as a result.
