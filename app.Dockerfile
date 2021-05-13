# Goをビルド
FROM golang:1.16 AS go-builder
WORKDIR /workspace/server
COPY ./go.mod /workspace/go.mod
COPY ./go.sum /workspace/go.sum
COPY ./server /workspace/server
RUN echo $DEV && echo $POSTGRES_DB
RUN go build server.go

# Next.jsをビルド
FROM node:14 AS node-builder
WORKDIR /workspace/client
COPY ./client /workspace/client
RUN yarn && yarn export

# Ubuntuに各ビルド結果をコピーする
FROM ubuntu
COPY --from=go-builder /workspace/server/server /workspace/server/server
COPY --from=node-builder /workspace/client/out /workspace/client/out
RUN chmod 755 /workspace/server/server

# サーバを起動
CMD ["/workspace/server/server"]
