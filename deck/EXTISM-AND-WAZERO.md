


Since yesterday, the Go SDK of Extism is officially rewritten with the Wasm Wazero runtime. The new repository is accessible here: https://github.com/extism/go-sdk. And this is great news, because it means that we can now compile host applications without any C dependencies.

However, the new SDK is not released yet. So to be able to use it and build your projects, you must:

- clone the new repository: git clone git@github.com:extism/go-sdk.git
- add a local reference of this git repository into the go.mod file of your projects: replace github.com/extism/extism => ../go-sdk


## Build Static application

This means that you will be able to build static and more portable applications (shortcut: once compiled for the target system, you will only need the executable to run and nothing else).

For example, here is how I compile my applications:

```bash
export TAG="v0.0.0"
env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o slingshot-${TAG}-darwin-arm64
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o slingshot-${TAG}-darwin-amd64
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o slingshot-${TAG}-linux-arm64
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o slingshot-${TAG}-linux-amd64
```

And for me, my most important use case is that of containerizing my applications. I will be able to create extremely small Docker images thanks to the `scratch` image. And at a time when web service and FaaS execution platforms often run on Kubernetes, we will gain efficiency, and therefore reduce costs:

- small executable
- small wasm plugin
- small container

You will be able to deploy and re-deploy faster with lower memory consumption and disk space.

## Create a Docker image

I'm working on project of an HTTP server to serve Extism Wasm plugins (https://github.com/bots-garden/slingshot/tree/main/slingshot-http-server), 

And this is how I create a Docker image of my application:

Create a new Dockerfile:
```Dockerfile
FROM scratch

ADD slingshot-v0.0.0-linux-arm64 ./
ADD simple.wasm ./

EXPOSE 8080

CMD ["./slingshot-v0.0.0-linux-arm64", "./simple.wasm", "handle", "8080"]

```

> I could have my executable built by docker using the [multi stage technique](https://docs.docker.com/build/building/multi-stage/). But to simplify, I build my executable before and copy it in the image during the construction of this one.


Type the following commands to build the image:

```bash
IMAGE_NAME="demo-slingshot"
docker build -t ${IMAGE_NAME} . 

docker images | grep ${IMAGE_NAME}
```

You should get something like this:
```bash
demo-slingshot   latest    92e0aa64e692   Less than a second ago   8.46MB
```

My docker image size is less than 9MB!!! 🚀

And I can run my service like this:
```bash
IMAGE_NAME="demo-slingshot"
docker run \
  -p 8080:8080 \
  --rm ${IMAGE_NAME}
``````

Well, this was a short article to explain this very good news. I have updated the source codes of my series. And as now I am able to build very light Docker images, we will see in a next article how to deploy Wasm services on Kubernetes.

Have a nice day!
