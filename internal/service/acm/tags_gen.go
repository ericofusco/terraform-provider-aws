// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package acm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// ListTags lists acm service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(ctx context.Context, conn acmiface.ACMAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &acm.ListTagsForCertificateInput{
		CertificateArn: aws.String(identifier),
	}

	output, err := conn.ListTagsForCertificateWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns acm service tags.
func Tags(tags tftags.KeyValueTags) []*acm.Tag {
	result := make([]*acm.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &acm.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from acm service tags.
func KeyValueTags(ctx context.Context, tags []*acm.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// UpdateTags updates acm service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(ctx context.Context, conn acmiface.ACMAPI, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &acm.RemoveTagsFromCertificateInput{
			CertificateArn: aws.String(identifier),
			Tags:           Tags(removedTags.IgnoreAWS()),
		}

		_, err := conn.RemoveTagsFromCertificateWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &acm.AddTagsToCertificateInput{
			CertificateArn: aws.String(identifier),
			Tags:           Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.AddTagsToCertificateWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags interface{}, newTags interface{}) error {
	return UpdateTags(ctx, meta.(*conns.AWSClient).ACMConn(), identifier, oldTags, newTags)
}
