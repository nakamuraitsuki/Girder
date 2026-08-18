package ovn

import "github.com/ovn-kubernetes/libovsdb/model"

func DatabaseModel() (model.ClientDBModel, error) {
	return model.NewClientDBModel("OVN_Northbound", map[string]model.Model{
		"Logical_Switch":      &LogicalSwitch{},
		"Logical_Switch_Port": &LogicalSwitchPort{},
		"Logical_Router":      &LogicalRouter{},
		"Logical_Router_Port": &LogicalRouterPort{},
		"Logical_Router_Static_Route": &LogicalRouterStaticRoute{},
	})
}
