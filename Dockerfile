FROM golang:1.16-alpine

ENV APP_NAME=std-dev-calculator

COPY . /tmp/$APP_NAME

RUN cd /tmp/$APP_NAME && go build -o /usr/local/bin/$APP_NAME .

ENV RANDOM_INTEGERS_GENERATOR_URL=http://localhost:8081/integers

ENTRYPOINT ["std-dev-calculator"]
