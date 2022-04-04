package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/msk"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		network, err := pulumi.NewStackReference(ctx, "shared-network-inf-" + conf.Require("env"), nil)
		if err != nil {
			return err
		}

		vpcID := network.GetOutput(pulumi.String("VPC ID")).
		    ApplyT(func (i interface {}) string { return i.(string) }).(pulumi.StringOutput)
		vpcCidr := network.GetOutput(pulumi.String("VPC CIDR")).
		    ApplyT(func (i interface {}) string { return i.(string) }).(pulumi.StringOutput)

		sg, err := ec2.NewSecurityGroup(ctx, "security-group", &ec2.SecurityGroupArgs{ // TODO: rename
			VpcId: vpcID,
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					Protocol: pulumi.String("tcp"),
					Self:     pulumi.Bool(true),
					FromPort: pulumi.Int(0),
					ToPort:   pulumi.Int(9094),
					CidrBlocks: pulumi.StringArray{ vpcCidr },
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
					FromPort: pulumi.Int(0),
					ToPort:   pulumi.Int(0),
					Protocol: pulumi.String("tcp"),
					Self:     pulumi.Bool(true),
					CidrBlocks: pulumi.StringArray{ vpcCidr },
				},
			},
		})
		if err != nil {
			return err
		}

		subnets := network.GetOutput(pulumi.String("Subnet IDs")).
		    ApplyT(func (sIds interface {}) []string {
				subnets := make([]string, len(sIds.([]interface{})))
				for i, id := range sIds.([]interface{}) {
					subnets[i] = id.(string)
				}
				return subnets
			}).(pulumi.StringArrayOutput)

		cluster, err := msk.NewCluster(ctx, "kafka", &msk.ClusterArgs{
			KafkaVersion:        pulumi.String("2.8.1"),
			NumberOfBrokerNodes: subnets.ApplyT(func(ids []string) int { return len(ids) }).(pulumi.IntOutput),
			BrokerNodeGroupInfo: &msk.ClusterBrokerNodeGroupInfoArgs{
				InstanceType: pulumi.String("kafka.t3.small"),
				ClientSubnets: subnets,
				EbsVolumeSize: pulumi.Int(1),
				SecurityGroups: sg.ID().
				    ToStringOutput().
					ApplyT(func(id string) []string { return []string{ id }}).(pulumi.StringArrayOutput),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("Zookeeper Connect String", cluster.ZookeeperConnectString)
		ctx.Export("Bootstrap Brokers TLS", cluster.BootstrapBrokersTls)

		// pulumi.All(
		// 	network.GetOutput(pulumi.String("Subnet IDs")),
		// 	sg.ID(),
		// ).ApplyT(func (ids interface{}, sg string) (interface {}) {
		// 	subnets := make([]string, len(ids.([]interface{})))
		// 	for i, id := range ids.([]interface{}) {
		// 		subnets[i] = id.(string)
		// 	}

		// 	cluster, e := msk.NewCluster(ctx, "kafka", &msk.ClusterArgs{
		// 		KafkaVersion:        pulumi.String("2.8.1"),
		// 		NumberOfBrokerNodes: pulumi.Int(2),
		// 		BrokerNodeGroupInfo: &msk.ClusterBrokerNodeGroupInfoArgs{
		// 			InstanceType: pulumi.String("kafka.t3.small"),
		// 			ClientSubnets: pulumi.ToStringArray(subnets),
		// 			EbsVolumeSize: pulumi.Int(1),
		// 			SecurityGroups: pulumi.ToStringArray([]string{ sg }),
		// 		},
		// 	})
		// 	ctx.Export("Zookeeper Connect String", cluster.ZookeeperConnectString)
		// 	ctx.Export("Bootstrap Brokers Tls", cluster.BootstrapBrokersTls)

		// 	if e != nil {
		// 		err = e
		// 		return nil
		// 	}

		// 	return nil
		// })

		return nil
	})
}
