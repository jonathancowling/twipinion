package main

import (
	"github.com/pulumi/pulumi-snowflake/sdk/go/snowflake"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		_, err := snowflake.NewDatabase(ctx, "test2", &snowflake.DatabaseArgs{
			Comment: pulumi.String("test comment 2"),
		})
		if err != nil {
			return err
		}

		schema, err := snowflake.NewSchema(ctx, "schema", &snowflake.SchemaArgs{
			Database:          pulumi.String("database"),
			DataRetentionDays: pulumi.Int(1),
		})
		if err != nil {
			return err
		}
		sequence, err := snowflake.NewSequence(ctx, "sequence", &snowflake.SequenceArgs{
			Database: schema.Database,
			Schema:   schema.Name,
		})
		if err != nil {
			return err
		}
		_, err = snowflake.NewTable(ctx, "table", &snowflake.TableArgs{
			Database: schema.Database,
			Schema:   schema.Name,
			Comment:  pulumi.String("A table."),
			ClusterBies: pulumi.StringArray{
				pulumi.String("to_date(DATE)"),
			},
			DataRetentionDays: schema.DataRetentionDays,
			ChangeTracking:    pulumi.Bool(false),
			Columns: TableColumnArray{
				&TableColumnArgs{
					Name:     pulumi.String("id"),
					Type:     pulumi.String("int"),
					Nullable: pulumi.Bool(true),
					Default: &TableColumnDefaultArgs{
						Sequence: sequence.FullyQualifiedName,
					},
				},
				&TableColumnArgs{
					Name:     pulumi.String("identity"),
					Type:     pulumi.String("NUMBER(38,0)"),
					Nullable: pulumi.Bool(true),
					Identity: &TableColumnIdentityArgs{
						StartNum: pulumi.Int(1),
						StepNum:  pulumi.Int(3),
					},
				},
				&TableColumnArgs{
					Name:     pulumi.String("data"),
					Type:     pulumi.String("text"),
					Nullable: pulumi.Bool(false),
				},
				&TableColumnArgs{
					Name: pulumi.String("DATE"),
					Type: pulumi.String("TIMESTAMP_NTZ(9)"),
				},
				&TableColumnArgs{
					Name:    pulumi.String("extra"),
					Type:    pulumi.String("VARIANT"),
					Comment: pulumi.String("extra data"),
				},
			},
			PrimaryKey: &TablePrimaryKeyArgs{
				Name: pulumi.String("my_key"),
				Keys: pulumi.StringArray{
					pulumi.String("data"),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = snowflake.NewWarehouse(ctx, "warehouse", &snowflake.WarehouseArgs{
			Comment:       pulumi.String("foo"),
			WarehouseSize: pulumi.String("small"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
