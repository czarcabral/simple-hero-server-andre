FROM golang:1.15.6
ENV GO111MODULE=on
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go build -o main .
CMD ["/app/main"]

# To Build:
#   cd backend
#   docker build -t fullstack1-backend .
# To Run:
#   docker run -it -p 8080:8081 fullstack1-backend    # note: tcp port 8081 must match main.go: http.ListenAndServe(":8081")
