############################
# STEP 1: Build the executable
############################

FROM golang:alpine AS builder

# Install git and bzr
# They are required for fetching the dependencies
RUN apk update && apk add --no-cache git bzr

WORKDIR /src/app/
COPY . .

# Fetch dependencies

# Using go get
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/twilight

############################
# STEP 2: Run the executable
############################

FROM scratch

# Copy the executable
COPY --from=builder /go/bin/twilight /go/bin/twilight

# Copy static dependencies
COPY --from=builder /src/app/index.html .
COPY --from=builder /src/app/static ./static
COPY --from=builder /src/app/maps ./maps

EXPOSE 8080 5555

# Run the hello binary.
ENTRYPOINT ["/go/bin/twilight"]