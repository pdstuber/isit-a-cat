# isit-a-cat-bff

This service serves as a backend for the ui part of the isit-a-cat application.

It uses [NATS](https://nats.io) messaging for triggering downstream services and [google cloud storage](https://cloud.google.com/storage) for storing/retrieving data.


## HTTP Endpoints

| HTTP Method | Endpoint          | Description                                                                                                                        |
|-------------|-------------------|------------------------------------------------------------------------------------------------------------------------------------|
| POST        | /images           | Receive an uploaded image, save in object storage, trigger prediction.  Returns an unique ID for retrieving the prediction result. |
| GET         | /predictions/{id} | Fetch prediction results for a given id                                                                                            |
| GET         | /images{id}       | Get the uploaded image with the given ID                                                                                           |
| GET         | /ping             | Health check                                                                                                                       |
| GET         | /metrics          | Endpoint for prometheus metrics                                                                                                    |

## Configuration

The following environment variables are relevant for the application:

| Variable Name           | Mandatory | Default value    | Description                                                        |
|-------------------------|-----------|------------------|--------------------------------------------------------------------|
| OBJECT_STORAGE_ENDPOINT | yes       | -                | The endpoint of the min.io server                                  |
| MINIO_ACCESS_KEY        | yes       | -                | The access key for the min.io server                               |
| MINIO_SECRET_KEY        | yes       | -                | The access key secret for the min.io server                        |
| CORS_ALLOWED_ORIGIN     | no        | *                | The allowed origin(s) for CORS                                     |
| OBJECT_STORAGE_USE_TLS  | no        | false            | This toggles whether TLS is used for communication with min.io     |
| SERVICE_HOST_PORT       | no        | 0.0.0.0:8080     | The host and port the service should bind to                       |
| STORAGE_BUCKET_NAME     | no        | isit-a-cat       | The google cloud storage bucket name to use                        |
| STORAGE_OBJECT_FOLDER   | no        | uploaded-images/ | The google cloud folder to upload images to                        |
| MESSAGING_SUBJECT_NAME  | no        | predictions      | The NATS topic name for predictions                                |
| MESSAGING_BROKER_URL    | no        | 0.0.0.0:4222     | The url of the NATS broker                                         |                     |