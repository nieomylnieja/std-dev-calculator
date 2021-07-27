FROM golang:1.16-alpine

ENV APP_NAME=std-deviation-calculator

COPY . /tmp/$APP_NAME

RUN cd /tmp/$APP_NAME && go build -o /usr/local/bin/$APP_NAME .

ENTRYPOINT ["std-deviation-calculator"]
