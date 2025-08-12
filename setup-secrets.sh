#!/bin/bash
# Setup script for secrets

set -e

echo "Setting up Playlister secrets..."

# Create secrets directory
mkdir -p k8s/secrets

# Check if files exist
if [ ! -f "k8s/secrets/api-key.txt" ]; then
    echo -n "Enter your YouTube API key: "
    read -s api_key
    echo
    echo -n "$api_key" > k8s/secrets/api-key.txt
    echo "✓ API key saved"
fi

if [ ! -f "k8s/secrets/id_rsa" ]; then
    echo "Enter path to your SSH private key:"
    echo "  For RSA: ~/.ssh/id_rsa"
    echo "  For ed25519: ~/.ssh/id_ed25519"
    echo "  For ECDSA: ~/.ssh/id_ecdsa"
    read -p "Path (default: ~/.ssh/id_ed25519): " ssh_path
    ssh_path=${ssh_path:-~/.ssh/id_ed25519}
    
    if [ -f "$ssh_path" ]; then
        cp "$ssh_path" k8s/secrets/id_rsa
        chmod 600 k8s/secrets/id_rsa
        echo "✓ SSH key copied ($(file -b "$ssh_path" | cut -d' ' -f1-2))"
    else
        echo "❌ SSH key not found at $ssh_path"
        exit 1
    fi
fi

# Update git settings in kustomization.yaml
echo "Enter your git repository URL (e.g., git@github.com:user/repo.git): "
read git_url

echo "Enter git commit author name: "
read git_name

echo "Enter git commit author email: "
read git_email

# Update kustomization.yaml
sed -i "s|git@github.com:youruser/playlister-data.git|$git_url|g" k8s/kustomization.yaml
sed -i "s|Playlister Bot|$git_name|g" k8s/kustomization.yaml
sed -i "s|bot@yourcompany.com|$git_email|g" k8s/kustomization.yaml

echo "✓ Secrets configured!"
echo ""
echo "To deploy:"
echo "  kubectl apply -k k8s/"
echo ""
echo "To update container image:"
echo "  cd k8s && kustomize edit set image your-registry/playlister-downloader:latest"