FROM node:14-alpine as webui

WORKDIR /app
COPY ./noid-webui/package.* .
RUN npm install

COPY ./noid-webui .
RUN npm run build

FROM golang:1.15-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
COPY noid ./noid
RUN CGO_ENABLED=0 go build -o main .

FROM scratch

COPY --from=builder /app/main /noid
COPY --from=webui /app/dist /public
ENTRYPOINT [ "/noid" ]