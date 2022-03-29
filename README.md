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

## Root

The root bootstrap is necessary to ensure applications can be deployed to the environment and enables ci-cd etc.
This must be deployed first (and must not require secrets).

Re create bootstrap (when sandbox has expired)
```sh
pulumi login file://.
cd bootstrap/
PULUMI_CONFIG_PASSPHRASE='' pulumi stack rm --preserve-config --yes --force
PULUMI_CONFIG_PASSPHRASE='' pulumi stack init --stack bootstrap-dev
```

```sh
cd bootstrap
PULUMI_CONFIG_PASSPHRASE='' pulumi up
```

## Application

Applications may need bootstrapping too (look for an ${application}/bootstrap directory).
These are deployed using root bootstrap outputs.

```
pulumi login file://.
BUCKET_NAME="$(PULUMI_CONFIG_PASSPHRASE='' pulumi stack output --stack bootstrap-dev 'bucket name')"
SECRETS_PROVIDER="$(PULUMI_CONFIG_PASSPHRASE='' pulumi stack output --stack bootstrap-dev 'secrets provider')"
pulumi login "${BUCKET_NAME}"

application='application-to-deploy'
cd ${application}/bootstrap
sed '/^encryptedkey:/d' Pulumi.${application}-bootstrap-dev.yaml > tmp.yml && mv tmp.yml Pulumi.${application}-bootstrap-dev.yaml
pulumi stack init --stack ${application}-bootstrap-dev --secrets-provider "${SECRETS_PROVIDER}"
pulumi up
```

# Proxy

install mitmproxy
```sh
brew install mitmproxy
sudo keytool -import -trustcacerts -file ~/.mitmproxy/mitmproxy-ca-cert.pem -alias mitmproxy -cacerts
mitmweb --port 8000
mvn spring-boot:run -Dspring-boot.run.jvmArguments='-Dhttps.proxyHost=localhost -Dhttps.proxyPort=8000'
```
