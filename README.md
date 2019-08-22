RailwaysBot takes the route information of Kanto and tweets the delay occurrence and the resolution to Slack.
The Bot periodically runs Lambda from CloudWatchEvents.


# Configuration

Change the following files:

* configs/railways.yaml     : target lines list
* configs/slack.yaml        : your slack app information

# Usage

Cross-compilation is required.

```
ex:Mac OS

$ CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o railways

$ zip railways.zip railways ./configs/*
```

Upload the railways.zip to Lambda.