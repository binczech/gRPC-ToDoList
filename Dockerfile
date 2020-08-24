FROM golang:1.13.8 as build-EnforceVersion
WORKDIR /build

COPY go.mod .
COPY go.sum .
COPY . .
COPY /server/binczech-test-a273644ddbb5.json .

RUN go build -o todoserver ./server

ENV GOOGLE_APPLICATION_CREDENTIALS=binczech-test-a273644ddbb5.json

CMD ["./todoserver"]