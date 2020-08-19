# Fabric-sdk-go demo 

 `fabric-sdk-go`的使用文档较少，这是使用`fabric-sdk-go`的一个案例。

本案例使用的`fabric`版本为`fabric 1.4.8`，单机部署的联盟链结构为：3`orderer`+2`org1`+2`org2`，`orderer`节点采用`Raft`共识。

多机部署参考：[Hyperledger Fabric Raft排序多机部署](http://blog.hubwiz.com/2019/12/24/fabric-raft-multi-host/)

## TODOs:

- [x] 编写`crypto-config.yaml`文件并生成秘钥文件
- [x] 编写`configtx.yaml`文件并生成创世快文件
- [x] 编写`docker-compose.yaml`文件，并运行成功
- [x] 编写链码
- [x] 为`fabric-sdk-go`编写配置文件，并调用成功

- [x] 使用`fabric-sdk-go`创建通道
- [x] 使用`fabric-sdk-go`把`org1`、`org2`（所有节点）加入通道
- [x] 使用`fabric-sdk-go`在指定节点安装自己编写的链码
- [x] 使用`fabric-sdk-go`实例化链码（配置背书策略）
- [x] 使用`fabric-sdk-go`调用链码
- [x] 使用`fabric-sdk-go`更新链码背书策略
- [x] 使用`fabric-sdk-go`调用新的链码（新的背书策略）

## 快速开始

1. Clone 该项目到你的电脑上

   ```shell
   git clone https://github.com/AlexTan-b-z/fabric-sdk-go-demo.git
   ```

2. 运行`fixtures/getFabric.sh`文件，来获取`fabric1.4.8`版本的docker镜像

3. 进入`fixtures`目录，运行`docker-compose-local.yaml`文件

   ```shell
   cd fixtures
   # 先创建local test网络
   docker network create local-test
   # 启动容器
   docker-compose -f docker-compose-local.yaml up -d
   ```

   运行成功后，使用`docker ps`命令能查看正在运行的容器

4. 拷贝项目根目录下的`chaincode`目录（链码），到你的`$GOPATH/src`目录下

   ```shell
   cd ..
   cp -r chaincode $GOPATH/src/
   ```

5. 分别运行`samplse/chaincode/main.go`和`samplse/event/main.go`

   ```
   go run samplse/chaincode/main.go
   go run samplse/event/main.go
   ```



#### [原文博客](https://blog.csdn.net/AlexTan_/article/details/108110927)