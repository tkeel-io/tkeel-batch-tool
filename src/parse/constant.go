package parse

const (
	//Extension Fields								  //
	ExtensionKey_PointTag        = "extp-tag"         //{true |string|测点tag}
	ExtensionKey_PointName       = "extp-name"        //{true |string|测点tag}
	ExtensionKey_PointPointType  = "extp-type"        //{true |string|测点tag} enum{"AI"|"AO"|"DI"|"DO"|"Alarm"}
	ExtensionKey_PointStatus     = "extp-status"      //{false|string|测点tag}
	ExtensionKey_PointUnit       = "extp-unit"        //{false|string|测点tag}
	ExtensionKey_PointAlarmLevel = "extp-alarm-level" //{false|integer|测点tag}
	ExtensionKey_PointAlarmType  = "extp-alarm-type"  //{true |integer|测点tag}
	ExtensionKey_PointPeriod     = "extp-period"      //{true |integer|测点tag}
	ExtensionKey_PointPercentage = "extp-percetage"   //{false|floag |测点tag}
	ExtensionKey_PointAbsValue   = "extp-abs-val"     //{false|floag |测点tag}
	ExtensionKey_PointAoBound    = "extp-ao-bound"    //{false|string|测点tag}
	ExtensionKey_PointDeviceTag  = "extp-device-tag"  //{true |string|测点tag}
	ExtensionKey_PointDataType   = "extp-data-type"   //{false|string|测点tag}
	ExtensionKey_PointDataScope  = "extp-data-scope"  //{false|string|测点tag}
	ExtensionKey_PointPrecision  = "extp-precision"   //{false|float |测点tag}
	ExtensionKey_PointDesc       = "extp-desc"        //{false|string|测点描述}
	ExtensionKey_PointAlias      = "extp-alias"       //{false|string|测点别名}
	ExtensionKey_Point           = "extp-"            //{false|string|测点tag}//保留

	//device Extension Fields
	ExtensionKey_DeviceTag          = "extd-tag"           //true|string
	ExtensionKey_DevicePath         = "extd-path"          //false|string
	ExtensionKey_DeviceName         = "extd-name"          //true|string
	ExtensionKey_DeviceEdgeId       = "extd-edge-id"       //true|string
	ExtensionKey_DeviceSpaceGuid    = "extd-space-guid"    //false|string
	ExtensionKey_DeviceType         = "extd-type"          //false|string
	ExtensionKey_DeviceDesc         = "extd-desc"          //false|string
	ExtensionKey_DeviceProtocol     = "extd-protocol"      //false|string
	ExtensionKey_DeviceRegisterType = "extd-register-type" //false|string
	ExtensionKey_DeviceDriverId     = "extd-driverid"      //false|string
	ExtensionKey_DeviceGatewayId    = "extd-gwid"          //false|string
	ExtensionKey_DeviceAlias        = "extd-alias"         //false|string
	ExtensionKey_Device             = "extd-"              //false|string//保留
)
