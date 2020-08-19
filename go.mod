module github.com/AlexTan-b-z/fabric-sdk-go-demo

go 1.13

require (
	github.com/hyperledger/fabric-protos-go v0.0.0-20200804194919-0c472b9cf543
	github.com/hyperledger/fabric-sdk-go v1.0.0-beta2
	github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190822125948-d2b42602e52e // indirect
	github.com/pkg/errors v0.9.1
)

replace github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190822125948-d2b42602e52e => github.com/hyperledger/fabric-sdk-go v1.0.0-beta2
