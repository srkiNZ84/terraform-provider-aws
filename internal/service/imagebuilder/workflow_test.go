// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package imagebuilder_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/imagebuilder"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfimagebuilder "github.com/hashicorp/terraform-provider-aws/internal/service/imagebuilder"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccImageBuilderWorkflow_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_workflow.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ImageBuilderServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowConfig_name(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "imagebuilder", regexache.MustCompile(fmt.Sprintf("workflow/test/%s/1.0.0/[1-9][0-9]*", rName))),
					resource.TestCheckResourceAttr(resourceName, "change_description", ""),
					resource.TestMatchResourceAttr(resourceName, "data", regexache.MustCompile(`schemaVersion`)),
					acctest.CheckResourceAttrRFC3339(resourceName, "date_created"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "kms_key_id", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					acctest.CheckResourceAttrAccountID(resourceName, "owner"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "type", imagebuilder.WorkflowTypeTest),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccImageBuilderWorkflow_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_workflow.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ImageBuilderServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowConfig_name(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfimagebuilder.ResourceWorkflow(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccImageBuilderWorkflow_changeDescription(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_workflow.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ImageBuilderServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowConfig_changeDescription(rName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "change_description", "description1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccImageBuilderWorkflow_description(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_workflow.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ImageBuilderServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowConfig_description(rName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccImageBuilderWorkflow_kmsKeyID(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_workflow.test"
	kmsKeyResourceName := "aws_kms_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ImageBuilderServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowConfig_kmsKeyID(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_id", kmsKeyResourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccImageBuilderWorkflow_tags(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_workflow.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ImageBuilderServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWorkflowConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccWorkflowConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccImageBuilderWorkflow_uri(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_imagebuilder_workflow.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ImageBuilderServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckWorkflowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowConfig_uri(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWorkflowExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "data"),
					resource.TestCheckResourceAttrSet(resourceName, "uri"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"uri"},
			},
		},
	})
}

func testAccCheckWorkflowDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ImageBuilderConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_imagebuilder_workflow" {
				continue
			}

			input := &imagebuilder.GetWorkflowInput{
				WorkflowBuildVersionArn: aws.String(rs.Primary.ID),
			}

			output, err := conn.GetWorkflowWithContext(ctx, input)

			if tfawserr.ErrCodeEquals(err, imagebuilder.ErrCodeResourceNotFoundException) {
				continue
			}

			if err != nil {
				return fmt.Errorf("error getting Image Builder Workflow (%s): %w", rs.Primary.ID, err)
			}

			if output != nil {
				return fmt.Errorf("Image Builder Workflow (%s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckWorkflowExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ImageBuilderConn(ctx)

		input := &imagebuilder.GetWorkflowInput{
			WorkflowBuildVersionArn: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetWorkflowWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("error getting Image Builder Workflow (%s): %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccWorkflowConfig_name(rName string) string {
	return fmt.Sprintf(`
resource "aws_imagebuilder_workflow" "test" {
  name    = %[1]q
  version = "1.0.0"
  type    = "TEST"

  data = <<-EOT
  name: test-image
  description: Workflow to test an image
  schemaVersion: 1.0

  parameters:
    - name: waitForActionAtEnd
      type: boolean

  steps:
    - name: LaunchTestInstance
      action: LaunchInstance
      onFailure: Abort
      inputs:
        waitFor: "ssmAgent"

    - name: TerminateTestInstance
      action: TerminateInstance
      onFailure: Continue
      inputs:
        instanceId.$: "$.stepOutputs.LaunchTestInstance.instanceId"

    - name: WaitForActionAtEnd
      action: WaitForAction
      if:
        booleanEquals: true
        value: "$.parameters.waitForActionAtEnd"
  EOT
}
`, rName)
}

func testAccWorkflowConfig_changeDescription(rName, changeDescription string) string {
	return fmt.Sprintf(`
resource "aws_imagebuilder_workflow" "test" {
  name    = %[1]q
  version = "1.0.0"
  type    = "TEST"

  change_description = %[2]q

  data = <<-EOT
  name: test-image
  description: Workflow to test an image
  schemaVersion: 1.0

  parameters:
    - name: waitForActionAtEnd
      type: boolean

  steps:
    - name: LaunchTestInstance
      action: LaunchInstance
      onFailure: Abort
      inputs:
        waitFor: "ssmAgent"

    - name: TerminateTestInstance
      action: TerminateInstance
      onFailure: Continue
      inputs:
        instanceId.$: "$.stepOutputs.LaunchTestInstance.instanceId"

    - name: WaitForActionAtEnd
      action: WaitForAction
      if:
        booleanEquals: true
        value: "$.parameters.waitForActionAtEnd"
  EOT
}
`, rName, changeDescription)
}

func testAccWorkflowConfig_description(rName, description string) string {
	return fmt.Sprintf(`
resource "aws_imagebuilder_workflow" "test" {
  name    = %[1]q
  version = "1.0.0"
  type    = "TEST"

  description = %[2]q

  data = <<-EOT
  name: test-image
  description: Workflow to test an image
  schemaVersion: 1.0

  parameters:
    - name: waitForActionAtEnd
      type: boolean

  steps:
    - name: LaunchTestInstance
      action: LaunchInstance
      onFailure: Abort
      inputs:
        waitFor: "ssmAgent"

    - name: TerminateTestInstance
      action: TerminateInstance
      onFailure: Continue
      inputs:
        instanceId.$: "$.stepOutputs.LaunchTestInstance.instanceId"

    - name: WaitForActionAtEnd
      action: WaitForAction
      if:
        booleanEquals: true
        value: "$.parameters.waitForActionAtEnd"
  EOT
}
`, rName, description)
}

func testAccWorkflowConfig_kmsKeyID(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  deletion_window_in_days = 7
}

resource "aws_imagebuilder_workflow" "test" {
  name    = %[1]q
  version = "1.0.0"
  type    = "TEST"

  kms_key_id = aws_kms_key.test.arn

  data = <<-EOT
  name: test-image
  description: Workflow to test an image
  schemaVersion: 1.0

  parameters:
    - name: waitForActionAtEnd
      type: boolean

  steps:
    - name: LaunchTestInstance
      action: LaunchInstance
      onFailure: Abort
      inputs:
        waitFor: "ssmAgent"

    - name: TerminateTestInstance
      action: TerminateInstance
      onFailure: Continue
      inputs:
        instanceId.$: "$.stepOutputs.LaunchTestInstance.instanceId"

    - name: WaitForActionAtEnd
      action: WaitForAction
      if:
        booleanEquals: true
        value: "$.parameters.waitForActionAtEnd"
  EOT
}
`, rName)
}

func testAccWorkflowConfig_tags1(rName string, tagKey1 string, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_imagebuilder_workflow" "test" {
  name    = %[1]q
  version = "1.0.0"
  type    = "TEST"

  tags = {
    %[2]q = %[3]q
  }

  data = <<-EOT
  name: test-image
  description: Workflow to test an image
  schemaVersion: 1.0

  parameters:
    - name: waitForActionAtEnd
      type: boolean

  steps:
    - name: LaunchTestInstance
      action: LaunchInstance
      onFailure: Abort
      inputs:
        waitFor: "ssmAgent"

    - name: TerminateTestInstance
      action: TerminateInstance
      onFailure: Continue
      inputs:
        instanceId.$: "$.stepOutputs.LaunchTestInstance.instanceId"

    - name: WaitForActionAtEnd
      action: WaitForAction
      if:
        booleanEquals: true
        value: "$.parameters.waitForActionAtEnd"
  EOT
}
`, rName, tagKey1, tagValue1)
}

func testAccWorkflowConfig_tags2(rName string, tagKey1 string, tagValue1 string, tagKey2 string, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_imagebuilder_workflow" "test" {
  name    = %[1]q
  version = "1.0.0"
  type    = "TEST"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }

  data = <<-EOT
  name: test-image
  description: Workflow to test an image
  schemaVersion: 1.0

  parameters:
    - name: waitForActionAtEnd
      type: boolean

  steps:
    - name: LaunchTestInstance
      action: LaunchInstance
      onFailure: Abort
      inputs:
        waitFor: "ssmAgent"

    - name: TerminateTestInstance
      action: TerminateInstance
      onFailure: Continue
      inputs:
        instanceId.$: "$.stepOutputs.LaunchTestInstance.instanceId"

    - name: WaitForActionAtEnd
      action: WaitForAction
      if:
        booleanEquals: true
        value: "$.parameters.waitForActionAtEnd"
  EOT
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccWorkflowConfig_uri(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_object" "test" {
  bucket = aws_s3_bucket.test.bucket
  key    = "test.yml"

  content = <<-EOT
  name: test-image
  description: Workflow to test an image
  schemaVersion: 1.0

  parameters:
    - name: waitForActionAtEnd
      type: boolean

  steps:
    - name: LaunchTestInstance
      action: LaunchInstance
      onFailure: Abort
      inputs:
        waitFor: "ssmAgent"

    - name: TerminateTestInstance
      action: TerminateInstance
      onFailure: Continue
      inputs:
        instanceId.$: "$.stepOutputs.LaunchTestInstance.instanceId"

    - name: WaitForActionAtEnd
      action: WaitForAction
      if:
        booleanEquals: true
        value: "$.parameters.waitForActionAtEnd"
  EOT
}

resource "aws_imagebuilder_workflow" "test" {
  name    = %[1]q
  version = "1.0.0"
  type    = "TEST"

  uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.test.key}"
}
`, rName)
}
