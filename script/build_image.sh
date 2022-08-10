git describe --tags
rm -rf docker && mkdir docker && cd docker

cp ../cloud-wallet-exe ./
cp -r ../db ./
chmod +x ./cloud-wallet-exe

cat > Dockerfile <<EOF
FROM ubuntu:18.04
RUN apt-get -y update && apt-get -y install ca-certificates
COPY . /opt
WORKDIR /opt

ENTRYPOINT ["./cloud-wallet-exe"]
EOF

docker build -t $1 .
