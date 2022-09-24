# go-nsq

## How To Run
1. Make sure you are already install `Go`. In this version, we use `v1.19.1` of `Go`.
2. Clone this repo, to desired location, then run `go mod tidy`
3. After the required packages are successfully installed, then run the app by typing `go run main.go`
4. Make sure you are already install `Docker`. After you cloned this repo, run `docker compose up`

## Note
Make sure to the env variable to use `MINIO_ACCESS_KEY_ID` and `MINIO_SECRET_ACCESS_KEY` as username and password when you want to login to Minio console
