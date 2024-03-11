// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package codecommit

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_codecommit_repository", name="Repository")
func dataSourceRepository() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceRepositoryRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clone_url_http": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clone_url_ssh": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
			},
		},
	}
}

func dataSourceRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CodeCommitClient(ctx)

	name := d.Get("repository_name").(string)
	repository, err := findRepositoryByName(ctx, conn, name)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading CodeCommit Repository (%s): %s", name, err)
	}

	d.SetId(aws.ToString(repository.RepositoryName))
	d.Set("arn", repository.Arn)
	d.Set("clone_url_http", repository.CloneUrlHttp)
	d.Set("clone_url_ssh", repository.CloneUrlSsh)
	d.Set("kms_key_id", repository.KmsKeyId)
	d.Set("repository_id", repository.RepositoryId)
	d.Set("repository_name", repository.RepositoryName)

	return diags
}
