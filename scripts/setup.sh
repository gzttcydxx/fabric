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
        if ! grep -q "/usr/local/bin/fabric-bin" ~/.zshrc; then
            echo 'export PATH=$PATH:/usr/local/bin/fabric-bin' >> ~/.zshrc
            source ~/.zshrc
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
    if ! grep -q "GOPATH" ~/.zshrc; then
        echo 'export GOPATH=$HOME/go' >> ~/.zshrc
        echo 'export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin' >> ~/.zshrc
    fi
    
    # 清理下载文件
    rm go1.22.5.linux-amd64.tar.gz
    
    # 重新加载环境变量
    source ~/.bashrc
    source ~/.zshrc

    echo "Go 1.22.5 installation completed"
}

# 安装 Docker
install_docker() {
    echo "Installing Docker..."
    # 检查docker版本是否 >= 27.0
    if command -v docker >/dev/null 2>&1; then
        docker_version=$(docker --version | grep -oP 'Docker version \K\d+\.\d+\.\d+')
        major_minor_version=$(echo $docker_version | awk -F. '{print $1"."$2}')
        if awk "BEGIN {exit !(echo $major_minor_version >= 27.0)}"; then
            echo "Docker $docker_version is already installed"
            return
        fi
    fi
    # 卸载旧版本
    apt-get remove docker docker-engine docker.io containerd runc
    # 更新apt
    apt-get update
    # 安装依赖
    apt-get install ca-certificates curl gnupg lsb-release make
    # 创建目录
    sudo install -m 0755 -d /etc/apt/keyrings
    # 下载并添加GPG密钥
    curl -fsSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    # 设置正确的权限
    sudo chmod a+r /etc/apt/keyrings/docker.gpg
    # 添加Docker源
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    # 安装 Docker
    apt-get update
    apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    # 检查docker版本
    docker --version
    # 将所有用户添加到docker组
    for USER_HOME in /home/*/ /root/; do
        # 跳过不存在的目录
        [ ! -d "$USER_HOME" ] && continue
        USER=$(basename $USER_HOME)
        usermod -aG docker $USER
    done
    # 设置开机自启
    systemctl enable docker
    # 配置docker镜像源
    cat > /etc/docker/daemon.json << 'EOL'
{
    "registry-mirrors": ["https://docker.gzttc.top"]
}
EOL
    # 重启docker
    systemctl restart docker
}

# 安装 Zsh
install_zsh() {
    # 检查是否已安装 Zsh
    if command -v zsh >/dev/null 2>&1; then
        echo "Zsh is already installed"
    else
        echo "Installing Zsh..."
        
        # 安装 Zsh
        apt-get update
        apt-get install -y zsh
    fi

    # 安装 Oh My Zsh
    export REMOTE="https://gh.gzttc.top/https://github.com/ohmyzsh/ohmyzsh.git"
    export ZSH="/usr/share/oh-my-zsh"
    sh -c "$(curl -fsSL https://gh.gzttc.top/https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended

    # 为所有用户配置 Oh My Zsh
    # 创建全局默认配置
    cat > /etc/skel/.zshrc << 'EOL'
export ZSH="/usr/share/oh-my-zsh"
ZSH_THEME="powerlevel10k/powerlevel10k"
plugins=(
    git 
    jsontools 
    z 
    vi-mode 
    copypath 
    copyfile 
    sudo 
    zsh-syntax-highlighting 
    zsh-autosuggestions 
    you-should-use
)
source $ZSH/oh-my-zsh.sh
export ZSH_AUTOSUGGEST_STRATEGY=(history completion)
export YSU_MESSAGE_POSITION="after"
EOL

    # 安装插件和主题
    OH_MY_ZSH_CUSTOM="/usr/share/oh-my-zsh/custom"
    if [ ! -d "$OH_MY_ZSH_CUSTOM/plugins/zsh-autosuggestions" ]; then
        git clone https://gh.gzttc.top/https://github.com/zsh-users/zsh-autosuggestions $OH_MY_ZSH_CUSTOM/plugins/zsh-autosuggestions
    fi
    if [ ! -d "$OH_MY_ZSH_CUSTOM/plugins/zsh-syntax-highlighting" ]; then
        git clone https://gh.gzttc.top/https://github.com/zsh-users/zsh-syntax-highlighting.git $OH_MY_ZSH_CUSTOM/plugins/zsh-syntax-highlighting
    fi
    if [ ! -d "$OH_MY_ZSH_CUSTOM/plugins/you-should-use" ]; then
        git clone https://gh.gzttc.top/https://github.com/MichaelAquilina/zsh-you-should-use.git $OH_MY_ZSH_CUSTOM/plugins/you-should-use
    fi
    if [ ! -d "$OH_MY_ZSH_CUSTOM/themes/powerlevel10k" ]; then
        git clone --depth=1 https://gh.gzttc.top/https://github.com/romkatv/powerlevel10k.git $OH_MY_ZSH_CUSTOM/themes/powerlevel10k
    fi

    # 为现有用户配置
    for USER_HOME in /home/*/ /root/; do
        # 跳过不存在的目录
        [ ! -d "$USER_HOME" ] && continue

        # 复制配置文件
        cp /etc/skel/.zshrc $USER_HOME/.zshrc
        
        # 设置正确的所有权
        USER=$(basename $USER_HOME)
        chown $USER:$USER $USER_HOME/.zshrc

        # 修改默认shell为zsh
        chsh -s $(which zsh) $USER
    done

    echo "Zsh installation completed and set as default shell for all users"
    echo "You can use `p10k configure` to configure powerlevel10k theme"
}

# Main function to call specific functions based on command line arguments
main() {
    if [ $# -eq 0 ]; then
        echo "Executing all setup functions..."
        configure_binaries
        configure_symlink
        configure_hosts
        install_go
        install_docker
        install_zsh
        echo "All setup functions completed."
        exit 0
    fi

    case "$1" in
        configure_binaries)
            configure_binaries
            ;;
        configure_symlink)
            configure_symlink
            ;;
        configure_hosts)
            configure_hosts
            ;;
        install_go)
            install_go
            ;;
        install_docker)
            install_docker
            ;;
        install_zsh)
            install_zsh
            ;;
        *)
            echo "Usage: $0 {configure_binaries|configure_symlink|configure_hosts|install_go|install_docker|install_zsh}"
            echo "       $0 (without arguments to execute all functions)"
            exit 1
            ;;
    esac
}

# Call the main function with all the script arguments
main "$@"
