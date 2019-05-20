# pull the latest golang image from DockerHub
FROM golang

# create the directory in the container in which to work
WORKDIR /app

# copy any dependencies into WORKDIR
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

# download dependencies
RUN go mod download

# copy the rest of the local project directory into WORKDIR
COPY . .

# build the Go binary in WORKDIR/bin/app
RUN go build -o /bin/app

# execute the Go binary
CMD [ "app" ]


# docker-compose up --build
# docker build -t skillitzimberg/app .
# docker run skillitzimberg/app