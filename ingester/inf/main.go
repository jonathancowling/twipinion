package main

import (
	"encoding/json"
	"ingester/iampolicy"
	"ingester/pom"

	// native client doesn't yey support Service linked roles
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)


func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		
		pomFile, err := pom.LoadDefault()
		if err != nil {
			return err
		}
		
		policy, err := json.Marshal(iampolicy.IamPolicy{
			Version: "2012-10-17",
			Statement: iampolicy.Statement{
				Principal: &iampolicy.Principal{
					Service: []string{"lambda.amazonaws.com"},
				},
				Effect: "Allow",
				Action: []string{"sts:AssumeRole"},
			},
		})
		if err != nil {
			return err
		}

		role, err := iam.NewRole(
			ctx,
			pomFile.ArtifactId + "-function-role",
			&iam.RoleArgs{
				AssumeRolePolicy: pulumi.String(policy),
				ManagedPolicyArns: pulumi.ToStringArrayOutput([]pulumi.StringOutput{
					iam.ManagedPolicy("arn:aws:iam::aws:policy/AWSLambdaExecute").ToStringOutput(),
					iam.ManagedPolicy("arn:aws:iam::aws:policy/CloudWatchLogsFullAccess").ToStringOutput(),
				}),

			},
		)
		if err != nil {
			return err
		}

		bucket, err := s3.NewBucket(ctx, pomFile.ArtifactId + "-src", nil)
		if err != nil {
			return err
		}

		uploadedJar, err := s3.NewBucketObject(ctx, "lambda-jar", &s3.BucketObjectArgs{
			Key: pulumi.String("jar"),
			Bucket: bucket.ID(),
			Source: pulumi.NewFileArchive("../app/target/" + pomFile.SuffixedJar("aws")),
		})
		if err != nil {
			return err
		}

		function, err := lambda.NewFunction(ctx, pomFile.ArtifactId + "-function", &lambda.FunctionArgs{
			Runtime: pulumi.String("java11"),
			Role:    role.Arn,
			Handler: pulumi.String("org.springframework.cloud.function.adapter.aws.FunctionInvoker::handleRequest"),
			S3Bucket: bucket.ID(),
			S3Key: uploadedJar.Key,
			Environment: &lambda.FunctionEnvironmentArgs{
				Variables: pulumi.StringMap{
					"FUNCTION_NAME": pulumi.String("test"),
					"SPRING_PROFILES_ACTIVE": pulumi.String("dev"),
					"SPRING_CLOUD_BOOTSTRAP_NAME": pulumi.String("bootstrap_dev"),
					"LOGGING_LEVEL_ROOT": pulumi.String("DEBUG"),
				},
			},
			MemorySize: pulumi.Int(512),
			Description: pulumi.String(pomFile.Description),
			Timeout: pulumi.Int(120),
		})
		if err != nil {
			return err
		}

		rule, err := cloudwatch.NewEventRule(ctx, pomFile.ArtifactId + "-schedule", &cloudwatch.EventRuleArgs{
			ScheduleExpression: pulumi.String("rate(5 minutes)"),
			Description: pulumi.String("Schedule tweet ingest"),
		})
		if err != nil {
			return err
		}
		cloudwatch.NewEventTarget(ctx, pomFile.ArtifactId + "-schedule-target", &cloudwatch.EventTargetArgs{
			Arn: function.Arn,
			Rule: rule.Name,
		})
		

		// TODO: msk

		return nil
	})
}
