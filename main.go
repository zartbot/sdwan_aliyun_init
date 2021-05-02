package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/zartbot/sdwan_aliyun_init/alicloud"
	"github.com/zartbot/sdwan_aliyun_init/vmanageapi"
	"gopkg.in/yaml.v2"
)

var commandOptions = struct {
	VmanageCfgPath string
	RootCAPath     string
	ConfigPath     string
}{
	"vmanage.yml",
	"root-ca-chain.pem",
	"aliyun.csv",
}

func init() {
	flag.StringVar(&commandOptions.VmanageCfgPath, "vmanage", commandOptions.VmanageCfgPath, "Vmanage config file")
	flag.StringVar(&commandOptions.RootCAPath, "cert", commandOptions.RootCAPath, "CA file path")
	flag.StringVar(&commandOptions.ConfigPath, "config", commandOptions.ConfigPath, "Config file path")
	flag.Parse()
}

func main() {

	cfg := &vmanageapi.VmanageInfo{}

	yamlFile, err := ioutil.ReadFile(commandOptions.VmanageCfgPath)
	if err != nil {
		logrus.Fatalf("[Error]Config file fetch error, %v", err)
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		logrus.Fatalf("Unmarshal: %v", err)
	}

	vmanage := vmanageapi.NewVmanage(cfg.BaseURL, cfg.Username, cfg.Password)
	vbondip, _ := vmanage.GetvBond()
	orgName := vmanage.GetOrgnazation()
	uuid := vmanage.GetDevices()
	file, err := os.Open(commandOptions.ConfigPath)
	if err != nil {
		logrus.Fatal("[Aliyun C8000V Configfile] Error:", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		logrus.Fatal("[Aliyun C8000V Configfile] Error:", err)
	}

	//var result []*alicloud.CSRConfig
	for idx := 1; idx < len(records); idx++ {
		r := records[idx]
		vpnid, err := strconv.ParseInt(r[11], 10, 32)
		if err != nil {
			logrus.Fatal("Invalid VPNID", r[11])
		}

		siteid, err := strconv.ParseInt(r[12], 10, 32)
		if err != nil {
			logrus.Fatal("Invalid SiteID", r[12])
		}
		bw, err := strconv.ParseInt(r[16], 10, 32)
		if err != nil {
			logrus.Fatal("Invalid Bandwidth", r[16])
		}

		s := &alicloud.CSRConfig{
			RegionID:        r[0],
			ImageID:         r[1],
			ZoneID:          r[2],
			Latitude:        r[3],
			Longitude:       r[4],
			SecurityGroupID: r[5],
			WANVSwitchID:    r[6],
			WANPrivateIP:    r[7],
			LANVSwitchID:    r[8],
			LANPrivateIP:    r[9],
			LoopbackIP:      r[10],
			ServiceVPNID:    int(vpnid),
			SiteID:          int(siteid),
			SystemIP:        r[13],
			HostName:        r[14],
			Password:        r[15],
			WANBW:           int(bw),
			UUID:            uuid[idx-1].UUID,
			OTP:             uuid[idx-1].Token,
			OrgName:         orgName,
			VBond:           vbondip,
			RootCAPath:      commandOptions.RootCAPath,
			InstanceType:    "ecs.g5ne.large",
		}
		s.AliCloudCLi()
		fmt.Println()
	}

}
