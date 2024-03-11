// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dms_test

import (
	"context"
	"fmt"
	"testing"

	dms "github.com/aws/aws-sdk-go/service/databasemigrationservice"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdms "github.com/hashicorp/terraform-provider-aws/internal/service/dms"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccDMSEventSubscription_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var eventSubscription dms.EventSubscription
	resourceName := "aws_dms_event_subscription.test"
	snsTopicResourceName := "aws_sns_topic.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEventSubscriptionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionConfig_enabled(rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "event_categories.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "event_categories.*", "creation"),
					resource.TestCheckTypeSetElemAttr(resourceName, "event_categories.*", "failure"),
					resource.TestCheckResourceAttrPair(resourceName, "sns_topic_arn", snsTopicResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "source_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "replication-instance"),
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

func TestAccDMSEventSubscription_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var eventSubscription dms.EventSubscription
	resourceName := "aws_dms_event_subscription.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEventSubscriptionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionConfig_enabled(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdms.ResourceEventSubscription(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDMSEventSubscription_enabled(t *testing.T) {
	ctx := acctest.Context(t)
	var eventSubscription dms.EventSubscription
	resourceName := "aws_dms_event_subscription.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEventSubscriptionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionConfig_enabled(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventSubscriptionConfig_enabled(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccEventSubscriptionConfig_enabled(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccDMSEventSubscription_eventCategories(t *testing.T) {
	ctx := acctest.Context(t)
	var eventSubscription dms.EventSubscription
	resourceName := "aws_dms_event_subscription.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEventSubscriptionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionConfig_eventCategories(rName, "creation", "failure"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
					resource.TestCheckResourceAttr(resourceName, "event_categories.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "event_categories.*", "creation"),
					resource.TestCheckTypeSetElemAttr(resourceName, "event_categories.*", "failure"),
					resource.TestCheckResourceAttr(resourceName, "source_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventSubscriptionConfig_eventCategories(rName, "configuration change", "deletion"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
					resource.TestCheckResourceAttr(resourceName, "event_categories.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "event_categories.*", "configuration change"),
					resource.TestCheckTypeSetElemAttr(resourceName, "event_categories.*", "deletion"),
					resource.TestCheckResourceAttr(resourceName, "source_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccDMSEventSubscription_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var eventSubscription dms.EventSubscription
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_dms_event_subscription.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DMSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEventSubscriptionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscriptionConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
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
				Config: testAccEventSubscriptionConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccEventSubscriptionConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEventSubscriptionExists(ctx, resourceName, &eventSubscription),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckEventSubscriptionDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_dms_event_subscription" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).DMSConn(ctx)

			_, err := tfdms.FindEventSubscriptionByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DMS Event Subscription %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckEventSubscriptionExists(ctx context.Context, n string, v *dms.EventSubscription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DMSConn(ctx)

		output, err := tfdms.FindEventSubscriptionByName(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccEventSubscriptionConfig_base(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_dms_replication_subnet_group" "test" {
  replication_subnet_group_description = %[1]q
  replication_subnet_group_id          = %[1]q
  subnet_ids                           = aws_subnet.test[*].id
}

resource "aws_dms_replication_instance" "test" {
  apply_immediately           = true
  replication_instance_class  = data.aws_partition.current.partition == "aws" ? "dms.t2.micro" : "dms.c4.large"
  replication_instance_id     = %[1]q
  replication_subnet_group_id = aws_dms_replication_subnet_group.test.replication_subnet_group_id
}

resource "aws_sns_topic" "test" {
  name = %[1]q
}
`, rName))
}

func testAccEventSubscriptionConfig_enabled(rName string, enabled bool) string {
	return acctest.ConfigCompose(testAccEventSubscriptionConfig_base(rName), fmt.Sprintf(`
resource "aws_dms_event_subscription" "test" {
  name             = %[1]q
  enabled          = %[2]t
  event_categories = ["creation", "failure"]
  source_type      = "replication-instance"
  sns_topic_arn    = aws_sns_topic.test.arn
}
`, rName, enabled))
}

func testAccEventSubscriptionConfig_eventCategories(rName string, eventCategory1, eventCategory2 string) string {
	return acctest.ConfigCompose(testAccEventSubscriptionConfig_base(rName), fmt.Sprintf(`
resource "aws_dms_event_subscription" "test" {
  name             = %[1]q
  enabled          = false
  event_categories = [%[2]q, %[3]q]
  source_type      = "replication-instance"
  source_ids       = [aws_dms_replication_instance.test.replication_instance_id]
  sns_topic_arn    = aws_sns_topic.test.arn
}
`, rName, eventCategory1, eventCategory2))
}

func testAccEventSubscriptionConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccEventSubscriptionConfig_base(rName), fmt.Sprintf(`
resource "aws_dms_event_subscription" "test" {
  name             = %[1]q
  enabled          = true
  event_categories = ["creation", "failure"]
  source_type      = "replication-instance"
  source_ids       = [aws_dms_replication_instance.test.replication_instance_id]
  sns_topic_arn    = aws_sns_topic.test.arn

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccEventSubscriptionConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccEventSubscriptionConfig_base(rName), fmt.Sprintf(`
resource "aws_dms_event_subscription" "test" {
  name             = %[1]q
  enabled          = true
  event_categories = ["creation", "failure"]
  source_type      = "replication-instance"
  source_ids       = [aws_dms_replication_instance.test.replication_instance_id]
  sns_topic_arn    = aws_sns_topic.test.arn

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}
