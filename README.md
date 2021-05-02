# Cisco SDWAN Alicloud Init

This script is used to generate alicloud init command for Cisco 
Catalyst Edge8000v.

## Usage

```bash
Usage of ./sdwan-aliyun_linux:
  -cert string
        CA file path (default "root-ca-chain.pem")
  -config string
        Config file path (default "aliyun.csv")
  -vmanage string
        Vmanage config file (default "vmanage.yml")
```

Config Alicloud VPC and collect the related information then store it in csv(`aliyun.csv`) format as below

```
regionId,imageID,zoneId,lat,long,securitygroup,wanSwitch,WANIP,lanSwitch,LANIp,Lo,serviceVPN,siteid,siteip,hostname,loginpassword,bw
cn-shanghai,m-uf6hwh6xxxxx,cn-shanghai-l,31.232,121.469,sg-123456,vsw-123456,10.116.0.1,vsw-123456,10.116.1.1,10.99.1.116,99,116001,10.116.0.1,ali_shanghai1,adminpwdhaha,1
```
Config vmanage access info in the following yaml(`vmanage.yml`) file

```yaml
baseurl: <vmanageIP:port>
username: admin
password: AdminPassword
```

then execute the command:

```bash
./sdwan-aliyun_linux -cert=root-ca-chain.pem -config=aliyun.csv -vmanage=vmanage.yml
```

finally execute the output in shell or aliyun cloud shell.




