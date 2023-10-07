// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudfront_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcloudfront "github.com/hashicorp/terraform-provider-aws/internal/service/cloudfront"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccCloudFrontPublicKey_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudfront_public_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPublicKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicKeyConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPublicKeyExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "caller_reference"),
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
					resource.TestCheckResourceAttrSet(resourceName, "encoded_key"),
					resource.TestCheckResourceAttrSet(resourceName, "etag"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func TestAccCloudFrontPublicKey_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudfront_public_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPublicKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicKeyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicKeyExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfcloudfront.ResourcePublicKey(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCloudFrontPublicKey_namePrefix(t *testing.T) {
	ctx := acctest.Context(t)
	startsWithPrefix := regexache.MustCompile("^tf-acc-test-")
	resourceName := "aws_cloudfront_public_key.example"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPublicKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicKeyConfig_namePrefix(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicKeyExists(ctx, resourceName),
					resource.TestMatchResourceAttr("aws_cloudfront_public_key.example", "name", startsWithPrefix),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name_prefix",
				},
			},
		},
	})
}

func TestAccCloudFrontPublicKey_update(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudfront_public_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudfront.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPublicKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicKeyConfig_comment(rName, "comment 1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicKeyExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", "comment 1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPublicKeyConfig_comment(rName, "comment 2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPublicKeyExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", "comment 2"),
				),
			},
		},
	})
}

func testAccCheckPublicKeyExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFrontConn(ctx)

		_, err := tfcloudfront.FindPublicKeyByID(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckPublicKeyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFrontConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_cloudfront_public_key" {
				continue
			}

			_, err := tfcloudfront.FindPublicKeyByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("CloudFront Public Key %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccPublicKeyConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudfront_public_key" "test" {
  encoded_key = file("test-fixtures/cloudfront-public-key.pem")
  name        = %[1]q
}
`, rName)
}

func testAccPublicKeyConfig_namePrefix() string {
	return `
resource "aws_cloudfront_public_key" "example" {
  comment     = "test key"
  encoded_key = file("test-fixtures/cloudfront-public-key.pem")
  name_prefix = "tf-acc-test-"
}
`
}

func testAccPublicKeyConfig_comment(rName, comment string) string {
	return fmt.Sprintf(`
resource "aws_cloudfront_public_key" "test" {
  comment     = %[2]q
  encoded_key = file("test-fixtures/cloudfront-public-key.pem")
  name        = %[1]q
}
`, rName, comment)
}
