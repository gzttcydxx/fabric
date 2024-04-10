source $LOCAL_ROOT_PATH/envpeer1soft

# 查询已提交的链码
output=$(peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name $CHAINCODE_NAME)

# 使用 awk 命令来解析版本和序列
version=$(echo "$output" | awk -F ', ' '/Version: / {print $2}' | awk -F ': ' '{print $2}')
sequence=$(echo "$output" | awk -F ', ' '/Sequence: / {print $2}' | awk -F ': ' '{print $2}')

export VERSION=$(echo "scale=2; $version + 0.1" | bc)
export SEQUENCE=$((sequence+1))

source scripts/code.sh
