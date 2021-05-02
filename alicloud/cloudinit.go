package alicloud

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func (cfg *CSRConfig) CloudinitData() string {

	var buf bytes.Buffer
	/*
		outputFileName := "cloud_init_" + cfg.HostName + ".txt"
		f, err := os.Create(outputFileName)
		if err != nil {
			logrus.Fatal("Create cloud-init file failed:  ", err)
		}
		defer f.Close()
	*/
	cert, err := os.Open(cfg.RootCAPath)
	if err != nil {
		logrus.Fatal("Open CA file error  :  ", err)
	}
	defer cert.Close()

	fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=\"===============3067523750048488884==\"\n")
	fmt.Fprintf(&buf, "MIME-Version: 1.0\n\n")
	fmt.Fprintf(&buf, "--===============3067523750048488884==\n")
	fmt.Fprintf(&buf, "Content-Type: text/cloud-config; charset=\"us-ascii\"\n")
	fmt.Fprintf(&buf, "MIME-Version: 1.0\n")
	fmt.Fprintf(&buf, "Content-Transfer-Encoding: 7bit\n")
	fmt.Fprintf(&buf, "Content-Disposition: attachment; filename=\"cloud-config\"\n\n")

	fmt.Fprintf(&buf, "#cloud-config\n")
	fmt.Fprintf(&buf, "vinitparam:\n")
	fmt.Fprintf(&buf, " - otp : %s\n", cfg.OTP)
	fmt.Fprintf(&buf, " - vbond : %s\n", cfg.VBond)
	fmt.Fprintf(&buf, " - uuid : %s\n", cfg.UUID)
	fmt.Fprintf(&buf, " - org : %s\n", cfg.OrgName)
	fmt.Fprintf(&buf, " - rcc : true\n")
	fmt.Fprintf(&buf, "ca-certs:\n")
	fmt.Fprintf(&buf, " remove-defaults: false\n")
	fmt.Fprintf(&buf, " trusted:\n")
	fmt.Fprintf(&buf, " - |\n")

	reader := bufio.NewReader(cert)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		fmt.Fprintf(&buf, "   %s", line)
		if err != nil {
			break
		}
	}

	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "--===============3067523750048488884==\n")
	fmt.Fprintf(&buf, "Content-Type: text/cloud-boothook; charset=\"us-ascii\"\n")
	fmt.Fprintf(&buf, "MIME-Version: 1.0\n")
	fmt.Fprintf(&buf, "Content-Transfer-Encoding: 7bit\n")
	fmt.Fprintf(&buf, "Content-Disposition: attachment;\n")
	fmt.Fprintf(&buf, " filename=\"config-%s.txt\"\n", cfg.UUID)

	fmt.Fprintf(&buf, "#cloud-boothook\n")
	fmt.Fprintf(&buf, "  system\n")
	fmt.Fprintf(&buf, "   ztp-status            success\n")
	fmt.Fprintf(&buf, "   pseudo-confirm-commit 300\n")
	fmt.Fprintf(&buf, "   personality           vedge\n")
	fmt.Fprintf(&buf, "   device-model          vedge-C8000V\n")
	fmt.Fprintf(&buf, "   chassis-number        %s\n", cfg.UUID)
	fmt.Fprintf(&buf, "   host-name             %s\n", cfg.HostName)
	fmt.Fprintf(&buf, "   system-ip             %s\n", cfg.SystemIP)
	fmt.Fprintf(&buf, "   overlay-id            1\n")
	fmt.Fprintf(&buf, "   site-id               %d\n", cfg.SiteID)
	fmt.Fprintf(&buf, "   gps-location latitude %s\n", cfg.Latitude)
	fmt.Fprintf(&buf, "   gps-location longitude %s\n", cfg.Longitude)
	fmt.Fprintf(&buf, "   no port-offset\n")
	fmt.Fprintf(&buf, "   control-session-pps   300\n")
	fmt.Fprintf(&buf, "   admin-tech-on-failure\n")
	fmt.Fprintf(&buf, "   sp-organization-name  \"%s\"\n", cfg.OrgName)
	fmt.Fprintf(&buf, "   organization-name     \"%s\"\n", cfg.OrgName)
	fmt.Fprintf(&buf, "   vbond %s  port 12346\n", cfg.VBond)
	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "\n")

	fmt.Fprintf(&buf, "  vrf definition %d\n", cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "   rd 1:%d\n", cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "   address-family ipv4\n")
	fmt.Fprintf(&buf, "    exit-address-family\n")
	fmt.Fprintf(&buf, "   !\n")
	fmt.Fprintf(&buf, "  !\n")

	fmt.Fprintf(&buf, "  ip pim vrf %d autorp listener\n", cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "  ip pim vrf %d send-rp-announce Loopback%d scope 12\n", cfg.ServiceVPNID, cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "  ip pim vrf %d send-rp-discovery Loopback%d scope 12\n", cfg.ServiceVPNID, cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "  ip pim vrf %d ssm default\n", cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "  !\n")

	fmt.Fprintf(&buf, "  ip multicast-routing vrf %d distributed\n", cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "  !\n")
	fmt.Fprintf(&buf, "  sdwan\n")
	fmt.Fprintf(&buf, "   interface GigabitEthernet1\n")
	fmt.Fprintf(&buf, "    tunnel-interface\n")
	fmt.Fprintf(&buf, "     encapsulation ipsec weight 1\n")
	fmt.Fprintf(&buf, "     color default\n")
	fmt.Fprintf(&buf, "     allow-service all\n")
	fmt.Fprintf(&buf, "    exit\n")
	fmt.Fprintf(&buf, "   exit\n")
	fmt.Fprintf(&buf, "   multicast\n")
	fmt.Fprintf(&buf, "    address-family ipv4 vrf %d\n", cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "     replicator threshold 7500\n")

	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "  hostname %s\n", cfg.HostName)
	fmt.Fprintf(&buf, "  username admin privilege 15 secret %s\n", cfg.Password)

	fmt.Fprintf(&buf, "  interface GigabitEthernet1\n")
	fmt.Fprintf(&buf, "   no shutdown\n")
	fmt.Fprintf(&buf, "   arp timeout 1200\n")
	fmt.Fprintf(&buf, "   ip address dhcp client-id GigabitEthernet1\n")
	fmt.Fprintf(&buf, "   no ip redirects\n")
	fmt.Fprintf(&buf, "   ip dhcp client default-router distance 1\n")
	fmt.Fprintf(&buf, "   ip mtu    1500\n")
	fmt.Fprintf(&buf, "   load-interval 30\n")
	fmt.Fprintf(&buf, "   mtu         1500\n")
	fmt.Fprintf(&buf, "   negotiation auto\n")
	fmt.Fprintf(&buf, "  exit\n")
	fmt.Fprintf(&buf, "  interface Tunnel1\n")
	fmt.Fprintf(&buf, "   no shutdown\n")
	fmt.Fprintf(&buf, "   ip unnumbered GigabitEthernet1\n")
	fmt.Fprintf(&buf, "   no ip redirects\n")
	fmt.Fprintf(&buf, "   ipv6 unnumbered GigabitEthernet1\n")
	fmt.Fprintf(&buf, "   no ipv6 redirects\n")
	fmt.Fprintf(&buf, "   tunnel source GigabitEthernet1\n")
	fmt.Fprintf(&buf, "   tunnel mode sdwan\n")
	fmt.Fprintf(&buf, "  exit\n")

	fmt.Fprintf(&buf, "  interface GigabitEthernet2\n")
	fmt.Fprintf(&buf, "   no shutdown\n")
	fmt.Fprintf(&buf, "   vrf forwarding %d\n", cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "   arp timeout 1200\n")
	fmt.Fprintf(&buf, "   ip address dhcp client-id GigabitEthernet2\n")
	fmt.Fprintf(&buf, "   no ip redirects\n")
	fmt.Fprintf(&buf, "   ip dhcp client default-router distance 1\n")
	fmt.Fprintf(&buf, "   ip mtu    1500\n")
	fmt.Fprintf(&buf, "   load-interval 30\n")
	fmt.Fprintf(&buf, "   ip pim sparse-mode\n")
	fmt.Fprintf(&buf, "   ip igmp version 3\n")
	fmt.Fprintf(&buf, "   mtu         1500\n")
	fmt.Fprintf(&buf, "   negotiation auto\n")
	fmt.Fprintf(&buf, "  exit\n")

	fmt.Fprintf(&buf, "  interface Loopback%d\n", cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "   no shutdown\n")
	fmt.Fprintf(&buf, "   vrf forwarding %d\n", cfg.ServiceVPNID)
	fmt.Fprintf(&buf, "   ip address %s 255.255.255.255\n", cfg.LoopbackIP)
	fmt.Fprintf(&buf, "   ip pim sparse-mode\n")
	fmt.Fprintf(&buf, "   ip igmp version 3\n")
	fmt.Fprintf(&buf, "  exit\n")

	fmt.Fprintf(&buf, "  line vty 0 4\n")
	fmt.Fprintf(&buf, "   login authentication default\n")
	fmt.Fprintf(&buf, "   transport input ssh\n")
	fmt.Fprintf(&buf, "  !\n")
	fmt.Fprintf(&buf, "  line vty 5 80\n")
	fmt.Fprintf(&buf, "   transport input ssh\n")
	fmt.Fprintf(&buf, "  !\n")

	fmt.Fprintf(&buf, "\n--===============3067523750048488884==--\n")

	//DEBUG_ME: fmt.Println(buf.String())
	return base64.StdEncoding.EncodeToString(buf.Bytes())

}
