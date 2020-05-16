FROM golang:1.14
# create a working directory
WORKDIR /root/data/
# add source code
COPY . /root/data/
RUN go build -mod=vendor

RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -  && \
    apt-get install -y nodejs

WORKDIR /root/data/react-api/
RUN npm install
RUN npm run build 

EXPOSE 5432 80
# run main.go
WORKDIR /root/data/
CMD ["go","run","-mod=vendor","main.go"]