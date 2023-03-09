# Pre-Plagiarism Backend

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
### Environment Variables
Before you test this application, please supply these environment variables on your `.env` that you should locate on the root of this repo. Table below tells the key and the value of the environment variable
| Key  | Value  |
|---|---|
| MONGODB_DB_NAME   | documents  |
| MINIO_ENDPOINT   | localhost:9000  |
| MINIO_ACCESS_KEY_ID   | Q3AM3UQ867SPQQA43P2F  |
| MINIO_SECRET_ACCESS_KEY   | zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG  |
| MINIO_BUCKET   | documents  |
| RABBITMQ_URL_ADDRESS   | amqp://guest:guest@localhost:5672/  |
| KAFKA_BROKER_ADDR   | localhost:9092  |
### NSQ
1. If you encountered error `nsqadmin: UPSTREAM_ERROR: Failed to query any nsqd` on nsqadmin, change the value of `--broadcast-address` from current value to value as set on `container_name` on `docker-compose.yml` (ref: https://github.com/nsqio/nsq/issues/1040).
2. Make sure to add `nsqd` as value to `127.0.0.1` in your `etc\host` file
### Graphite Integration
In case you want to unlock the graphite integration with NSQ on this project using the non legacy namespace (in this project, the version of graphite using legacy namespace by default), please follow these steps
1. Go to CLI on `graphite` container, then move to `opt/graphite/conf`
2. Update the `storage-aggregation.conf` and `storage-schemas.conf` file as mentioned here [https://nsq.io/components/nsqd.html#statsd--graphite-integration](https://nsq.io/components/nsqd.html#statsd--graphite-integration)
3. Add the `graphite` key at `udp.js` that located at `opt/statsd/config`. The result would be like this
```json
{
  "graphiteHost": "127.0.0.1",
  "graphitePort": 2003,
  "port": 8125,
  "flushInterval": 10000,
  "graphite":{"legacyNamespace": false},
  "servers": [
    { server: "./servers/udp", address: "0.0.0.0", port: 8125 }
  ]
}
```
4. Access the graphite browser. Now you'll see the `counter` folder beneath the `stats` folder

## How To Run
1. Make sure you are already install `Go`. In this version, we use `v1.19.1` of `Go`.
2. Clone this repo, to desired location, then run `go mod tidy`
3. Make sure you are already install `Docker`. After you cloned this repo, run `docker compose up`
4. Run the app by typing `go run main.go`

## Note
1. Make sure to the env variable to use `MINIO_ACCESS_KEY_ID` and `MINIO_SECRET_ACCESS_KEY` as username and password when you want to login to Minio console

