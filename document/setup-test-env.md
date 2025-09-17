# Setup Testing OpenStack Environment

## Install Tools

```bash
# Install Docker
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

```


```bash
# Install K3s
curl -sfL https://get.k3s.io | sh -
sudo chmod 644 /etc/rancher/k3s/k3s.yaml
echo 'export KUBECONFIG=/etc/rancher/k3s/k3s.yaml' >> ~/.bashrc
echo 'source <(kubectl completion bash)' >> ~/.bashrc


# Install Helm
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
sudo chmod 700 get_helm.sh
./get_helm.sh
rm get_helm.sh
```

### Prepare LDAP Service and PegasusIAM Image

```bash

cd ~
git clone https://github.com/Zillaforge/utility-image-builder
cd utility-image-builder
sudo make release-image-debugger
sudo make release-image-golang


repos=("ldapservice" "pegasusiam" "eventpublishplugin")

for repo in "${repos[@]}"; do
    cd ~
    git clone "https://github.com/Zillaforge/$repo"
    cd "$repo"
    sudo make RELEASE_MODE=prod release-image
done

cd ~
sudo docker save -o ldapservice.tar Zillaforge/ldapservice:0.0.5
sudo docker save -o pegasusiam.tar Zillaforge/iam:1.8.3
sudo docker save -o event-publish-plugin.tar Zillaforge/event-publish-plugin:0.1.2
sudo docker save -o debugger.tar Zillaforge/debugger:1.0.2

sudo ctr --address /run/k3s/containerd/containerd.sock -n k8s.io images import ~/ldapservice.tar
sudo ctr --address /run/k3s/containerd/containerd.sock -n k8s.io images import ~/pegasusiam.tar
sudo ctr --address /run/k3s/containerd/containerd.sock -n k8s.io images import ~/event-publish-plugin.tar
sudo ctr --address /run/k3s/containerd/containerd.sock -n k8s.io images import ~/debugger.tar

sudo ctr --address /run/k3s/containerd/containerd.sock -n k8s.io images ls -q | while read -r img; do
    echo "Labeling $img ..."
    sudo ctr --address /run/k3s/containerd/containerd.sock -n k8s.io image label "$img" io.cri-containerd.image=managed
done

```

### Install Testing OpenStack Environment

```bash
cd ~
git clone https://github.com/Zillaforge/openstack-deploy

cd openstack-deploy

# Update kolla_external_vip_address with real IP
EXTERNAL_IP=$(curl -s ifconfig.me)
sed -i "s/#kolla_external_vip_address: \"{{ kolla_internal_vip_address }}\"/kolla_external_vip_address: \"$EXTERNAL_IP\"/" config/globals.yml

./install.sh
```


### Install LDAP Service and PegasusIAM

```bash
cd ~
git clone https://github.com/Zillaforge/mini-zillaforge-setup.git
cd mini-zillaforge-setup

HOSTNAME=$(hostname)
HOSTIP=$(curl -s ipinfo.io/ip)
HOSTIP_DASH=$(echo "$HOSTIP" | sed 's/\./-/g')

# Update configuration files with hostname
sed -i "s/instance-hx9bq8/$HOSTNAME/g" ./helm/mariadb-galera/values-trustedcloud.yaml
sed -i "s/instance-hx9bq8/$HOSTNAME/g" ./helm/redis-sentinel/values-trustedcloud.yaml

sudo mkdir -p /trusted-cloud/local/redis
sudo chmod -R 775 /trusted-cloud

helm install test-redis ./helm/redis-sentinel -f ./helm/redis-sentinel/values-trustedcloud.yaml
helm install test-mariadb ./helm/mariadb-galera -f ./helm/mariadb-galera/values-trustedcloud.yaml


kubectl create serviceaccount pegasus-system-admin

helm install pegasusiam ./helm/pegasusiam -f ./helm/pegasusiam/values-trustedcloud.yaml
helm install ldap-opsk ./helm/ldap -f ./helm/ldap/values-openstack.yaml

```

### Setup OpenStack domain with LDAPService

```bash
source ~/venv/bin/activate

export OS_CLIENT_CONFIG_FILE=/etc/kolla/clouds.yaml
export OS_CLOUD=kolla-admin
export USER_NAME="test@trusted-cloud.nchc.org.tw"
export PROJECT_NAME="trustedcloud"
export DOMAIN_NAME="trustedcloud"
export ROLE_NAME="admin"
export container_name="keystone"

openstack domain create $DOMAIN_NAME
openstack project create $PROJECT_NAME --domain $DOMAIN_NAME
sudo docker exec -it $container_name service apache2 restart


USER_ID=$(openstack user list --domain "$DOMAIN_NAME" -f value -c ID -c Name | grep "$USER_NAME" | awk '{print $1}')

PROJECT_ID=$(openstack project list -f value -c ID -c Name | grep "$PROJECT_NAME" | awk '{print $1}')

openstack role add --user "$USER_ID" --project "$PROJECT_ID" "$ROLE_NAME" 
openstack role add --user "$USER_ID" --domain "$DOMAIN_NAME" "$ROLE_NAME"


```

### Modify docker-compose to match testing environment

* Change admin project UUID

    ```bash
    # Replace UUID in all yaml files under docker-compose directory
    find docker-compose -name "*.yaml" -o -name "*.yml" | xargs sed -i 's/14735dfa-5553-46cc-b4bd-405e711b223f/14735dfa-5553-46cc-b4bd-405e711b223f/g'
    ```
* Change OpenStack identity endpoint

* Change PegasusIAM endpoint
