package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/router"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func applyRouterUpdates(meta interface{}, routerID string) error {
	clt := meta.(*QingCloudClient).router
	params := router.UpdateRoutersRequest{}
	params.RoutersN.Add(routerID)
	if _, err := clt.UpdateRouters(params); err != nil {
		return err
	}
	if _, err := RouterTransitionStateRefresh(clt, routerID); err != nil {
		return err
	}
	return nil
}

func modifyRouterAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.ModifyRouterAttributesInput)
	input.Routers = []*string{qc.String(d.Id())}

	if create {
		if description := d.Get("description").(string); description != "" {
			input.Description = qc.String(description)
		}
	} else {
		if d.HasChange("description") {
			input.Description = qc.String(d.Get("description").(string))
		}
		if d.HasChange("name") {
			input.Description = qc.String(d.Get("name").(string))
		}
		if d.HasChange("eip") {
			input.EIP = qc.String(d.Get("eip").(string))
		}
		if d.HasChange("securitygroup") {
			input.SecurityGroup = qc.String(d.Get("eip").(string))
		}
	}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error modify router attributes input validate: %s", err)
	}
	output, err := clt.ModifyRouterAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify router attributes: %s", err)
	}
	if output.RetCode != 0 {
		return fmt.Errorf("Error modify router attrubites: %s", output.Message)
	}
	return nil
}

// func modifyRouterVxnets(d *schema.ResourceData, meta interface{}, create bool) error {
// 	clt := meta.(*QingCloudClient).router
// 	if create {
// 		map
// 	} else {
// 		if
// 	}
// }

// func getEIPInfoMap(data *qc.EIP) map[string]interface{} {
// 	var a = make(map[string]interface{}, 3)
// 	a["eip_id"] = qc.EIP.EIPID
// 	a["eip_name"] = qc.EIP.EIPName
// 	a["eip_addr"] = qc.EIP.EIPAddr
// 	return a
// }

// func getVxnetsMap(data []*qc.VxNet) map[string]interface{} {
// 	length := len(data)
// 	if data > 0 {
// 		var a = make(map[string]interface{}, length)
// 		for _, vxnet := range data {
// 			a[vxnet.VxNetID] = vxnet.NICID
// 		}
// 		return a
// 	}
// 	return nil
// }
