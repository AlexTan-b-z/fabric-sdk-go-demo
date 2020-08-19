/*
 * @Author: AlexTan
 * @GIthub: https://github.com/AlexTan-b-z
 * @Date: 2020-08-04 20:35:47
 * @LastEditors: AlexTan
 * @LastEditTime: 2020-08-12 22:18:49
 */
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
)

/*
func stringToByte(args string) []byte{
	var as []byte
	for _, a := range args {
		as = append(as, []byte(a)...)
	}
	return as
}*/

type ScoreChaincode struct {

}

func (t *ScoreChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response{

	return shim.Success(nil)
}

func (t *ScoreChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()

	if fun == "addScore"{
		stub.SetEvent("addScore", []byte{})
		return t.addScore(stub, args)		// 添加成绩
	}else if fun == "queryScoreByNameAndGrade" {
		stub.SetEvent("queryScoreByNameAndGrade", []byte{})
		return t.queryScoreByNameAndGrade(stub, args)		// 根据姓名及年级查询成绩，假设不同年级存在同名
	}else if fun == "queryScoreDetailByStuID" {
		stub.SetEvent("queryScoreDetailByStuID", []byte{})
		return t.queryScoreDetailByStuID(stub, args)		// 根据学号查询成绩详细信息(包括历史信息)
	}else if fun == "updateScore" {
		stub.SetEvent("updateScore", []byte{})
		return t.updateScore(stub, args)		// 根据证书编号更新信息
	}else if fun == "delScore"{
		stub.SetEvent("delScore", []byte{})
		return t.delScore(stub, args)	// 根据证书编号删除信息
	}

	return shim.Error("指定的函数名称错误")

}

func main(){
	err := shim.Start(new(ScoreChaincode))
	if err != nil{
		fmt.Printf("启动ScoreChaincode时发生错误: %s", err)
	}
}