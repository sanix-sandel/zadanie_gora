FROM golang:alpine
RUN apk add build-base
WORKDIR /app 
COPY . /app
RUN go build -o images_upload
ENV PORT 8000
EXPOSE 8000
ENTRYPOINT [ "./images_upload" ]