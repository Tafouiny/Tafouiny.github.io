FROM golang:1.20-alpine3.18
LABEL MAINTAINER: "niangpmaxgatte@gmail.com, ivan,ndiayetapha7@gmail.com, oulam, djibsow " 
RUN mkdir /app
RUN apk update && apk add bash && apk add tree
RUN apk add --no-cache gcc musl-dev
COPY . /app
WORKDIR /app
RUN go build -o forum
CMD ["./forum"]