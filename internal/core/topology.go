package core

func RouterPortName(switchOrGatewayName string) string {
	return "to-" + switchOrGatewayName
}

func LogicalSwitchPortName(routerName, switchOrGatewayName string) string {
	return routerName + "-to-" + switchOrGatewayName
}