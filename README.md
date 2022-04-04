# twipinion

Understand a brands public perception

# Getting Started

1. Run bootstrap script to bootstrap the project (scripts/bootstrap.sh)
2. Push commits to main branch
   (push an empty commit to deploy to a new environments without a code change)

# Bootstrapping

To get all the infrastructure up and running for the first time,
or deploy the project to a new AWS account use the refresh script
(requires being logged in to AWS locally).

The `/bootstrap` directory containins infrastructure code to enable
applications to be created & deployed,
while a `<application>/bootstrap` directory exists for any application
that needs extra base infrastructure in place to deploy it (e.g. secrets).

# Technology

- IaC [Pulumi](https://www.pulumi.com/)
- Cloud Provider [AWS](https://aws.amazon.com/)
- CI/CD [GitHub Actions](https://docs.github.com/en/actions)
- Applications [Spring Boot](https://spring.io/projects/spring-boot)
- Message Passing [Apache Kafka](https://kafka.apache.org/)
- (Coming Soon) Store [Snowflake]()
- (Coming Soon) Frontend

# Architecture Principles

- Applications are decoupled with message passing
- Applications should be serverless if possible
- Infrastructure should require minimal manual intervention
