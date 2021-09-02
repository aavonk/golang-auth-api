FROM golang:1.16

# Create a working directory
WORKDIR /app

# fetch dependencies on a separate layer as they are less likely to
# change on every build and will be cached, speeding up each build
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy the source from the host to the working directory
COPY . . 

# expoxe port 7777 to the oustide world
EXPOSE 7777

