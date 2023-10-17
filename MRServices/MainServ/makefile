all:
	go build ./

lint:
	golangci-lint run

vet:
	go vet ./*.go;
	go vet accountservice/*.go;
	go vet apiproxy/*.go;
	go vet bot/*.go;
	go vet common/*.go;
	go vet gmtool/*.go;
	go vet httpgateway/*.go;
	go vet itemservice/*.go;
	go vet loginservice/*.go;
	go vet matchingservice/*.go;
	go vet paymentservice/*.go;
	go vet socialservice/*.go;
	go vet sspb/*.go;
	go vet sspberr/*.go;
	go vet tcpgateway/*.go;
	go vet teamservice/*.go;
	go vet userservice/*.go;
	go vet util/*.go;

p:
	./mp.sh

cfg:
	cp -R ../config/ServerCSV/* ./resources/