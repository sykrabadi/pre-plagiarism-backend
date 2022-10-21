# go-nsq

## System Architecure
![assets\system_architecture.jpg](assets/system_architecture.jpg)
This project use docker to pack every dependencies. However (currently) the Pre-Falsification Engine is unavailable, but you still able to run this software.

The `Entry Point Service` act as gate to access the prefalsification engine. The contract of the `Entry Point Service` provided below. We store the document on MinIO object storage and results from `prefalsification engine` on MongoDB (since we don't see any urgencies to use relational databases). The document information flows from `Entry Point Service` to `prefalsification engine` via message brokers.

## API Contract
This system use REST API to connect from frontend. The API contract is provided below

### **sendDocument**
| HTTP Method  | MIME type  |
|---|---|
| POST   | .pdf  |

This screenshot from postman shows how to use the `/sendDocument` endpoint properly

![assets\success_request.png](assets/success_request.png)

## Before You Run
Before you test this application, please supply these environment variables on your `.env` that you should locate on the root of this repo. Table below tells the key and the value of the environment variable
| Key  | Value  |
|---|---|
| MONGODB_DB_NAME   | documents  |
| MINIO_ENDPOINT   | localhost:9000  |
| MINIO_ACCESS_KEY_ID   | Q3AM3UQ867SPQQA43P2F  |
| MINIO_SECRET_ACCESS_KEY   | zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG  |
| MINIO_BUCKET   | documents  |
| RABBITMQ_URL_ADDRESS   | amqp://guest:guest@localhost:5672/  |

## How To Run
1. Make sure you are already install `Go`. In this version, we use `v1.19.1` of `Go`.
2. Clone this repo, to desired location, then run `go mod tidy`
3. After the required packages are successfully installed, then run the app by typing `go run main.go`
4. Make sure you are already install `Docker`. After you cloned this repo, run `docker compose up`

## Note
1. Make sure to the env variable to use `MINIO_ACCESS_KEY_ID` and `MINIO_SECRET_ACCESS_KEY` as username and password when you want to login to Minio console
2. If you encountered error `nsqadmin: UPSTREAM_ERROR: Failed to query any nsqd` on nsqadmin, change the value of `--broadcast-address` from current value to value as set on `container_name` on `docker-compose.yml` (ref: https://github.com/nsqio/nsq/issues/1040).
3. Make sure to add `nsqd` as value to `127.0.0.1` in your `etc\host` file
