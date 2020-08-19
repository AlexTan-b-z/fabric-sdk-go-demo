package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"fmt"
	"bytes"
)

const DOC_TYPE = "ScoreObj"

// 保存sco
// args: Score
func PutScore(stub shim.ChaincodeStubInterface, sco Score) ([]byte, bool) {

	sco.ObjectType = DOC_TYPE

	b, err := json.Marshal(sco)
	if err != nil {
		return nil, false
	}

	// 保存sco状态
	err = stub.PutState(sco.StuID, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// 根据学号查询信息
// args: stuID
func GetScoreInfo(stub shim.ChaincodeStubInterface, stuID string) (Score, bool)  {
	var sco Score
	// 根据学号查询信息状态
	b, err := stub.GetState(stuID)
	if err != nil {
		return sco, false
	}

	if b == nil {
		return sco, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &sco)
	if err != nil {
		return sco, false
	}

	// 返回结果
	return sco, true
}

// 根据指定的查询字符串实现富查询
func getEduByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer  resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}

// 添加信息
// args: scoreObject
// 学生号为 key, Score 为 value
func (t *ScoreChaincode) addScore(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1{
		return shim.Error("给定的参数个数不符合要求")
	}

	var sco Score
	err := json.Unmarshal([]byte(args[0]), &sco)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// 查重: 身份证号码必须唯一
	_, exist := GetScoreInfo(stub, sco.StuID)
	if exist {
		return shim.Error("要添加的学生号已存在")
	}

	_, bl := PutScore(stub, sco)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	return shim.Success([]byte("信息添加成功"))
}

// 根据姓名及年级查询信息
// args: Name, Grade
func (t *ScoreChaincode) queryScoreByNameAndGrade(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	name := args[0]
	grade := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"Name\":\"%s\", \"Grade\":\"%s\"}}", DOC_TYPE, name, grade)

	// 查询数据
	result, err := getEduByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据证书姓名及年级查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的姓名及年级没有查询到相关的信息")
	}
	return shim.Success(result)
}

// 根据学号查询详情（溯源）
// args: StuID
func (t *ScoreChaincode) queryScoreDetailByStuID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}

	// 根据学号查询sco状态
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据学号查询信息失败")
	}

	if b == nil {
		return shim.Error("根据学号没有查询到相关的信息")
	}

	// 对查询到的状态进行反序列化
	var sco Score
	err = json.Unmarshal(b, &sco)
	if err != nil {
		return  shim.Error("反序列化sco信息失败")
	}

	// 获取历史变更数据
	iterator, err := stub.GetHistoryForKey(sco.StuID)
	if err != nil {
		return shim.Error("根据指定的学号查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	// 迭代处理
	var historys []HistoryItem
	var hisSco Score
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取sco的历史变更数据失败")
		}

		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisSco)

		if hisData.Value == nil {
			var empty Score
			historyItem.Score = empty
		}else {
			historyItem.Score = hisSco
		}

		historys = append(historys, historyItem)

	}

	sco.Historys = historys

	// 返回
	result, err := json.Marshal(sco)
	if err != nil {
		return shim.Error("序列化sco信息时发生错误")
	}
	return shim.Success(result)
}

// 根据学号更新信息
// args: Score Object
func (t *ScoreChaincode) updateScore(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1{
		return shim.Error("给定的参数个数不符合要求")
	}

	var info Score
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}

	// 根据学号查询信息
	result, bl := GetScoreInfo(stub, info.StuID)
	if !bl{
		return shim.Error("根据学号查询信息时发生错误")
	}

	result.Name = info.Name
	result.Gender = info.Gender
	result.Grade = info.Grade
	result.Result = info.Result
	result.Time = info.Time

	_, bl = PutScore(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}

	return shim.Success([]byte("信息更新成功"))
}

// 根据学号删除信息（暂不对外提供）
// args: StuID
func (t *ScoreChaincode) delScore(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1{
		return shim.Error("给定的参数个数不符合要求")
	}

	/*var edu Education
	result, bl := GetEduInfo(stub, info.EntityID)
	err := json.Unmarshal(result, &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}*/

	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}
	return shim.Success([]byte("信息删除成功"))
}