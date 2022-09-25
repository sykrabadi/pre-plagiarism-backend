# go-nsq

## How To Run
1. Make sure you are already install `Go`. In this version, we use `v1.19.1` of `Go`.
2. Clone this repo, to desired location, then run `go mod tidy`
3. After the required packages are successfully installed, then run the app by typing `go run main.go`
4. Make sure you are already install `Docker`. After you cloned this repo, run `docker compose up`

## Note
1. Make sure to the env variable to use `MINIO_ACCESS_KEY_ID` and `MINIO_SECRET_ACCESS_KEY` as username and password when you want to login to Minio console
2. If you encountered error `nsqadmin: UPSTREAM_ERROR: Failed to query any nsqd` on nsqadmin, change the value of `--broadcast-address` from current value to value as set on `container_name` on `docker-compose.yml` (ref: https://github.com/nsqio/nsq/issues/1040)
