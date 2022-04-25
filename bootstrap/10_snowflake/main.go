package main

import (
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi-snowflake/sdk/go/snowflake"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {

		password, err := random.NewRandomPassword(ctx, "password", &random.RandomPasswordArgs{
			Length:          pulumi.Int(64),
			Special:         pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		user, err := snowflake.NewUser(ctx, "ci-user", &snowflake.UserArgs{
			Comment:            pulumi.String("Snowflake CI user"),
			DefaultRole:        pulumi.String("SYSADMIN"),
			DisplayName:        pulumi.String("CI User"),
			Email:              pulumi.String("jonathan.cowling+twipinion-ci@infinityworks.com"),
			FirstName:          pulumi.String("Twipinion"),
			LastName:           pulumi.String("CI User"),
			LoginName:          pulumi.String("Twipinion CI User"),
			Password:       	password.Result,
		})
		if err != nil {
			return err
		}

		ctx.Export("snowflake login", user.LoginName)
		ctx.Export("snowflake password", password.Result)

		return nil
	})
}
