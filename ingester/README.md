# Ingester

This project will stream relevant tweets from twitter and make them available to subscribers.

# Running the Project

## Local

```sh
mvn spring-boot:run -Dspring-boot.run.jvmArguments='-Dspring.profiles.active=local -Dspring.cloud.bootstrap.name=bootstrap_local'
```

## Dev

```sh
mvn spring-boot:run -Dspring-boot.run.jvmArguments='-Dspring.profiles.active=dev -Dspring.cloud.bootstrap.name=bootstrap_dev' 
```

# Proxy

install mitmproxy
```sh
brew install mitmproxy
sudo keytool -import -trustcacerts -file ~/.mitmproxy/mitmproxy-ca-cert.pem -alias mitmproxy -cacerts
mitmweb --port 8000
mvn spring-boot:run -Dspring-boot.run.jvmArguments='-Dhttps.proxyHost=localhost -Dhttps.proxyPort=8000'
```