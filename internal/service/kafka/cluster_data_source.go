// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kafka

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kafka/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKDataSource("aws_msk_cluster", name="Cluster")
func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceClusterRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bootstrap_brokers": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bootstrap_brokers_public_sasl_iam": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bootstrap_brokers_public_sasl_scram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bootstrap_brokers_public_tls": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bootstrap_brokers_sasl_iam": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bootstrap_brokers_sasl_scram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bootstrap_brokers_tls": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"cluster_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kafka_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"number_of_broker_nodes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"zookeeper_connect_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zookeeper_connect_string_tls": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).KafkaClient(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	clusterName := d.Get("cluster_name").(string)
	input := &kafka.ListClustersInput{
		ClusterNameFilter: aws.String(clusterName),
	}
	cluster, err := findCluster(ctx, conn, input, func(v *types.ClusterInfo) bool {
		return aws.ToString(v.ClusterName) == clusterName
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading MSK Cluster (%s): %s", clusterName, err)
	}

	clusterARN := aws.ToString(cluster.ClusterArn)
	bootstrapBrokersOutput, err := findBootstrapBrokersByARN(ctx, conn, clusterARN)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading MSK Cluster (%s) bootstrap brokers: %s", clusterARN, err)
	}

	d.SetId(clusterARN)
	d.Set("arn", clusterARN)
	d.Set("bootstrap_brokers", SortEndpointsString(aws.ToString(bootstrapBrokersOutput.BootstrapBrokerString)))
	d.Set("bootstrap_brokers_public_sasl_iam", SortEndpointsString(aws.ToString(bootstrapBrokersOutput.BootstrapBrokerStringPublicSaslIam)))
	d.Set("bootstrap_brokers_public_sasl_scram", SortEndpointsString(aws.ToString(bootstrapBrokersOutput.BootstrapBrokerStringPublicSaslScram)))
	d.Set("bootstrap_brokers_public_tls", SortEndpointsString(aws.ToString(bootstrapBrokersOutput.BootstrapBrokerStringPublicTls)))
	d.Set("bootstrap_brokers_sasl_iam", SortEndpointsString(aws.ToString(bootstrapBrokersOutput.BootstrapBrokerStringSaslIam)))
	d.Set("bootstrap_brokers_sasl_scram", SortEndpointsString(aws.ToString(bootstrapBrokersOutput.BootstrapBrokerStringSaslScram)))
	d.Set("bootstrap_brokers_tls", SortEndpointsString(aws.ToString(bootstrapBrokersOutput.BootstrapBrokerStringTls)))
	d.Set("cluster_name", cluster.ClusterName)
	clusterUUID, _ := clusterUUIDFromARN(clusterARN)
	d.Set("cluster_uuid", clusterUUID)
	d.Set("kafka_version", cluster.CurrentBrokerSoftwareInfo.KafkaVersion)
	d.Set("number_of_broker_nodes", cluster.NumberOfBrokerNodes)
	d.Set("zookeeper_connect_string", SortEndpointsString(aws.ToString(cluster.ZookeeperConnectString)))
	d.Set("zookeeper_connect_string_tls", SortEndpointsString(aws.ToString(cluster.ZookeeperConnectStringTls)))

	if err := d.Set("tags", KeyValueTags(ctx, cluster.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}

func findCluster(ctx context.Context, conn *kafka.Client, input *kafka.ListClustersInput, filter tfslices.Predicate[*types.ClusterInfo]) (*types.ClusterInfo, error) {
	output, err := findClusters(ctx, conn, input, filter)

	if err != nil {
		return nil, err
	}

	return tfresource.AssertFirstValueResult(output)
}

func findClusters(ctx context.Context, conn *kafka.Client, input *kafka.ListClustersInput, filter tfslices.Predicate[*types.ClusterInfo]) ([]types.ClusterInfo, error) {
	var output []types.ClusterInfo

	pages := kafka.NewListClustersPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		for _, v := range page.ClusterInfoList {
			if filter(&v) {
				output = append(output, v)
			}
		}
	}

	return output, nil
}
