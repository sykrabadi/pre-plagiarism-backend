FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go build -o pre-plagiarism

EXPOSE 8000

CMD ./pre-plagiarism