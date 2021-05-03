# Hyperledger

Intended as a foundation for developing applications or solutions with a modular architecture, Hyperledger Fabric allows components, such as consensus and membership services, to be plug-and-play.

## Fabric-sdk-rest
	$ cd fabric-sdk-rest/
	$ ./setup.sh -f /root/go/src/github.com/hyperledger/fabric-samples-cgu/basic-network/ -ukat
	$ nohup node . > /dev/null 2>&1 &

## Blockchain-explorer
	$ mysql -u root -p < db/fabricexplorer.sql 

# Socket
## Paper name
	### Hybrid Blockchain-Based Log Management Scheme with Accountability for Smart Girds

## CA
	$openssl req -new -x509 -days 365 -nodes -out cert.pem -keyout key.pem

## Environment

* Hyperledger 1.0.0 -> 1.1.0
* Golang Version 1.9
* Docker Version 18.03.1-ce
* Docker-Compose Version 1.21.1
* Kafka 4.4
* Redis 4.0.9

# Package Installation
### Ubuntu
	$ sudo apt-get install libcurl4-openssl-dev libssl-dev
	
### Mac 
	$ brew install openssl
	$ brew install curl
	$ export PYCURL_SSL_LIBRARY=openssl

## Docker
### Redis
	$ docker pull redis:4.0.9
	$ docker run -p 6379:6379 -d redis:4.0.9 redis-server

	# Using redis password

	$ docker run -p 6379:6379 -d redis:4.0.9 /
	redis-server --requirepass dockeredis
### Kafka
    $ docker-compose up -d

## Golang Build and Pack
Golang開發的程序都會比較大，這是因為Golang是靜態編譯的，編譯打包之後基本就不會再對其他類庫有依賴了，所以會比較大。
舉個例子：C++程序可以調用dll，所以打包的時候可以不把dll打進去，包自然就小了。
之前還有看到過有人使用GO -> C -- dll --> C -> GO的方式間接實現了Golang的偽動態鏈接，有興趣的同學可以研究一下。

	$ go build -ldflags '-w -s'
	$ upx --brute -o Server

* Author:gwpp
* Link:<a herf="https://www.jianshu.com/p/cd3c766b893c">Source</a>

# Web

> Golang Gin
>
> Template : https://github.com/modularcode/modular-admin-html

kubectl describe pod blockchain-orderer | grep IP | awk '{print $2}' |head -1
kubectl describe pod blockchain-org1peer1 | grep IP | awk '{print $2}' |head -1
kubectl describe pod blockchain-org2peer1 | grep IP | awk '{print $2}' |head -1
kubectl describe pod blockchain-org3peer1 | grep IP | awk '{print $2}' |head -1
kubectl describe pod blockchain-org4peer1 | grep IP | awk '{print $2}' |head -1
kubectl describe pod blockchain-ca | grep IP | awk '{print $2}' |head -1

172.17.0.8  orderer
172.17.0.10 peer0.org1.example.com
172.17.0.7  peer0.org2.example.com
172.17.0.9  peer0.org3.example.com
172.17.0.11 peer0.org4.example.com
172.17.0.6 ca.example.com