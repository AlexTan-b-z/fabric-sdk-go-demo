package main

import (
	"encoding/hex"
	"log"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	
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
	org1Client := cli.New(org1CfgPath, "Org1", "Admin", "User1")
	org2Client := cli.New(org2CfgPath, "Org2", "Admin", "User1")
	defer org1Client.Close()
	defer org2Client.Close()

	// New event client
	cp := org1Client.SDK.ChannelContext(org1Client.ChannelID, fabsdk.WithUser(org1Client.OrgUser))

	ec, err := event.New(
		cp,
		event.WithBlockEvents(), // 如果没有，会是filtered
		// event.WithBlockNum(1), // 从指定区块获取，需要此参数
		event.WithSeekType(seek.Newest))
	if err != nil {
		log.Printf("Create event client error: %v", err)
	}

	// block event listen
	defer ec.Unregister(blockListener(ec))
	defer ec.Unregister(filteredBlockListener(ec))

	// tx listen
	txIDCh := make(chan string, 100)
	go txListener(ec, txIDCh)

	// chaincode event listen
	defer ec.Unregister(chainCodeEventListener(nil, ec))

	DoChainCode(org1Client, txIDCh)
	close(txIDCh)

	time.Sleep(time.Second * 10)
}

func blockListener(ec *event.Client) fab.Registration {
	// Register monitor block event
	beReg, beCh, err := ec.RegisterBlockEvent()
	if err != nil {
		log.Printf("Register block event error: %v", err)
	}
	log.Println("Registered block event")

	// Receive block event
	go func() {
		for e := range beCh {
			log.Printf("Receive block event:\nSourceURL: %v\nNumber: %v\nHash"+
				": %v\nPreviousHash: %v\n\n",
				e.SourceURL,
				e.Block.Header.Number,
				hex.EncodeToString(e.Block.Header.DataHash),
				hex.EncodeToString(e.Block.Header.PreviousHash))
		}
	}()

	return beReg
}

func filteredBlockListener(ec *event.Client) fab.Registration {
	// Register monitor filtered block event
	fbeReg, fbeCh, err := ec.RegisterFilteredBlockEvent()
	if err != nil {
		log.Printf("Register filtered block event error: %v", err)
	}
	log.Println("Registered filtered block event")

	// Receive filtered block event
	go func() {
		for e := range fbeCh {
			log.Printf("Receive filterd block event:\nNumber: %v\nlen("+
				"transactions): %v\nSourceURL: %v",
				e.FilteredBlock.Number, len(e.FilteredBlock.
					FilteredTransactions), e.SourceURL)

			for i, tx := range e.FilteredBlock.FilteredTransactions {
				log.Printf("tx index %d: type: %v, txid: %v, "+
					"validation code: %v", i,
					tx.Type, tx.Txid,
					tx.TxValidationCode)
			}
			log.Println() // Just go print empty log for easy to read
		}
	}()

	return fbeReg
}

func txListener(ec *event.Client, txIDCh chan string) {
	log.Println("Transaction listener start")
	defer log.Println("Transaction listener exit")

	for id := range txIDCh {
		// Register monitor transaction event
		log.Printf("Register transaction event for: %v", id)
		txReg, txCh, err := ec.RegisterTxStatusEvent(id)
		if err != nil {
			log.Printf("Register transaction event error: %v", err)
			continue
		}
		defer ec.Unregister(txReg)

		// Receive transaction event
		go func() {
			for e := range txCh {
				log.Printf("Receive transaction event: txid: %v, "+
					"validation code: %v, block number: %v",
					e.TxID,
					e.TxValidationCode,
					e.BlockNumber)
			}
		}()
	}
}

func chainCodeEventListener(c *cli.Client, ec *event.Client) fab.Registration {
	eventName := ".*"
	log.Printf("Listen chaincode event: %v", eventName)

	var (
		ccReg   fab.Registration
		eventCh <-chan *fab.CCEvent
		err     error
	)
	if c != nil {
		log.Println("Using client to register chaincode event")
		ccReg, eventCh, err = c.RegisterChaincodeEvent("mycc", eventName)
	} else {
		log.Println("Using event client to register chaincode event")
		ccReg, eventCh, err = ec.RegisterChaincodeEvent("mycc", eventName)
	}
	if err != nil {
		log.Printf("Register chaincode event error: %v", err.Error())
		return nil
	}

	// consume event
	go func() {
		for e := range eventCh {
			log.Printf("Receive cc event, ccid: %v \neventName: %v\n"+
				"payload: %v \ntxid: %v \nblock: %v \nsourceURL: %v\n",
				e.ChaincodeID, e.EventName, string(e.Payload), e.TxID, e.BlockNumber, e.SourceURL)
		}
	}()

	return ccReg
}

// Install、Deploy、Invoke、Query、Upgrade
func DoChainCode(cli1 *cli.Client, txCh chan<- string) {
	var (
		txid fab.TransactionID
		err  error
	)

	// ccVersion := "v1"
	// if err := cli1.InstallCC(ccVersion, peer0Org1); err != nil {
	// 	log.Panicf("Intall chaincode error: %v", err)
	// }
	// log.Println("Chaincode has been installed on org1's peers")
	//
	// // InstantiateCC chaincode only need once for each channel
	// if txid, err = cli1.InstantiateCC(ccVersion, peer0Org1); err != nil {
	// 	log.Panicf("Instantiated chaincode error: %v", err)
	// }
	// if txid != "" {
	// 	txCh <- string(txid)
	// }
	// log.Println("Chaincode has been instantiated")
	
	sco3 := model.Score{
		Name: "eventTest",
		Gender: "男",
		StuID: "888888",
		Grade: "2018",
		Result: "100",
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
	if txid, err = cli1.InvokeCCadd([]string{peer0Org1, "peer0.org2.example.com"}, sco3); err != nil {
		log.Panicf("InvokeCCadd chaincode error: %v", err)
	}
	log.Println("InvokeCCadd chaincode success 1")
	if txid != "" {
		txCh <- string(txid)
	}
	log.Println("Invoke chaincode success")

	if err = cli1.QueryCCByID("peer0.org1.example.com", "888888"); err != nil {
		log.Panicf("Query chaincode error: %v", err)
	}
	log.Println("Query chaincode success on peer0.org1")

	if txid, err = cli1.InvokeCCDelete([]string{"peer0.org1.example.com", "peer0.org2.example.com"}, "888888"); err != nil {
		log.Panicf("InvokeCCDelete chaincode error: %v", err)
	}
}
