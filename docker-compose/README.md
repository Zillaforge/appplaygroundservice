# APS Setup


## Replace OpenStack Identity Endpoint

```shell
APS_YAML="docker-compose/etc/aps.yaml"
VPS_YAML="docker-compose/etc/vps.yaml"
VRM_YAML="docker-compose/etc/vrm.yaml"


KEYSTONE_ENDPOINT="http://140.110.160.230:5000/v3"

# 替換 keystone endpoint
sed -i.bak "s/REPLACE_KEYSTONE_ENDPOINT/$KEYSTONE_ENDPOINT/"     "$APS_YAML" && rm "$APS_YAML.bak"
sed -i.bak "s/REPLACE_KEYSTONE_ENDPOINT/$KEYSTONE_ENDPOINT/"     "$VPS_YAML" && rm "$VPS_YAML.bak"
sed -i.bak "s/REPLACE_KEYSTONE_ENDPOINT/$KEYSTONE_ENDPOINT/"     "$VRM_YAML" && rm "$VRM_YAML.bak"

```