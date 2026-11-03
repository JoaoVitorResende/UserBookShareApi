FROM golang:1.26-alpine AS build
ENV GOTOOLCHAIN=auto
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/userapi ./cmd/api

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/userapi /userapi
USER nonroot
ENTRYPOINT ["/userapi"]