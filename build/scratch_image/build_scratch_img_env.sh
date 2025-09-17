TMP=build/scratch_image/tmp/
PACKAGE=AppPlaygroundService
EXEC_FILE=$TMP/$PACKAGE
TMP_LINK_FILES=$TMP/link_files

echo $(pwd)
mkdir -p $TMP_LINK_FILES
linkPkgs=$(ldd $EXEC_FILE | egrep -o '\/[\/a-z0-9.+\_-]+')

for pkg in $linkPkgs; 
do 
    lPath=$TMP_LINK_FILES/$(dirname "$pkg" | cut -c2-)
    mkdir -p $lPath
    cp $pkg $lPath
done

cp /lib/$(uname -m)-linux-gnu/libresolv.so.2 $TMP_LINK_FILES/lib/$(uname -m)-linux-gnu/libresolv.so.2
cp /lib/$(uname -m)-linux-gnu/libnss_dns.so.2 $TMP_LINK_FILES/lib/$(uname -m)-linux-gnu/libnss_dns.so.2
cp /usr/local/go/lib/time/zoneinfo.zip $TMP_LINK_FILES/zoneinfo.zip

# Download terraform exec file
TERRAFORM_VERSION=1.9.8
ARCH=$(uname -m)
if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  TERRAFORM_ARCH=arm64
else
  TERRAFORM_ARCH=amd64
fi
TERRAFORM_URL="https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${TERRAFORM_ARCH}.zip"
TERRAFORM_ZIP=$TMP/terraform_${TERRAFORM_VERSION}_linux_${TERRAFORM_ARCH}.zip
TERRAFORM_BIN=$TMP_LINK_FILES/usr/bin

mkdir -p $TERRAFORM_BIN
wget -O $TERRAFORM_ZIP $TERRAFORM_URL
unzip $TERRAFORM_ZIP -d $TERRAFORM_BIN
rm $TERRAFORM_ZIP


# sed -i '' 's/$(PREVERSION)/$(VERSION)/g' $(PWD)/build/AETriton.spec

