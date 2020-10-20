FROM golang:1.15.3-buster

ADD . /chat
WORKDIR /chat

RUN apt-get update
RUN apt-get install apt-transport-https
RUN sh -c 'wget -qO- https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add -'
RUN sh -c 'wget -qO- https://storage.googleapis.com/download.dartlang.org/linux/debian/dart_stable.list > /etc/apt/sources.list.d/dart_stable.list'


RUN apt-get update
RUN apt-get install dart


ENV PATH=$PATH:/usr/lib/dart/bin

RUN pub global activate webdev
RUN cd static && pub upgrade && pub get 
RUN ls -al


ENV SERVER_PORT 7070
ENV SERVER_HOST 0.0.0.0
ENV DATABASE_NAME chat
ENV DATABASE_HOST rethink
ENV DATABASE_PORT 28015
ENV STATICS_PATH ../../../library-fe/build
ENV STATICS_PATH static
ENV SESSION_DB_NAME 0
#ENV SESSION_DB_PASSWORD 
ENV SESSION_DB_HOST redis
ENV SESSION_DB_PORT 6379

EXPOSE 7070

CMD ["go", "run", "cmd/chat/main.go"]