// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package lakeformation

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	lakeformation_sdkv2 "github.com/aws/aws-sdk-go-v2/service/lakeformation"
	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	lakeformation_sdkv1 "github.com/aws/aws-sdk-go/service/lakeformation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newResourceDataCellsFilter,
			Name:    "Data Cells Filter",
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceDataLakeSettings,
			TypeName: "aws_lakeformation_data_lake_settings",
		},
		{
			Factory:  DataSourcePermissions,
			TypeName: "aws_lakeformation_permissions",
		},
		{
			Factory:  DataSourceResource,
			TypeName: "aws_lakeformation_resource",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceDataLakeSettings,
			TypeName: "aws_lakeformation_data_lake_settings",
		},
		{
			Factory:  ResourceLFTag,
			TypeName: "aws_lakeformation_lf_tag",
		},
		{
			Factory:  ResourcePermissions,
			TypeName: "aws_lakeformation_permissions",
		},
		{
			Factory:  ResourceResource,
			TypeName: "aws_lakeformation_resource",
			Name:     "Resource",
		},
		{
			Factory:  ResourceResourceLFTags,
			TypeName: "aws_lakeformation_resource_lf_tags",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.LakeFormation
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*lakeformation_sdkv1.LakeFormation, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return lakeformation_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*lakeformation_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return lakeformation_sdkv2.NewFromConfig(cfg, func(o *lakeformation_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
