// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rds

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_db_subnet_group", name="DB Subnet Group")
// @Tags(identifierAttribute="arn")
func resourceSubnetGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSubnetGroupCreate,
		ReadWithoutTimeout:   resourceSubnetGroupRead,
		UpdateWithoutTimeout: resourceSubnetGroupUpdate,
		DeleteWithoutTimeout: resourceSubnetGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validSubnetGroupName,
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validSubnetGroupNamePrefix,
			},
			"subnet_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"supported_network_types": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceSubnetGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RDSClient(ctx)

	name := create.Name(d.Get("name").(string), d.Get("name_prefix").(string))
	input := &rds.CreateDBSubnetGroupInput{
		DBSubnetGroupDescription: aws.String(d.Get("description").(string)),
		DBSubnetGroupName:        aws.String(name),
		SubnetIds:                flex.ExpandStringValueSet(d.Get("subnet_ids").(*schema.Set)),
		Tags:                     getTagsInV2(ctx),
	}

	output, err := conn.CreateDBSubnetGroup(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating RDS DB Subnet Group (%s): %s", name, err)
	}

	d.SetId(aws.ToString(output.DBSubnetGroup.DBSubnetGroupName))

	return append(diags, resourceSubnetGroupRead(ctx, d, meta)...)
}

func resourceSubnetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RDSClient(ctx)

	v, err := findDBSubnetGroupByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] RDS DB Subnet Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading RDS DB Subnet Group (%s): %s", d.Id(), err)
	}

	d.Set("arn", v.DBSubnetGroupArn)
	d.Set("description", v.DBSubnetGroupDescription)
	d.Set("name", v.DBSubnetGroupName)
	d.Set("name_prefix", create.NamePrefixFromName(aws.ToString(v.DBSubnetGroupName)))
	d.Set("subnet_ids", tfslices.ApplyToAll(v.Subnets, func(v types.Subnet) string {
		return aws.ToString(v.SubnetIdentifier)
	}))
	d.Set("supported_network_types", v.SupportedNetworkTypes)
	d.Set("vpc_id", v.VpcId)

	return diags
}

func resourceSubnetGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RDSClient(ctx)

	if d.HasChanges("description", "subnet_ids") {
		input := &rds.ModifyDBSubnetGroupInput{
			DBSubnetGroupDescription: aws.String(d.Get("description").(string)),
			DBSubnetGroupName:        aws.String(d.Id()),
			SubnetIds:                flex.ExpandStringValueSet(d.Get("subnet_ids").(*schema.Set)),
		}

		_, err := conn.ModifyDBSubnetGroup(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating RDS DB Subnet Group (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceSubnetGroupRead(ctx, d, meta)...)
}

func resourceSubnetGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RDSClient(ctx)

	log.Printf("[DEBUG] Deleting RDS DB Subnet Group: %s", d.Id())
	_, err := conn.DeleteDBSubnetGroup(ctx, &rds.DeleteDBSubnetGroupInput{
		DBSubnetGroupName: aws.String(d.Id()),
	})

	if errs.IsA[*types.DBSubnetGroupNotFoundFault](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting RDS Subnet Group (%s): %s", d.Id(), err)
	}

	_, err = tfresource.RetryUntilNotFound(ctx, 3*time.Minute, func() (interface{}, error) {
		return findDBSubnetGroupByName(ctx, conn, d.Id())
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for RDS Subnet Group (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func findDBSubnetGroupByName(ctx context.Context, conn *rds.Client, name string) (*types.DBSubnetGroup, error) {
	input := &rds.DescribeDBSubnetGroupsInput{
		DBSubnetGroupName: aws.String(name),
	}
	output, err := findDBSubnetGroup(ctx, conn, input, tfslices.PredicateTrue[*types.DBSubnetGroup]())

	if err != nil {
		return nil, err
	}

	// Eventual consistency check.
	if aws.ToString(output.DBSubnetGroupName) != name {
		return nil, &retry.NotFoundError{
			LastRequest: input,
		}
	}

	return output, nil
}

func findDBSubnetGroup(ctx context.Context, conn *rds.Client, input *rds.DescribeDBSubnetGroupsInput, filter tfslices.Predicate[*types.DBSubnetGroup]) (*types.DBSubnetGroup, error) {
	output, err := findDBSubnetGroups(ctx, conn, input, filter)

	if err != nil {
		return nil, err
	}

	return tfresource.AssertSingleValueResult(output)
}

func findDBSubnetGroups(ctx context.Context, conn *rds.Client, input *rds.DescribeDBSubnetGroupsInput, filter tfslices.Predicate[*types.DBSubnetGroup]) ([]types.DBSubnetGroup, error) {
	var output []types.DBSubnetGroup

	pages := rds.NewDescribeDBSubnetGroupsPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if errs.IsA[*types.DBSubnetGroupNotFoundFault](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: input,
			}
		}

		if err != nil {
			return nil, err
		}

		for _, v := range page.DBSubnetGroups {
			if filter(&v) {
				output = append(output, v)
			}
		}
	}

	return output, nil
}
