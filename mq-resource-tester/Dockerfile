FROM golang:1.5

#install tools
RUN go get -v github.com/FiloSottile/gvt github.com/onsi/ginkgo/ginkgo
RUN echo "" > /tmp/.keep

# Install Supervisor
RUN apt-get update && apt-get -qy install supervisor
ADD supervisord-app.conf /etc/supervisor/conf.d/supervisord-app.conf
ADD run.sh /run.sh
RUN echo "#!/bin/sh \n" > /start.sh
RUN chmod +x /*.sh

ENV MQ_RESOURCE_DEST $GOPATH/src/github.com/aminjam/mq-resource
ENV GO15VENDOREXPERIMENT 1
WORKDIR $MQ_RESOURCE_DEST
