/*
 * @Author: test1Tan
 * @GIthub: https://github.com/test1Tan-b-z
 * @Date: 2020-08-11 21:03:54
 * @LastEditors: AlexTan
 * @LastEditTime: 2020-08-19 20:20:44
 */
package main

import (
	"log"
	"time"

	"github.com/AlexTan-b-z/fabric-sdk-go-demo/cli"
	"github.com/AlexTan-b-z/fabric-sdk-go-demo/model"
)

const (
	org1CfgPath = "../../sdkConfig/org1_config.yaml"
	org2CfgPath = "../../sdkConfig/org2_config.yaml"
)

var (
	peer0Org1 = "peer0.org1.example.com"
	peer0Org2 = "peer0.org2.example.com"
)

func main() {

	// init
	// 第一次运行后记得注释掉
	cli.CreateChannel()

	org1Client := cli.New(org1CfgPath, "Org1", "Admin", "User1")
	org2Client := cli.New(org2CfgPath, "Org2", "Admin", "User1")

	defer org1Client.Close()
	defer org2Client.Close()

	// Install, instantiate, invoke, query
	Phase1(org1Client, org2Client)
	// Install, upgrade, invoke, query
	Phase2(org1Client, org2Client)
}

func Phase1(cli1, cli2 *cli.Client) {
	log.Println("=================== Phase 1 begin ===================")
	defer log.Println("=================== Phase 1 end ===================")

	if err := cli1.InstallCC("v1", peer0Org1); err != nil {
		log.Panicf("Intall chaincode error: %v", err)
	}
	log.Println("Chaincode has been installed on org1's peers")

	if err := cli2.InstallCC("v1", peer0Org2); err != nil {
		log.Panicf("Intall chaincode error: %v", err)
	}
	log.Println("Chaincode has been installed on org2's peers")

	// InstantiateCC chaincode only need once for each channel
	if _, err := cli1.InstantiateCC("v1", peer0Org1); err != nil {
		log.Panicf("Instantiated chaincode error: %v", err)
	}
	log.Println("Chaincode has been instantiated")

	sco1 := model.Score{
		Name: "test1",
		Gender: "男",
		StuID: "123",
		Grade: "2015",
		Result: "100",
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}

	sco2 := model.Score{
		Name: "test2",
		Gender: "女",
		StuID: "1234",
		Grade: "2017",
		Result: "100",
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}

	if _, err := cli1.InvokeCCadd([]string{peer0Org1}, sco1); err != nil {
		log.Panicf("InvokeCCadd test1 chaincode error: %v", err)
	}
	log.Println("InvokeCCadd test1 chaincode success 1")

	if _, err := cli1.InvokeCCadd([]string{peer0Org1}, sco2); err != nil {
		log.Panicf("InvokeCCadd test2 chaincode error: %v", err)
	}
	log.Println("InvokeCCadd test2 chaincode success 2")

	if err := cli1.QueryCCByNameAndGrade("peer0.org1.example.com", "test1", "2015"); err != nil {
		log.Panicf("QueryCCByNameAndGrade chaincode error: %v", err)
	}
	log.Println("QueryCCByNameAndGrade chaincode success on peer0.org1")

	if err := cli1.QueryCCByID("peer0.org1.example.com", "1234"); err != nil {
		log.Panicf("QueryCCByID chaincode error: %v", err)
	}
	log.Println("QueryCCByID chaincode success on peer0.org1")

	new_sco := model.Score{
		Name: "test1",
		Gender: "男",
		StuID: "123",
		Grade: "2015",
		Result: "99",
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
	if _, err := cli1.UpdateCCScore([]string{peer0Org1}, new_sco); err != nil {
		log.Panicf("updateCCScore chaincode error: %v", err)
	}
	log.Println("updateCCScore chaincode success 1")

	if err := cli1.QueryCCByID("peer0.org1.example.com", "123"); err != nil {
		log.Panicf("QueryCCByID chaincode error: %v", err)
	}
	log.Println("QueryCCByID chaincode success on peer0.org1")

	if _, err := cli1.InvokeCCDelete([]string{"peer0.org1.example.com"}, "123"); err != nil {
		log.Panicf("InvokeCCDelete chaincode error: %v", err)
	}
	log.Println("InvokeCCDelete chaincode success on peer0.org1")
}

func Phase2(cli1, cli2 *cli.Client) {
	log.Println("=================== Phase 2 begin ===================")
	defer log.Println("=================== Phase 2 end ===================")

	v := "v2"

	// Install new version chaincode
	if err := cli1.InstallCC(v, peer0Org1); err != nil {
		log.Panicf("Intall chaincode error: %v", err)
	}
	log.Println("Chaincode has been installed on org1's peers")

	if err := cli2.InstallCC(v, peer0Org2); err != nil {
		log.Panicf("Intall chaincode error: %v", err)
	}
	log.Println("Chaincode has been installed on org2's peers")

	// Upgrade chaincode only need once for each channel
	if err := cli1.UpgradeCC(v, peer0Org1); err != nil {
		log.Panicf("Upgrade chaincode error: %v", err)
	}
	log.Println("Upgrade chaincode success for channel")

	sco1 := model.Score{
		Name: "test3",
		Gender: "男",
		StuID: "12345",
		Grade: "2015",
		Result: "100",
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}

	sco2 := model.Score{
		Name: "test4",
		Gender: "女",
		StuID: "123456",
		Grade: "2017",
		Result: "100",
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}

	if _, err := cli1.InvokeCCadd([]string{"peer0.org1.example.com", "peer0.org2.example.com"}, sco1); err != nil {
		log.Panicf("InvokeCCadd test3 chaincode error: %v", err)
	}
	log.Println("InvokeCCadd test3 chaincode success 1")

	if _, err := cli1.InvokeCCadd([]string{peer0Org1, "peer0.org2.example.com"}, sco2); err != nil {
		log.Panicf("InvokeCCadd test4 chaincode error: %v", err)
	}
	log.Println("InvokeCCadd test4 chaincode success 2")

	if err := cli1.QueryCCByNameAndGrade("peer0.org2.example.com", "test3", "2015"); err != nil {
		log.Panicf("QueryCCByNameAndGrade chaincode error: %v", err)
	}
	log.Println("QueryCCByNameAndGrade chaincode success on peer0.org2")

	if err := cli1.QueryCCByID("peer0.org2.example.com", "12345"); err != nil {
		log.Panicf("QueryCCByID chaincode error: %v", err)
	}
	log.Println("QueryCCByID chaincode success on peer0.org2")

	new_sco := model.Score{
		Name: "test3",
		Gender: "男",
		StuID: "12345",
		Grade: "2015",
		Result: "99",
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
	if _, err := cli1.UpdateCCScore([]string{peer0Org1, "peer0.org2.example.com"}, new_sco); err != nil {
		log.Panicf("updateCCScore chaincode error: %v", err)
	}
	log.Println("updateCCScore chaincode success 1")

	if err := cli1.QueryCCByID("peer0.org2.example.com", "12345"); err != nil {
		log.Panicf("QueryCCByID chaincode error: %v", err)
	}
	log.Println("QueryCCByID chaincode success on peer0.org2")

	if _, err := cli1.InvokeCCDelete([]string{"peer0.org2.example.com", peer0Org1}, "12345"); err != nil {
		log.Panicf("InvokeCCDelete chaincode error: %v", err)
	}
	log.Println("InvokeCCDelete chaincode success on peer0.org2")
}
