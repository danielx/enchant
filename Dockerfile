FROM golang:stretch

RUN apt-get update && apt-get upgrade -y

RUN apt-get install -y libenchant-dev

RUN mkdir /app

ADD . /app/

WORKDIR /app

CMD ["go", "test", "-v", "--cover"]
