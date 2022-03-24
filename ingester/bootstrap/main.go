package main

import (
	"github.com/joho/godotenv"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		
		env, err := godotenv.Read("../app/src/main/resources/.env." + conf.Require("env"))
		if err != nil {
			return err
		}

		_, err = ssm.NewParameter(ctx, "twitter-bearer-parameter", &ssm.ParameterArgs{
			Name: pulumi.String("/config/ingester-dev/TWITTER_BEARER"),
			Type:  pulumi.String("SecureString"),
			Value: pulumi.String(env["TWITTER_BEARER"]), // TODO: check pulumi secret docs to see if this is secure
		})
		if err != nil {
			return err
		}
		return nil
	})
}
