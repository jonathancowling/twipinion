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
		
		zones, err := aws.GetAvailabilityZones(ctx, &aws.GetAvailabilityZonesArgs{
			ExcludeNames: []string { "us-east-1e" },
			Filters: []aws.GetAvailabilityZonesFilter{
				{
					Name: "zone-type",
					Values: []string{ "availability-zone" },
				},
			},
		})

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

		rtb, err := ec2.NewRouteTable(ctx, "routetable", &ec2.RouteTableArgs{
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

		subnets := make([]pulumi.IDOutput, len(zoneNames))
		for i, zone := range zoneNames {
			// FIXME: temporarily shift the zones to support recreation without deletion
			subnetCidr, err := cidr.Subnet(network, 8, i + len(zoneNames))
			if err != nil {
				return err
			}

			subnet, err := ec2.NewSubnet(ctx, "subnet-" + zone,  &ec2.SubnetArgs{
					VpcId: vpc.ID(),
					AvailabilityZone: pulumi.String(zone),
					CidrBlock: pulumi.String(subnetCidr.String()),
				},
		    )
			subnets[i] = subnet.ID()
			if err != nil {
				return err
			}
			_, err = ec2.NewRouteTableAssociation(ctx, "routeTableAssociation", &ec2.RouteTableAssociationArgs{
				SubnetId:     subnet.ID(),
				RouteTableId: rtb.ID(),
			})
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

		// Export the name of the bucket
		ctx.Export("VPC ID", vpc.ID())
		ctx.Export("VPC CIDR", pulumi.String(network.String()))
		ctx.Export("Subnet IDs", pulumi.ToIDArrayOutput(subnets))

		return nil
	})
}
