/**
 * Copyright (c) 2016 Magicshui
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */
/**
 * Copyright (c) 2017 yunify
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceTagColor = "color"
)

func resourceQingcloudTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudTagCreate,
		Read:   resourceQingcloudTagRead,
		Update: resourceQingcloudTagUpdate,
		Delete: resourceQingcloudTagDelete,
		Schema: map[string]*schema.Schema{
			resourceName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			resourceTagColor: &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      DEFAULT_TAG_COLOR,
				ValidateFunc: validateColorString,
			},
			resourceDescription: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQingcloudTagCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.CreateTagInput)
	input.TagName, _ = getNamePointer(d)
	var output *qc.CreateTagOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.CreateTag(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.TagID))
	return resourceQingcloudTagUpdate(d, meta)
}
func resourceQingcloudTagRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.DescribeTagsInput)
	input.Tags = []*string{qc.String(d.Id())}
	var output *qc.DescribeTagsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeTags(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.TagSet) == 0 {
		d.SetId("")
		return nil
	}
	tag := output.TagSet[0]
	d.Set(resourceName, qc.StringValue(tag.TagName))
	d.Set(resourceDescription, qc.StringValue(tag.Description))
	if qc.StringValue(tag.Color) == "default" {
		d.Set(resourceTagColor, DEFAULT_TAG_COLOR)
	} else {
		d.Set(resourceTagColor, qc.StringValue(tag.Color))
	}
	return nil
}
func resourceQingcloudTagUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := modifyTagAttributes(d, meta); err != nil {
		return err
	}
	return nil
}
func resourceQingcloudTagDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).tag
	input := new(qc.DeleteTagsInput)
	input.Tags = []*string{qc.String(d.Id())}
	var output *qc.DeleteTagsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DeleteTags(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
