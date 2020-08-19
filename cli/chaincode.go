package cli

import (
	"log"
	"net/http"
	"strings"
	"encoding/json"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/pkg/errors"

	"github.com/AlexTan-b-z/fabric-sdk-go-demo/model"
)

// InstallCC install chaincode for target peer
func (c *Client) InstallCC(v string, peer string) error {
	targetPeer := resmgmt.WithTargetEndpoints(peer)

	// pack the chaincode
	ccPkg, err := gopackager.NewCCPackage(c.CCPath, c.CCGoPath)
	if err != nil {
		return errors.WithMessage(err, "pack chaincode error")
	}

	// new request of installing chaincode
	req := resmgmt.InstallCCRequest{
		Name:    c.CCID,
		Path:    c.CCPath,
		Version: v,
		Package: ccPkg,
	}

	resps, err := c.rc.InstallCC(req, targetPeer)
	if err != nil {
		return errors.WithMessage(err, "installCC error")
	}

	// check other errors
	var errs []error
	for _, resp := range resps {
		log.Printf("Install  response status: %v", resp.Status)
		if resp.Status != http.StatusOK {
			errs = append(errs, errors.New(resp.Info))
		}
		if resp.Info == "already installed" {
			log.Printf("Chaincode %s already installed on peer: %s.\n",
				c.CCID+"-"+v, resp.Target)
			return nil
		}
	}

	if len(errs) > 0 {
		log.Printf("InstallCC errors: %v", errs)
		return errors.WithMessage(errs[0], "installCC first error")
	}
	return nil
}

func (c *Client) InstantiateCC(v string, peer string) (fab.TransactionID,
	error) {
	// endorser policy
	org1OrOrg2 := "OR('Org1MSP.member','Org2MSP.member')"
	ccPolicy, err := c.genPolicy(org1OrOrg2)
	if err != nil {
		return "", errors.WithMessage(err, "gen policy from string error")
	}

	// new request
	// Attention: args should include `init` for Request not
	// have a method term to call init
	args := packArgs([]string{"init"})
	req := resmgmt.InstantiateCCRequest{
		Name:    c.CCID,
		Path:    c.CCPath,
		Version: v,
		Args:    args,
		Policy:  ccPolicy,
	}

	// send request and handle response
	reqPeers := resmgmt.WithTargetEndpoints(peer)
	resp, err := c.rc.InstantiateCC(c.ChannelID, req, reqPeers)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return "", nil
		}
		return "", errors.WithMessage(err, "instantiate chaincode error")
	}

	log.Printf("Instantitate chaincode tx: %s", resp.TransactionID)
	return resp.TransactionID, nil
}

func (c *Client) genPolicy(p string) (*common.SignaturePolicyEnvelope, error) {
	// TODO bug, this any leads to endorser invalid
	if p == "ANY" {
		return cauthdsl.SignedByAnyMember([]string{c.OrgName}), nil
	}
	return cauthdsl.FromString(p)
}

func (c *Client) InvokeCCadd(peers []string, sco model.Score) (fab.TransactionID, error) {
	//将sco序列化成字节数组
	sco_b, err := json.Marshal(sco)
	if err != nil {
		log.Printf("指定的Score对象序列化时发生错误: %v", err)
		return "", err
	}

	// new channel request for invoke
	req := channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "addScore",
		Args:        [][]byte{sco_b},
	}

	// send request and handle response
	// peers is needed
	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := c.cc.Execute(req, reqPeers)
	log.Printf("Invoke chaincode response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)
	if err != nil {
		return "", errors.WithMessage(err, "invoke chaincode error")
	}

	return resp.TransactionID, nil
}

func (c *Client) QueryCCByNameAndGrade(peer, name string, grade string) error {
	// new channel request for query
	req := channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "queryScoreByNameAndGrade",
		Args:        packArgs([]string{name, grade}),
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer)
	resp, err := c.cc.Query(req, reqPeers)
	if err != nil {
		return errors.WithMessage(err, "query chaincode error")
	}

	log.Printf("Query chaincode tx response:\ntx: %s\nresult: %v\n\n",
		resp.TransactionID,
		string(resp.Payload))
	return nil
}

func (c *Client) QueryCCByID(peer, StuID string) error {
	// new channel request for query
	req := channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "queryScoreDetailByStuID",
		Args:        packArgs([]string{StuID}),
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer)
	resp, err := c.cc.Query(req, reqPeers)
	if err != nil {
		return errors.WithMessage(err, "query chaincode error")
	}

	log.Printf("Query chaincode tx response:\ntx: %s\nresult: %v\n\n",
		resp.TransactionID,
		string(resp.Payload))
	return nil
}

func (c *Client) UpdateCCScore(peers []string, sco model.Score) (fab.TransactionID, error) {
	//将sco序列化成字节数组
	sco_b, err := json.Marshal(sco)
	if err != nil {
		log.Printf("指定的Score对象序列化时发生错误: %v", err)
		return "", err
	}

	// new channel request for invoke
	req := channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "updateScore",
		Args:        [][]byte{sco_b},
	}

	// send request and handle response
	// peers is needed
	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := c.cc.Execute(req, reqPeers)
	log.Printf("Invoke chaincode response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)
	if err != nil {
		return "", errors.WithMessage(err, "invoke chaincode error")
	}

	return resp.TransactionID, nil
}

func (c *Client) InvokeCCDelete(peers []string, stuID string) (fab.TransactionID, error) {
	log.Println("Invoke delete")
	// new channel request for invoke
	req := channel.Request{
		ChaincodeID: c.CCID,
		Fcn:         "delScore",
		Args:        packArgs([]string{stuID}),
	}

	// send request and handle response
	// peers is needed
	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := c.cc.Execute(req, reqPeers)
	log.Printf("Invoke chaincode delete response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)
	if err != nil {
		return "", errors.WithMessage(err, "invoke chaincode error")
	}

	return resp.TransactionID, nil
}

func (c *Client) UpgradeCC(v string, peer string) error {
	// endorser policy
	org1AndOrg2 := "AND('Org1MSP.member','Org2MSP.member')"
	ccPolicy, err := c.genPolicy(org1AndOrg2)
	if err != nil {
		return errors.WithMessage(err, "gen policy from string error")
	}

	// new request
	// Attention: args should include `init` for Request not
	// have a method term to call init
	// Reset a b's value to test the upgrade
	args := packArgs([]string{"init"})
	req := resmgmt.UpgradeCCRequest{
		Name:    c.CCID,
		Path:    c.CCPath,
		Version: v,
		Args:    args,
		Policy:  ccPolicy,
	}

	// send request and handle response
	reqPeers := resmgmt.WithTargetEndpoints(peer)
	resp, err := c.rc.UpgradeCC(c.ChannelID, req, reqPeers)
	if err != nil {
		return errors.WithMessage(err, "instantiate chaincode error")
	}

	log.Printf("Instantitate chaincode tx: %s", resp.TransactionID)
	return nil
}

func (c *Client) QueryCCInfo(v string, peer string) {

}

func (c *Client) Close() {
	c.SDK.Close()
}

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}
