package main

import (
	"os"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"gopkg.in/yaml.v2"
)

type applicationConfig struct {
	Config struct {
		Twitter struct {
			Bearer string `yaml:"bearer"`
		} `yaml:"twitter"`
	} `yaml:"config"`
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")

		b, err := os.ReadFile("../app/src/main/resources/configserver/" + conf.Require("env") + "/application.yml")
		if err != nil {
			return err
		}
		
		appConfig := applicationConfig{}
		err = yaml.Unmarshal(b, &appConfig)
		
		if err != nil {
			return err
		}
		
		out := pulumi.String(appConfig.Config.Twitter.Bearer).ToStringOutput()

		_, err = ssm.NewParameter(ctx, "twitter-bearer-parameter", &ssm.ParameterArgs{
			Name: pulumi.String("/config/ingester-dev/config.twitter.bearer"),
			Type:  pulumi.String("SecureString"),
			Value: pulumi.ToSecret(out).(pulumi.StringOutput),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
