FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Install core utilities
RUN apt-get update && apt-get install -y \
    curl wget jq yq \
    git \
    gcc make python3 python3-pip \
    vim nano \
    netcat telnet dnsutils \
    unzip tar gzip \
    && rm -rf /var/lib/apt/lists/*

# AWS CLI
RUN pip3 install --no-cache-dir awscli boto3

# Azure CLI
RUN curl -sSL https://packages.microsoft.com/keys/microsoft.asc | apt-key add - \
    && echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list \
    && apt-get update && apt-get install -y azure-cli \
    && rm -rf /var/lib/apt/lists/*

# Google Cloud SDK
RUN curl -sSL https://dl.google.com/cloud-sdk/install.sh | bash
ENV PATH="/usr/local/google-cloud-sdk/bin:$PATH"

# Kubernetes
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl && mv kubectl /usr/local/bin/

# Helm
RUN curl -fsSL https://get.helm.sh/helm-v3.14.0-linux-amd64.tar.gz -o helm.tar.gz \
    && tar -zxvf helm.tar.gz \
    && mv linux-amd64/helm /usr/local/bin/ \
    && rm -rf helm.tar.gz linux-amd64

# Terraform
RUN curl -fsSL https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip -o terraform.zip \
    && unzip terraform.zip \
    && mv terraform /usr/local/bin/ \
    && rm terraform.zip

# Ansible
RUN pip3 install --no-cache-dir ansible ansible-core

# Docker
COPY --from=docker:latest /usr/local/bin/docker /usr/local/bin/

# Python libraries
RUN pip3 install --no-cache-dir \
    requests pyyaml kubernetes docker boto3 azure-identity google-auth flask celery

# Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# Working directory
RUN mkdir -p /workspace /scripts
WORKDIR /workspace

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
