// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticache_test

import (
	"context"
	"fmt"
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfelasticache "github.com/hashicorp/terraform-provider-aws/internal/service/elasticache"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccElastiCacheServerlessCache_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_elasticache_serverless_cache.test"
	var serverlessElasticCache awstypes.ServerlessCache

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ElastiCacheServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckServerlessCacheDestroy(ctx),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessCacheConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_usage_limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint.#"),
					resource.TestCheckResourceAttrSet(resourceName, "engine"),
					resource.TestCheckResourceAttrSet(resourceName, "full_engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "reader_endpoint.#"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ids.#"),
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

func TestAccElastiCacheServerlessCache_basicRedis(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_elasticache_serverless_cache.test"
	var serverlessElasticCache awstypes.ServerlessCache

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ElastiCacheServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckServerlessCacheDestroy(ctx),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessCacheConfig_basicRedis(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_usage_limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "daily_snapshot_time"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint.#"),
					resource.TestCheckResourceAttrSet(resourceName, "engine"),
					resource.TestCheckResourceAttrSet(resourceName, "full_engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "reader_endpoint.#"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ids.#"),
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

func TestAccElastiCacheServerlessCache_full(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_elasticache_serverless_cache.test"
	var serverlessElasticCache awstypes.ServerlessCache

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ElastiCacheServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckServerlessCacheDestroy(ctx),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessCacheConfig_full(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_usage_limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint.#"),
					resource.TestCheckResourceAttrSet(resourceName, "engine"),
					resource.TestCheckResourceAttrSet(resourceName, "full_engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "reader_endpoint.#"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ids.#"),
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

func TestAccElastiCacheServerlessCache_fullRedis(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_elasticache_serverless_cache.test"
	var serverlessElasticCache awstypes.ServerlessCache

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ElastiCacheServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckServerlessCacheDestroy(ctx),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessCacheConfig_fullRedis(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_usage_limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint.#"),
					resource.TestCheckResourceAttrSet(resourceName, "engine"),
					resource.TestCheckResourceAttrSet(resourceName, "full_engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "reader_endpoint.#"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ids.#"),
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

func TestAccElastiCacheServerlessCache_update(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	descriptionOld := "Memcached Serverless Cluster"
	descriptionNew := "Memcached Serverless Cluster updated"
	resourceName := "aws_elasticache_serverless_cache.test"
	var serverlessElasticCache awstypes.ServerlessCache

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ElastiCacheServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckServerlessCacheDestroy(ctx),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessCacheConfig_update(rName, descriptionOld),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_usage_limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionOld),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint.#"),
					resource.TestCheckResourceAttrSet(resourceName, "engine"),
					resource.TestCheckResourceAttrSet(resourceName, "full_engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "reader_endpoint.#"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ids.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServerlessCacheConfig_update(rName, descriptionNew),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_usage_limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionNew),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint.#"),
					resource.TestCheckResourceAttrSet(resourceName, "engine"),
					resource.TestCheckResourceAttrSet(resourceName, "full_engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "reader_endpoint.#"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ids.#"),
				),
			},
		},
	})
}

func TestAccElastiCacheServerlessCache_updatesc(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	descriptionOld := "Memcached Serverless Cluster"
	descriptionNew := "Memcached Serverless Cluster updated"
	resourceName := "aws_elasticache_serverless_cache.test"
	var serverlessElasticCache awstypes.ServerlessCache

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ElastiCacheServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckServerlessCacheDestroy(ctx),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessCacheConfig_updatesc(rName, descriptionOld, 1, 1000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_usage_limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionOld),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint.#"),
					resource.TestCheckResourceAttrSet(resourceName, "engine"),
					resource.TestCheckResourceAttrSet(resourceName, "full_engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "reader_endpoint.#"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ids.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServerlessCacheConfig_updatesc(rName, descriptionNew, 2, 1010),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "cache_usage_limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionNew),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint.#"),
					resource.TestCheckResourceAttrSet(resourceName, "engine"),
					resource.TestCheckResourceAttrSet(resourceName, "full_engine_version"),
					resource.TestCheckResourceAttrSet(resourceName, "reader_endpoint.#"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_ids.#"),
				),
			},
		},
	})
}

func TestAccElastiCacheServerlessCache_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_elasticache_serverless_cache.test"
	var serverlessElasticCache awstypes.ServerlessCache

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ElastiCacheServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckServerlessCacheDestroy(ctx),

		Steps: []resource.TestStep{
			{
				Config: testAccServerlessCacheConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfelasticache.ResourceServerlessCache, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccElastiCacheServerlessCache_tags(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_elasticache_serverless_cache.test"
	var serverlessElasticCache awstypes.ServerlessCache

	tags1 := `
  tags = {
    key1 = "value1"
  }
`
	tags2 := `
  tags = {
    key1 = "value1"
    key2 = "value2"
  }
`
	tags3 := `
  tags = {
    key2 = "value2"
  }
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ElastiCacheServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckServerlessCacheDestroy(ctx),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessCacheConfig_tags(rName, tags1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccServerlessCacheConfig_tags(rName, tags2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccServerlessCacheConfig_tags(rName, tags3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerlessCacheExists(ctx, resourceName, &serverlessElasticCache),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckServerlessCacheExists(ctx context.Context, resourceName string, v *awstypes.ServerlessCache) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ElastiCache Serverless ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ElastiCacheClient(ctx)
		out, err := tfelasticache.FindServerlessCacheByID(ctx, conn, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("retrieving ElastiCache Serverlesss (%s): %w", rs.Primary.ID, err)
		}

		*v = out

		return nil
	}
}

func testAccCheckServerlessCacheDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ElastiCacheClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_elasticache_serverless_cache" {
				continue
			}

			_, err := tfelasticache.FindServerlessCacheByID(ctx, conn, rs.Primary.ID)
			if tfresource.NotFound(err) {
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("ElastiCache Serverless (%s) still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccServerlessCacheConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_elasticache_serverless_cache" "test" {
  engine = "memcached"
  name   = %[1]q
}
`, rName)
}

func testAccServerlessCacheConfig_basicRedis(rName string) string {
	return fmt.Sprintf(`
resource "aws_elasticache_serverless_cache" "test" {
  engine = "redis"
  name   = %[1]q
}
`, rName)
}

func testAccServerlessCacheConfig_full(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
resource "aws_elasticache_serverless_cache" "test" {
  engine = "memcached"
  name   = %[1]q

  cache_usage_limits {
    data_storage {
      maximum = 10
      unit    = "GB"
    }
    ecpu_per_second {
      maximum = 1000
    }
  }

  description          = "Test Full Memcached Attributes"
  kms_key_id           = aws_kms_key.test.arn
  major_engine_version = "1.6"
  security_group_ids   = [aws_security_group.test.id]
  subnet_ids           = aws_subnet.test[*].id
  tags = {
    Name = %[1]q
  }
}

resource "aws_kms_key" "test" {
  description = "tf-test-cmk-kms-key-id"
}

resource "aws_security_group" "test" {
  name        = %[1]q
  description = %[1]q
  vpc_id      = aws_vpc.test.id

  ingress {
    from_port   = -1
    to_port     = -1
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
`, rName))
}

func testAccServerlessCacheConfig_fullRedis(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
resource "aws_elasticache_serverless_cache" "test" {
  engine = "redis"
  name   = %[1]q

  cache_usage_limits {
    data_storage {
      maximum = 10
      unit    = "GB"
    }
    ecpu_per_second {
      maximum = 1000
    }
  }

  daily_snapshot_time      = "09:00"
  description              = "Test Full Redis Attributes"
  kms_key_id               = aws_kms_key.test.arn
  major_engine_version     = "7"
  snapshot_retention_limit = 1
  security_group_ids       = [aws_security_group.test.id]
  subnet_ids               = aws_subnet.test[*].id
  tags = {
    Name = %[1]q
  }
}

resource "aws_kms_key" "test" {
  description = "tf-test-cmk-kms-key-id"
}

resource "aws_security_group" "test" {
  name        = %[1]q
  description = %[1]q
  vpc_id      = aws_vpc.test.id

  ingress {
    from_port   = -1
    to_port     = -1
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
`, rName))
}

func testAccServerlessCacheConfig_update(rName, desc string) string {
	return fmt.Sprintf(`
resource "aws_elasticache_serverless_cache" "test" {
  engine      = "memcached"
  name        = %[1]q
  description = %[2]q
}
`, rName, desc)
}

func testAccServerlessCacheConfig_updatesc(rName string, desc string, d1 int64, d2 int64) string {
	return fmt.Sprintf(`
resource "aws_elasticache_serverless_cache" "test" {
  engine      = "memcached"
  name        = %[1]q
  description = %[2]q
  cache_usage_limits {
    data_storage {
      maximum = %[3]d
      unit    = "GB"
    }
    ecpu_per_second {
      maximum = %[4]d
    }
  }
}
`, rName, desc, d1, d2)
}

func testAccServerlessCacheConfig_tags(rName, tags string) string {
	return fmt.Sprintf(`
resource "aws_elasticache_serverless_cache" "test" {
  engine = "memcached"
  name   = %[1]q

%[2]s
}
`, rName, tags)
}
