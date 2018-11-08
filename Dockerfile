FROM golang:1.11-stretch

RUN  apt-get update && apt-get install -y apt-transport-https && \
  curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
  echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list && \
  apt-get update && \
  apt-get install -y kubectl=1.11.3-00 jq && \
  apt-get clean

ADD . /src
WORKDIR /src
RUN go build  -o workshop-namespaces main.go

EXPOSE 9090

ENTRYPOINT ["./workshop-namespaces"]
