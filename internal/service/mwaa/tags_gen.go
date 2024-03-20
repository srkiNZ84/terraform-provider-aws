// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package mwaa

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mwaa"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types/option"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// map[string]string handling

// Tags returns mwaa service tags.
func Tags(tags tftags.KeyValueTags) map[string]string {
	return tags.Map()
}

// KeyValueTags creates tftags.KeyValueTags from mwaa service tags.
func KeyValueTags(ctx context.Context, tags map[string]string) tftags.KeyValueTags {
	return tftags.New(ctx, tags)
}

// getTagsIn returns mwaa service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) map[string]string {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets mwaa service tags in Context.
func setTagsOut(ctx context.Context, tags map[string]string) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(KeyValueTags(ctx, tags))
	}
}

// updateTags updates mwaa service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn *mwaa.Client, identifier string, oldTagsMap, newTagsMap any, optFns ...func(*mwaa.Options)) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.MWAA)
	if len(removedTags) > 0 {
		input := &mwaa.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			TagKeys:     removedTags.Keys(),
		}

		_, err := conn.UntagResource(ctx, input, optFns...)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.MWAA)
	if len(updatedTags) > 0 {
		input := &mwaa.TagResourceInput{
			ResourceArn: aws.String(identifier),
			Tags:        Tags(updatedTags),
		}

		_, err := conn.TagResource(ctx, input, optFns...)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

// UpdateTags updates mwaa service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).MWAAClient(ctx), identifier, oldTags, newTags)
}