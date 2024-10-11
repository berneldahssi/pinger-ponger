# Starling infrastructure recruitment challenge

## Overview

Two dummy micro services based on HTTP Rest and written in Go are interacting with each other. The Ponger service has one endpoint `GET /ping` which responds with body `pong` and a 200 success message. The Pinger service consumes the service provided by Ponger by sending requests to it in regular intervals. We provide a Makefile to run the services on a k3d local kubernetes cluster. 

As you complete this test, please create your own README with instructions on how to run your solution and chnages you've made to get things working cleanly. 

## Time
We haven’t set a time limit for this test. We want you to spend the time you need to hand something in that you’re happy with - at the same time we don’t want you to spend hours on end either. As a guideline, some of our engineers here take around an hour or 2 on average to get the job done cleanly.

## Submitting
Share a private github repo with the user starlingtechtest - ideally we’d like to avoid lots of public solutions. Please include your name in the repository description, if it is not obvious from your github profile, so we know whose code we are looking at!

Prerequisites:
* Having a container runtime installed on your machine (https://docs.docker.com/install/)
* Install k3d  

## Challenges:

1) Create docker images locally called "pinger" and "ponger" for the two services. Their source codes are respectively in `app/pinger` and `app/pinger`. The Dockerfiles should be stored in these directories. The docker images can be tagged as latest. Here is a command example to build the binary of the pinger service :
```
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a --installsuffix cgo --ldflags="-s -w" -o /pinger
```

2) The Makefile found along with this README may require changes depending on your project structure so please review before you carry on. It's advised to make sure the intermediary targets work to build individual images before you carry on. 

3) The ponger service can be provided with a certificate path in its configuration file:

```yaml
---
  service:
    protocol: https
    tlsCertificate: /path/to/cert
    tlsPrivateKey: /path/to/private_key
```

For the purposes of this assignment, we do not require a CA-signed certificate, a self-signed certificate that you generate on your machine will do. However, the pinger service needs to be instructed to accept the self-signed certificate. The pinger service can be provided with a certificate to accept as valid:
```yaml
---
  ponger:
    url: https://localhost:8080
    acceptCert: /path/to/cert
```



4) Run a local kubernetes cluster with the below make target. The kubernetes context of this cluster is `k3d-cluster`
```
make run-local-kube-with-ping-pong-app
```

5) Are the services up and running? If not, please fix the issues and provide details about them in your README. 

6) What would you do differently or add to run these micro services in a production environment ?

7) Bonus : feel free to implement some suggestions made in the previous question, if any.


