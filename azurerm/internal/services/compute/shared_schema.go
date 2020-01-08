package compute

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SourceImageReferenceSchema() *schema.Schema {
	// whilst originally I was hoping we could use the 'id' from `azurerm_platform_image' unfortunately Azure doesn't
	// like this as a value for the 'id' field:
	// Id /...../Versions/16.04.201909091 is not a valid resource reference."
	// as such the image is split into two fields (source_image_id and source_image_reference) to provide better validation
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"source_image_id"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"publisher": {
					Type:     schema.TypeString,
					Required: true,
				},
				"offer": {
					Type:     schema.TypeString,
					Required: true,
				},
				"sku": {
					Type:     schema.TypeString,
					Required: true,
				},
				"version": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func ExpandSourceImageReference(referenceInput []interface{}, imageId string) (*compute.ImageReference, error) {
	if imageId != "" {
		return &compute.ImageReference{
			ID: utils.String(imageId),
		}, nil
	}

	if len(referenceInput) == 0 {
		return nil, fmt.Errorf("Either a `source_image_id` or a `source_image_reference` block must be specified!")
	}

	raw := referenceInput[0].(map[string]interface{})
	return &compute.ImageReference{
		Publisher: utils.String(raw["publisher"].(string)),
		Offer:     utils.String(raw["offer"].(string)),
		Sku:       utils.String(raw["sku"].(string)),
		Version:   utils.String(raw["version"].(string)),
	}, nil
}

func FlattenSourceImageReference(input *compute.ImageReference) []interface{} {
	// since the image id is pulled out as a separate field, if that's set we should return an empty block here
	if input == nil || input.ID != nil {
		return []interface{}{}
	}

	var publisher, offer, sku, version string

	if input.Publisher != nil {
		publisher = *input.Publisher
	}
	if input.Offer != nil {
		offer = *input.Offer
	}
	if input.Sku != nil {
		sku = *input.Sku
	}
	if input.Version != nil {
		version = *input.Version
	}

	return []interface{}{
		map[string]interface{}{
			"publisher": publisher,
			"offer":     offer,
			"sku":       sku,
			"version":   version,
		},
	}
}
