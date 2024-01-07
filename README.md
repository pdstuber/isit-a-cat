# Is it a cat?

## Overview

isit-a-cat is an application to predict whether an uploaded picture shows a cat or not.
 
The core components are a tensorflow machine learning model and a web application to upload pictures and trigger prediction.

There also is a telegram bot.

## Services

- [__frontend__](./frontend): Vue.js application where you can upload a picture and see if it's a cat or not
- [__backend__](./internal/api): golang service that handles RESTful HTTP requests from the frontend
- [__bot__](./internal/bot): golang telegram bot as an alternative to the web application
- [__learn__](./learn): the python machine learning code

## Other Components

- [min.io](https://min.io): Kubernetes-native object storage. Used for storing the machine learning model and uploaded images
- [Keras](https://keras.io/): Python deep learning API. Used with tensorflow backend to train a deep learning model

## Deployment

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

