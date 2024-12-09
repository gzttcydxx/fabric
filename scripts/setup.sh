#!/bin/bash

# 配置二进制工具
configure_binaries() {
    if [ ! -d "${LOCAL_ROOT_PATH}/bin" ]; then
        echo "Downloading Fabric binaries..."
        wget https://gh.gzttc.top/https://github.com/hyperledger/fabric/releases/download/v2.5.10/hyperledger-fabric-linux-amd64-2.5.10.tar.gz
        echo "Downloading Fabric CA client..."
        wget https://gh.gzttc.top/https://github.com/hyperledger/fabric-ca/releases/download/v1.5.13/hyperledger-fabric-ca-linux-amd64-1.5.13.tar.gz
        mkdir ${LOCAL_ROOT_PATH}/temp
        tar -xzf hyperledger-fabric-linux-amd64-2.5.10.tar.gz -C ${LOCAL_ROOT_PATH}/temp
        tar -xzf hyperledger-fabric-ca-linux-amd64-1.5.13.tar.gz -C ${LOCAL_ROOT_PATH}/temp
        mv ${LOCAL_ROOT_PATH}/temp/bin ${LOCAL_ROOT_PATH}/bin
        rm -rf ${LOCAL_ROOT_PATH}/temp
        rm ${LOCAL_ROOT_PATH}/hyperledger-fabric-linux-amd64-2.5.10.tar.gz
        rm ${LOCAL_ROOT_PATH}/hyperledger-fabric-ca-linux-amd64-1.5.13.tar.gz
    else
        echo "Fabric binaries already exist in ${LOCAL_ROOT_PATH}/bin"
    fi
}

# 判断是否需要创建符号链接
configure_symlink() {
    if [ ! -L "/usr/local/bin/fabric-bin" ]; then
        echo "Creating symbolic link for Fabric binaries..."
        ln -s ${LOCAL_ROOT_PATH}/bin /usr/local/bin/fabric-bin
        
        # 检查 PATH 中是否已经包含 fabric-bin
        if ! grep -q "/usr/local/bin/fabric-bin" ~/.bashrc; then
            echo 'export PATH=$PATH:/usr/local/bin/fabric-bin' >> ~/.bashrc
            source ~/.bashrc
        fi
    else
        echo "Symbolic link already exists at /usr/local/bin/fabric-bin"
    fi
}

# 配置 hosts 文件
configure_hosts() {
    echo "Configuring hosts file..."

    # 定义需要添加的 hosts
    HOSTS=(
        "127.0.0.1 traefik.example.com"
        "127.0.0.1 council.example.com"
        "127.0.0.1 soft.example.com"
        "127.0.0.1 web.example.com"
        "127.0.0.1 hard.example.com"
        "127.0.0.1 orderer1.council.example.com"
        "127.0.0.1 orderer2.council.example.com"
        "127.0.0.1 orderer3.council.example.com"
        "127.0.0.1 orderer1-admin.council.example.com"
        "127.0.0.1 orderer2-admin.council.example.com"
        "127.0.0.1 orderer3-admin.council.example.com"
        "127.0.0.1 peer1.soft.example.com"
        "127.0.0.1 peer1.web.example.com"
        "127.0.0.1 peer1.hard.example.com"
        "127.0.0.1 council.example.com"
        "127.0.0.1 orderer.example.com"
        "127.0.0.1 soft.example.com"
        "127.0.0.1 web.example.com"
        "127.0.0.1 couchdb.soft.example.com"
        "127.0.0.1 couchdb.web.example.com"
        "127.0.0.1 couchdb.hard.example.com"
        "127.0.0.1 orderer1.orderer.example.com"
        "127.0.0.1 peer1.soft.example.com"
        "127.0.0.1 peer1.web.example.com"
        "127.0.0.1 fabric.example.com"
    )

    # 检查并添加 hosts 记录
    for HOST in "${HOSTS[@]}"; do
        if ! grep -q "^$HOST" /etc/hosts; then
            echo "$HOST" >> /etc/hosts
            echo "Added: $HOST"
        else
            echo "Host already exists: $HOST"
        fi
    done
}

# 安装 Go 语言
install_go() {
    # 检查是否已安装 Go
    if command -v go >/dev/null 2>&1; then
        current_version=$(go version | awk '{print $3}')
        if [ "$current_version" = "go1.22.5" ]; then
            echo "Go 1.22.5 is already installed"
            return
        fi
    fi

    echo "Installing Go 1.22.5..."
    
    # 下载 Go
    wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
    
    # 删除旧版本（如果存在）
    sudo rm -rf /usr/local/go
    
    # 解压到 /usr/local
    sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
    
    # 设置环境变量（如果尚未设置）
    if ! grep -q "GOPATH" ~/.bashrc; then
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin' >> ~/.bashrc
    fi
    
    # 清理下载文件
    rm go1.22.5.linux-amd64.tar.gz
    
    # 重新加载环境变量
    source ~/.bashrc
    
    echo "Go 1.22.5 installation completed"
}

# 安装 Docker
install_docker() {
    echo "Installing Docker..."
    sudo apt-get update
    sudo apt-get install -y docker.io
    # TODO: 配置 Docker
}

configure_binaries
configure_symlink
configure_hosts
install_go
install_docker
