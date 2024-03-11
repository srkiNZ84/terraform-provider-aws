// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package efs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_efs_replication_configuration", name="Replication Configuration")
func ResourceReplicationConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceReplicationConfigurationCreate,
		ReadWithoutTimeout:   resourceReplicationConfigurationRead,
		DeleteWithoutTimeout: resourceReplicationConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: []string{"destination.0.availability_zone_name", "destination.0.region"},
						},
						"file_system_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"region": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: verify.ValidRegionName,
							AtLeastOneOf: []string{"destination.0.availability_zone_name", "destination.0.region"},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"original_source_file_system_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_file_system_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_file_system_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceReplicationConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EFSConn(ctx)

	fsID := d.Get("source_file_system_id").(string)
	input := &efs.CreateReplicationConfigurationInput{
		SourceFileSystemId: aws.String(fsID),
	}

	if v, ok := d.GetOk("destination"); ok && len(v.([]interface{})) > 0 {
		input.Destinations = expandDestinationsToCreate(v.([]interface{}))
	}

	_, err := conn.CreateReplicationConfigurationWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EFS Replication Configuration (%s): %s", fsID, err)
	}

	d.SetId(fsID)

	if _, err := waitReplicationConfigurationCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for EFS Replication Configuration (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceReplicationConfigurationRead(ctx, d, meta)...)
}

func resourceReplicationConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EFSConn(ctx)

	replication, err := FindReplicationConfigurationByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EFS Replication Configuration (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EFS Replication Configuration (%s): %s", d.Id(), err)
	}

	destinations := flattenDestinations(replication.Destinations)

	// availability_zone_name and kms_key_id aren't returned from the AWS Read API.
	if v, ok := d.GetOk("destination"); ok && len(v.([]interface{})) > 0 {
		copy := func(i int, k string) {
			destinations[i].(map[string]interface{})[k] = v.([]interface{})[i].(map[string]interface{})[k]
		}
		// Assume 1 destination.
		copy(0, "availability_zone_name")
		copy(0, "kms_key_id")
	}

	d.Set("creation_time", aws.TimeValue(replication.CreationTime).String())
	if err := d.Set("destination", destinations); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting destination: %s", err)
	}
	d.Set("original_source_file_system_arn", replication.OriginalSourceFileSystemArn)
	d.Set("source_file_system_arn", replication.SourceFileSystemArn)
	d.Set("source_file_system_id", replication.SourceFileSystemId)
	d.Set("source_file_system_region", replication.SourceFileSystemRegion)

	return diags
}

func resourceReplicationConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EFSConn(ctx)

	// Deletion of the replication configuration must be done from the Region in which the destination file system is located.
	destination := expandDestinationsToCreate(d.Get("destination").([]interface{}))[0]
	region := aws.StringValue(destination.Region)
	session, err := conns.NewSessionForRegion(&conn.Config, region, meta.(*conns.AWSClient).TerraformVersion)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating AWS session (%s): %s", region, err)
	}

	log.Printf("[DEBUG] Deleting EFS Replication Configuration: %s", d.Id())
	if err := deleteReplicationConfiguration(ctx, efs.New(session), d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	// Delete also in the source Region.
	if err := deleteReplicationConfiguration(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	return diags
}

func deleteReplicationConfiguration(ctx context.Context, conn *efs.EFS, fsID string, timeout time.Duration) error {
	_, err := conn.DeleteReplicationConfigurationWithContext(ctx, &efs.DeleteReplicationConfigurationInput{
		SourceFileSystemId: aws.String(fsID),
	})

	if tfawserr.ErrCodeEquals(err, efs.ErrCodeFileSystemNotFound, efs.ErrCodeReplicationNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting EFS Replication Configuration (%s): %w", fsID, err)
	}

	if _, err := waitReplicationConfigurationDeleted(ctx, conn, fsID, timeout); err != nil {
		return fmt.Errorf("waiting for EFS Replication Configuration (%s) delete: %w", fsID, err)
	}

	return nil
}

func findReplicationConfiguration(ctx context.Context, conn *efs.EFS, input *efs.DescribeReplicationConfigurationsInput) (*efs.ReplicationConfigurationDescription, error) {
	output, err := findReplicationConfigurations(ctx, conn, input)

	if err != nil {
		return nil, err
	}

	return tfresource.AssertSinglePtrResult(output)
}

func findReplicationConfigurations(ctx context.Context, conn *efs.EFS, input *efs.DescribeReplicationConfigurationsInput) ([]*efs.ReplicationConfigurationDescription, error) {
	var output []*efs.ReplicationConfigurationDescription

	err := conn.DescribeReplicationConfigurationsPagesWithContext(ctx, input, func(page *efs.DescribeReplicationConfigurationsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Replications {
			if v != nil {
				output = append(output, v)
			}
		}

		return !lastPage
	})

	if tfawserr.ErrCodeEquals(err, efs.ErrCodeFileSystemNotFound, efs.ErrCodeReplicationNotFound) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

func FindReplicationConfigurationByID(ctx context.Context, conn *efs.EFS, id string) (*efs.ReplicationConfigurationDescription, error) {
	input := &efs.DescribeReplicationConfigurationsInput{
		FileSystemId: aws.String(id),
	}

	output, err := findReplicationConfiguration(ctx, conn, input)

	if err != nil {
		return nil, err
	}

	if len(output.Destinations) == 0 || output.Destinations[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func statusReplicationConfiguration(ctx context.Context, conn *efs.EFS, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindReplicationConfigurationByID(ctx, conn, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Destinations[0].Status), nil
	}
}

func waitReplicationConfigurationCreated(ctx context.Context, conn *efs.EFS, id string, timeout time.Duration) (*efs.ReplicationConfigurationDescription, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{efs.ReplicationStatusEnabling},
		Target:  []string{efs.ReplicationStatusEnabled},
		Refresh: statusReplicationConfiguration(ctx, conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*efs.ReplicationConfigurationDescription); ok {
		return output, err
	}

	return nil, err
}

func waitReplicationConfigurationDeleted(ctx context.Context, conn *efs.EFS, id string, timeout time.Duration) (*efs.ReplicationConfigurationDescription, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{efs.ReplicationStatusDeleting},
		Target:                    []string{},
		Refresh:                   statusReplicationConfiguration(ctx, conn, id),
		Timeout:                   timeout,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*efs.ReplicationConfigurationDescription); ok {
		return output, err
	}

	return nil, err
}

func expandDestinationToCreate(tfMap map[string]interface{}) *efs.DestinationToCreate {
	if tfMap == nil {
		return nil
	}

	apiObject := &efs.DestinationToCreate{}

	if v, ok := tfMap["availability_zone_name"].(string); ok && v != "" {
		apiObject.AvailabilityZoneName = aws.String(v)
	}

	if v, ok := tfMap["kms_key_id"].(string); ok && v != "" {
		apiObject.KmsKeyId = aws.String(v)
	}

	if v, ok := tfMap["region"].(string); ok && v != "" {
		apiObject.Region = aws.String(v)
	}

	if v, ok := tfMap["file_system_id"].(string); ok && v != "" {
		apiObject.FileSystemId = aws.String(v)
	}

	return apiObject
}

func expandDestinationsToCreate(tfList []interface{}) []*efs.DestinationToCreate {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*efs.DestinationToCreate

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandDestinationToCreate(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func flattenDestination(apiObject *efs.Destination) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.FileSystemId; v != nil {
		tfMap["file_system_id"] = aws.StringValue(v)
	}

	if v := apiObject.Region; v != nil {
		tfMap["region"] = aws.StringValue(v)
	}

	if v := apiObject.Status; v != nil {
		tfMap["status"] = aws.StringValue(v)
	}

	return tfMap
}

func flattenDestinations(apiObjects []*efs.Destination) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenDestination(apiObject))
	}

	return tfList
}
