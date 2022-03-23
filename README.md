# twipinion
Visualise a brands public perception

# Getting Started

# Bootstrapping

To get all the infrastructure up and running for the first time see the bootstrap project.

# Proxy

install mitmproxy
```sh
brew install mitmproxy
sudo keytool -import -trustcacerts -file ~/.mitmproxy/mitmproxy-ca-cert.pem -alias mitmproxy -cacerts
mvn spring-boot:run -Dspring-boot.run.jvmArguments='-Dhttps.proxyHost=localhost -Dhttps.proxyPort=8000'
```
