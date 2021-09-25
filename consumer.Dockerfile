FROM golang:1.17-alpine

# Set destination for COPY
WORKDIR /root

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY cmd/consumer/ ./
COPY entity/ ./entity/
COPY repository/ ./repository/
COPY usecase/ ./usecase/
COPY util/ ./util/
COPY config/consumer.yaml ./config/consumer.yaml

# Build
RUN go build -o /consumer ./main.go

# Run
CMD [ "/consumer" ]