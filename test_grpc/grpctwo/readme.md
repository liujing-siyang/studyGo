# 生成私钥文件
openssl genrsa -des -out ca.key 2048

输入密码如：123456

# 创建证书请求
openssl req -new -key ca.key -out ca.csr

# 生成ca.crt
openssl x509 -req -days 365 -in ca.csr -signkey ca.key -out ca.crt

# 服务端
## 生成证书私钥文件
openssl genpkey -algorithm RSA -out server.key

## 通过私钥server.key 生成证书请求文件server.crs
openssl req -new -nodes -key server.key -out server.csr  -days 365 -config ./openssl.cnf -extensions v3_req

## 生成SAN证书
openssl x509 -req -days 365 -in server.csr -out server.pem -CA ca.crt -CAkey ca.key -CAcreateserial -extfile ./openssl.cnf -extensions  v3_req

# 客户端
## 生成证书私钥文件
openssl genpkey -algorithm RSA -out client.key

## 通过私钥server.key 生成证书请求文件server.crs
openssl req -new -nodes -key client.key -out client.csr  -days 365 -config ./openssl.cnf -extensions v3_req

## 生成SAN证书
openssl x509 -req -days 365 -in client.csr -out client.pem -CA ca.crt -CAkey ca.key -CAcreateserial -extfile ./openssl.cnf -extensions  v3_req