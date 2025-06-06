echo "Installing system dependencies..."
apt-get update
apt-get install -y \
    curl \
    git \
    unzip \
    apt-transport-https \
    ca-certificates \
    software-properties-common \
    jq \
    make \
    gcc

echo "Installing Docker..."
apt-get remove -y docker docker-engine docker.io containerd runc || true
apt-get update
apt-get install -y docker.io
systemctl enable --now docker
usermod -aG docker "$USER"

echo "Installing Docker Compose..."
COMPOSE_VERSION="1.29.2"
curl -L "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

echo "Installing Golang..."
GO_VERSION="1.20.5"
curl -L https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -o go.tar.gz
tar -xzf go.tar.gz
mv go /usr/local/
cp /usr/local/go/bin/* /usr/local/bin/
chmod +x /usr/local/bin/*
rm go.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz
echo "Golang version $(go version) installed."

echo "Installing specific Hyperledger Fabric binaries..."

HYPERLEDGER_VERSION="2.5.11"
FABRIC_BINARIES_URL="https://github.com/hyperledger/fabric/releases/download/v${HYPERLEDGER_VERSION}/hyperledger-fabric-linux-amd64-${HYPERLEDGER_VERSION}.tar.gz"
curl -L ${FABRIC_BINARIES_URL} -o fabric-binaries.tar.gz
tar -xzf fabric-binaries.tar.gz
cp bin/{peer,orderer,configtxgen,configtxlator,cryptogen,osnadmin,discover,ledgerutil} /usr/local/bin/
chmod +x /usr/local/bin/{peer,orderer,configtxgen,configtxlator,cryptogen,osnadmin,discover,ledgerutil}
rm -rf bin builders config fabric-binaries.tar.gz
echo "Selected Hyperledger Fabric binaries installed."

# Clone the repository
echo "Cloning repository from https://github.com/Duskmate/chaincare.git..."
git clone https://github.com/Duskmate/chaincare.git $(pwd)

# Check if clone was successful
if [ $? -ne 0 ]; then
    echo "Failed to clone repository. Exiting..."
    exit 1
fi

echo "Setup complete. Files downloaded."