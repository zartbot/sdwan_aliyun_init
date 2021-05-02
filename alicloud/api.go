package alicloud

import (
	"fmt"
)

type CSRConfig struct {
	VBond           string
	UUID            string
	OTP             string
	OrgName         string
	RootCAPath      string
	Password        string
	SiteID          int
	Latitude        string
	Longitude       string
	SystemIP        string
	HostName        string
	ServiceVPNID    int
	InstanceType    string
	RegionID        string
	ZoneID          string
	ImageID         string
	SecurityGroupID string
	WANVSwitchID    string
	WANPrivateIP    string
	LANVSwitchID    string
	LANPrivateIP    string
	LoopbackIP      string
	WANBW           int
}

func (cfg *CSRConfig) AliCloudCLi() string {
	fmt.Printf("aliyun ecs RunInstances \\\n")
	fmt.Printf("--RegionId %s \\\n", cfg.RegionID)
	fmt.Printf("--ImageId %s \\\n", cfg.ImageID)
	fmt.Printf("--InstanceType %s \\\n", cfg.InstanceType)
	fmt.Printf("--ZoneId %s \\\n", cfg.ZoneID)
	fmt.Printf("--InstanceName %s \\\n", cfg.HostName)
	fmt.Printf("--SecurityGroupId %s \\\n", cfg.SecurityGroupID)
	fmt.Printf("--VSwitchId %s \\\n", cfg.WANVSwitchID)
	fmt.Printf("--PrivateIpAddress %s \\\n", cfg.WANPrivateIP)
	fmt.Printf("--NetworkInterface.1.VSwitchId  %s \\\n", cfg.LANVSwitchID)
	fmt.Printf("--NetworkInterface.1.PrimaryIpAddress  %s \\\n", cfg.LANPrivateIP)
	fmt.Printf("--UserData %s \\\n", cfg.CloudinitData())
	fmt.Printf("--InstanceChargeType PostPaid \\\n")
	fmt.Printf("--InternetMaxBandwidthOut %d\n", cfg.WANBW)
	return ""
}
