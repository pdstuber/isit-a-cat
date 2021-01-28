# Is it a cat?

## Overview

isit-a-cat is an application to predict whether an uploaded picture shows a cat or not.
 
The core components are a tensorflow machine learning model and a web application to upload pictures and trigger prediction.

## Services

- [__frontend__](./frontend): Vue.js application where you can upload a picture and see if it's a cat or not
- [__bff__](./bff): golang backend for frontend service that handles RESTful HTTP requests from the frontend and delegates work to the prediction service via NATS messaging
- [__predict__](./predict): Uses the golang tensorflow bindings to load the tensorflow model and make a prediction on a given picture
- [__learn__](./learn): the python machine learning code

## Other Components

- [NATS](https://nats.io): Cloud native messaging. Used for communication between microservices
- [min.io](https://min.io): Kubernetes-native object storage. Used for storing the machine learning model and uploaded images
- [Prometheus](https://prometheus.io/): Used for metrics via the golang client library
- [MessagePack](https://msgpack.org/): Used for Serialization / Deserialization of binary messages
- [Keras](https://keras.io/): Python deep learning API. Used with tensorflow backend to train a deep learning model

## Deployment

### Kubernetes

__Prerequisites__
- kubectl
- a valid kubeconfig pointing to your cluster
- helm v3 installed in your cluster

The folder __kubernetes__ includes the following scripts 
- __apply_all.sh__
- __delete_all.sh__

The following components are installed on the cluster with __helm__ and  __kubectl__:
- Prometheus operator
- Nats operator
- cert-manager operator and pods
- the isit-a-cat namespace with the microservices, nats broker, minio object storage service and cert-manager resources

__The ingress configuration for the frontend and the cert-manager issuer needs to be adapted!__

### Docker-Compose

__Prerequisites__

- Docker
- docker-compose

Deploy:

Clone this repo to a folder of your choice. 

Rename the __env_template__ to __.env__ and set the two environment variables within to random values.

E.g. with:
```bash
openssl rand --base64 10
openssl rand --base64 40
```

Then build and run everything with:

```bash
docker-compose up --build -d --remove-orphans
```

