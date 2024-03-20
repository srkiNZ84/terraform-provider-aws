// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package cognitoidp

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	cognitoidentityprovider_sdkv1 "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{
		{
			Factory: newDataSourceDataSourceUserGroup,
			Name:    "User Group",
		},
		{
			Factory: newDataSourceDataSourceUserGroups,
			Name:    "User Groups",
		},
	}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newResourceManagedUserPoolClient,
		},
		{
			Factory: newResourceUserPoolClient,
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceUserPoolClient,
			TypeName: "aws_cognito_user_pool_client",
		},
		{
			Factory:  DataSourceUserPoolClients,
			TypeName: "aws_cognito_user_pool_clients",
		},
		{
			Factory:  DataSourceUserPoolSigningCertificate,
			TypeName: "aws_cognito_user_pool_signing_certificate",
		},
		{
			Factory:  DataSourceUserPools,
			TypeName: "aws_cognito_user_pools",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceIdentityProvider,
			TypeName: "aws_cognito_identity_provider",
		},
		{
			Factory:  ResourceResourceServer,
			TypeName: "aws_cognito_resource_server",
		},
		{
			Factory:  ResourceRiskConfiguration,
			TypeName: "aws_cognito_risk_configuration",
		},
		{
			Factory:  ResourceUser,
			TypeName: "aws_cognito_user",
		},
		{
			Factory:  resourceUserGroup,
			TypeName: "aws_cognito_user_group",
			Name:     "User Group",
		},
		{
			Factory:  ResourceUserInGroup,
			TypeName: "aws_cognito_user_in_group",
		},
		{
			Factory:  ResourceUserPool,
			TypeName: "aws_cognito_user_pool",
			Name:     "User Pool",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  ResourceUserPoolDomain,
			TypeName: "aws_cognito_user_pool_domain",
		},
		{
			Factory:  ResourceUserPoolUICustomization,
			TypeName: "aws_cognito_user_pool_ui_customization",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CognitoIDP
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*cognitoidentityprovider_sdkv1.CognitoIdentityProvider, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return cognitoidentityprovider_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}