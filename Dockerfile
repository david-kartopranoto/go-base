FROM golang:1.17-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./
COPY entity/ ./entity/
COPY repository/ ./repository/
COPY usecase/ ./usecase/
COPY util/ ./util/
COPY rest/ ./rest/
COPY config/app.yaml ./config/app.yaml

# Build
RUN go build -o /go-base-main

# Run
CMD [ "/go-base-main" ]