# twipinion
Visualise a brands public perception

# Getting Started

1. Deploy bootstrap projects:
    1. /bootstrap (local backend)
    2. /ingester/bootstrap (s3 backend)
2. Init stacks (TODO: automate this)
    1. /injester/inf


# Bootstrapping

To get all the infrastructure up and running for the first time see the bootstrap project.

# Proxy

install mitmproxy
```sh
brew install mitmproxy
sudo keytool -import -trustcacerts -file ~/.mitmproxy/mitmproxy-ca-cert.pem -alias mitmproxy -cacerts
mvn spring-boot:run -Dspring-boot.run.jvmArguments='-Dhttps.proxyHost=localhost -Dhttps.proxyPort=8000'
```
