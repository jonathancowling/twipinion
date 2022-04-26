package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws-native/sdk/go/aws"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/iam"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/kms"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")

		region, err := aws.GetRegion(ctx)
		if err != nil {
			return err
		}
		ctx.Export("aws region", pulumi.String(region.Region))

		bucket, err := s3.NewBucket(ctx, "backend-bucket", nil)
		if err != nil {
			return err
		}
		ctx.Export("pulumi backend", pulumi.Sprintf("s3://%s", bucket.BucketName.Elem()))

		oidc, err := iam.NewOIDCProvider(
			ctx,
			"github-oidc",
			&iam.OIDCProviderArgs{
				Url: pulumi.String("https://token.actions.githubusercontent.com"),
				ClientIdList: pulumi.ToStringArray([]string{
					"sts.amazonaws.com",
				}),
				ThumbprintList: pulumi.ToStringArray([]string{
					"6938fd4d98bab03faadb97b34396831e3780aea1", // GitHub's thumbprint
				}),
			},
		)
		if err != nil {
			return err
		}
		ctx.Export("OIDC", oidc.Url)

		ciPolicyDocument := oidc.Arn.ApplyT(func(arn string) string {
			return `{
				"Version": "2012-10-17",
				"Statement": {
					"Effect": "Allow",
					"Principal": {"Federated": "` + arn + `"},
					"Action": "sts:AssumeRoleWithWebIdentity",
					"Condition": {
						"StringLike": {
						  "token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
						  "token.actions.githubusercontent.com:sub": "repo:` + conf.Require("repo") + `:*"
						}
					}
				}
			}`
		})

		ciRole, err := iam.NewRole(ctx, "ci-role", &iam.RoleArgs{
			AssumeRolePolicyDocument: ciPolicyDocument,
			Description: pulumi.String("IAM Role for CI servers to use"),
			ManagedPolicyArns: pulumi.ToStringArray([]string{
				"arn:aws:iam::aws:policy/IAMFullAccess",
				"arn:aws:iam::aws:policy/AmazonS3FullAccess",
				"arn:aws:iam::aws:policy/AWSKeyManagementServicePowerUser",
				"arn:aws:iam::aws:policy/CloudWatchEventsFullAccess",
				"arn:aws:iam::aws:policy/AWSCloudFormationFullAccess",
				"arn:aws:iam::aws:policy/AWSLambda_FullAccess",
				"arn:aws:iam::aws:policy/AmazonEC2FullAccess",
				"arn:aws:iam::aws:policy/AmazonMSKFullAccess",
				"arn:aws:iam::aws:policy/AmazonSSMFullAccess",
			}),
		})
		if err != nil {
			return err
		}
		ctx.Export("aws ci role", ciRole.Arn)

		accountId, err := aws.GetAccountId(ctx)
		if err != nil {
			return err
		}

		template, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"AWS": []string{
							"arn:aws:iam::" + accountId.AccountId + ":user/cloud_user",
							"%s",
						},
					},
					"Action":   "kms:*",
					"Resource": "*",
				},
			},
		})
		if err != nil {
			return err
		}
		policy := pulumi.Sprintf(string(template), ciRole.Arn)
		key, err := kms.NewKey(ctx, "secret-key", &kms.KeyArgs{
			Description: pulumi.String("secret key for pulumi"),
			KeyPolicy:   policy,
		})
		if err != nil {
			return err
		}
		alias, err := kms.NewAlias(ctx, "secret-key-alias", &kms.AliasArgs{
			AliasName:   pulumi.StringPtr("alias/pulumi"),
			TargetKeyId: key.KeyId,
		})
		if err != nil {
			return err
		}
		
		ctx.Export(
			"pulumi secrets provider",
			pulumi.Sprintf("awskms://%s?region=%s", alias.AliasName, region.Region),
		)

		return nil
	})
}
