/*
 * @Author: AlexTan
 * @GIthub: https://github.com/AlexTan-b-z
 * @Date: 2020-08-12 21:13:19
 * @LastEditors: AlexTan
 * @LastEditTime: 2020-08-12 22:07:26
 */

package cli

import (
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
)

const (
	org1CfgPath = "../../sdkConfig/org1_config.yaml"
	org2CfgPath = "../../sdkConfig/org2_config.yaml"

	channelConfig = "../../fixtures/channel-artifacts/channel.tx"
	ordererID = "orderer0.example.com"
	ordererOrgName = "ordererorg"
	org1Name = "Org1"
	org2Name = "Org2"
	orgAdmin = "Admin"
	channelID = "mychannel"
)

func CreateChannel(){
	sdk1, err := fabsdk.New(config.FromFile(org1CfgPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk1: %s", err)
	}
	sdk2, err := fabsdk.New(config.FromFile(org2CfgPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk2: %s", err)
	}

	clientContext := sdk1.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(ordererOrgName))
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Panicf("failed to create resMgmtClient in createChannel: %s", err)
	}

	mspClient, err := mspclient.New(sdk1.Context(), mspclient.WithOrg(org1Name))
	if err != nil {
		log.Panicf("failed to create msp client: %s", err)
	}
	adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)
	if err != nil {
		log.Panicf("failed to GetSigningIdentity: %s", err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: channelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
		
	_, err = resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererID))
	if err != nil {
		log.Panicf("failed to GetSigningIdentity: %s", err)
	}
	log.Println("created fabric channel")

	// join Channel
	org1Context := sdk1.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(org1Name))
	org2Context := sdk2.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(org2Name))

	org1ResMgmt, err := resmgmt.New(org1Context)
	if err != nil {
		log.Panicf("failed to create org1ResMgmt: %s", err)
	}
	org2ResMgmt, err := resmgmt.New(org2Context)
	if err != nil {
		log.Panicf("failed to create org2ResMgmt: %s", err)
	}
	if err = org1ResMgmt.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererID)); err != nil {
		log.Panicf("Org1 peers failed to JoinChannel: %s", err)
	}
	log.Println("org1 joined channel")
	if err = org2ResMgmt.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererID)); err != nil {
		log.Panicf("Org2 peers failed to JoinChannel: %s", err)
	}
	log.Println("org2 joined channel")
}