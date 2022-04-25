This is a baseline pulumi project to set up infrastructure other projects depend on (i.e. s3 & IAM).
Note that the infrastructure here is so we can create other pulumi projects not to create application infrastructure

# Deployment

To get all the infrastructure up and running you'll need to create the stack for it (bootstrap-dev),
login to the local backend and use pulumi up to get everything up and running.

> The bootstrap project must not contain any secrets, that way only the infrastructure contains secrets. (PULUMI_CONFIG_PASSPHRASE='')
> The project will however require being run manually using an account with AWS access.