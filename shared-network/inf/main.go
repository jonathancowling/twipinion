package main

import (
	"errors"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		
		zones, err := aws.GetAvailabilityZones(ctx, nil)

		if err != nil {
			return err
		}
		zoneNames := zones.Names[0:3]

		_, network, err := net.ParseCIDR("172.16.0.0/16")
		if err != nil {
			return err
		}

		vpc, err := ec2.NewVpc(ctx, "vpc", &ec2.VpcArgs{
			CidrBlock: pulumi.String(network.String()),
		})
		if err != nil {
			return err
		}

		gw, err := ec2.NewInternetGateway(ctx, "gateway", &ec2.InternetGatewayArgs{
			VpcId: vpc.ID(),
		})
		if err != nil {
			return err
		}

		subnets := make([]pulumi.IDOutput, len(zoneNames))
		for i, zone := range zoneNames {
			subnetCidr, err := cidr.Subnet(network, 8, i)
			if err != nil {
				return err
			}

			subnet, err := ec2.NewSubnet(ctx, "subnet-" + zone,  &ec2.SubnetArgs{
				VpcId: vpc.ID(),
				AvailabilityZone: pulumi.String(zone),
				CidrBlock: pulumi.String(subnetCidr.String()),
			})
			subnets[i] = subnet.ID()
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
		if len(subnets) == 0 {
			return errors.New("no default subnet not found in region")
		}

		// sg, err := ec2.NewSecurityGroup(ctx, "default", &ec2.SecurityGroupArgs{
		// 	VpcId: vpc.ID(),
		// 	Ingress: ec2.SecurityGroupIngressArray{
		// 		&ec2.SecurityGroupIngressArgs{
		// 			Protocol: pulumi.String("tcp"),
		// 			Self:     pulumi.Bool(true),
		// 			FromPort: pulumi.Int(0),
		// 			ToPort:   pulumi.Int(9094),
		// 			CidrBlocks: pulumi.StringArray{
		// 				vpc.CidrBlock,
		// 			},
		// 		},
		// 	},
		// 	Egress: ec2.SecurityGroupEgressArray{
		// 		&ec2.SecurityGroupEgressArgs{
		// 			FromPort: pulumi.Int(0),
		// 			ToPort:   pulumi.Int(0),
		// 			Protocol: pulumi.String("tcp"),
		// 			Self:     pulumi.Bool(true),
		// 			CidrBlocks: pulumi.StringArray{
		// 				vpc.CidrBlock,
		// 			},
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }
		_, err = ec2.NewRouteTable(ctx, "routetable", &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: gw.ID(),
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the name of the bucket
		ctx.Export("VPC ID", vpc.ID())
		ctx.Export("VPC CIDR", pulumi.String(network.String()))
		ctx.Export("Subnet IDs", pulumi.ToIDArrayOutput(subnets))
		// ctx.Export("Security Group", sg.ID())

		return nil
	})
}
