// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package elasticbeanstalk

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	elasticbeanstalk_sdkv2 "github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceApplication,
			TypeName: "aws_elastic_beanstalk_application",
		},
		{
			Factory:  DataSourceHostedZone,
			TypeName: "aws_elastic_beanstalk_hosted_zone",
		},
		{
			Factory:  DataSourceSolutionStack,
			TypeName: "aws_elastic_beanstalk_solution_stack",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceApplication,
			TypeName: "aws_elastic_beanstalk_application",
			Name:     "Application",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceApplicationVersion,
			TypeName: "aws_elastic_beanstalk_application_version",
			Name:     "Application Version",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceConfigurationTemplate,
			TypeName: "aws_elastic_beanstalk_configuration_template",
		},
		{
			Factory:  ResourceEnvironment,
			TypeName: "aws_elastic_beanstalk_environment",
			Name:     "Environment",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.ElasticBeanstalk
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*elasticbeanstalk_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return elasticbeanstalk_sdkv2.NewFromConfig(cfg, func(o *elasticbeanstalk_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}