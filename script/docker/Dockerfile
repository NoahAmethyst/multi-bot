FROM ubuntu:18.04
RUN apt-get -y update && apt-get -y install ca-certificates
COPY . /opt
WORKDIR /opt

CMD ./cloud-wallet-exe
