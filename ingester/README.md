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
