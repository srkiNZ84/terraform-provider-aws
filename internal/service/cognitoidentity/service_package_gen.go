// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package cognitoidentity

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	cognitoidentity_sdkv2 "github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
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
			Factory:  DataSourcePool,
			TypeName: "aws_cognito_identity_pool",
			Name:     "Pool",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourcePool,
			TypeName: "aws_cognito_identity_pool",
			Name:     "Pool",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourcePoolProviderPrincipalTag,
			TypeName: "aws_cognito_identity_pool_provider_principal_tag",
		},
		{
			Factory:  ResourcePoolRolesAttachment,
			TypeName: "aws_cognito_identity_pool_roles_attachment",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CognitoIdentity
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*cognitoidentity_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return cognitoidentity_sdkv2.NewFromConfig(cfg, func(o *cognitoidentity_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}