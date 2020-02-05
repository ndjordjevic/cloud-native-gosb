# cloud-native-gosb

- To regenerate swagger docs run 'swag init' (https://github.com/swaggo/swag). To see docs go to http://localhost:8080/swagger/index.html
- To generate certificate and private key for TLS
    - openssl req -newkey rsa:2048 -nodes -keyout domain.key -out domain.csr -subj "/C=IN/ST=Mumbai/L=Andheri East/O=Packt/CN=packtpub.com"
    - openssl req -key domain.key -new -x509 -days 365 -out domain.crt -subj "/C=IN/ST=Mumbai/L=Andheri East/O=Packt/CN=packtpub.com"
