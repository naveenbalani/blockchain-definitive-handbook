package main

import (
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

type TradeContract struct {
	tradeId string //used
	buyerTaxId string //used
	skuid string //used
	sellerTaxId string //used
	exportBankId string // used
	importBankId string // used
	deliveryDate string
    shipperId string
	status string // used

	tradePrice int //used
	shippingPrice int //used
	totalPrice int
}

func (t *TradeContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return setupTrade(stub);
}

func (t *TradeContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, _ := stub.GetFunctionAndParameters()
	if (function == "createLOC") {
		return t.createLOC(stub, t.tradeId)
	}

	return shim.Error("Invalid function name")
}

func setupTrade(stub shim.ChaincodeStubInterface) peer.Response {
	tradeContract := TradeContract {
		tradeId: "FTE_1", 
		buyerTaxId: "FTE_B_1", 
		sellerTaxId: "FTE_S_1", 
		skuid: "SKU001",
		tradePrice: 10000,
		shippingPrice: 1000,
		status: "Trade initiated"}

	tcBytes, _ := json.Marshal(tradeContract)
	stub.PutState(tradeContract.tradeId, tcBytes)
	
	return shim.Success(nil)
}

func (t *TradeContract) createLOC(stub shim.ChaincodeStubInterface, tradeId string) peer.Response {
	tcBytes, _ := stub.GetState(tradeId)
	tc := TradeContract{}
	json.Unmarshal(tcBytes, &tc)

	if (tc.status == "Trade initiated") {
		tc.importBankId = "BNK_I_1"
		tc.status = "LOC created"
	} else {
		fmt.Printf("Trade not initiated yet")
	}

	tcBytes, _ = json.Marshal(tc)
	stub.PutState(tc.tradeId, tcBytes)
	
	return shim.Success(nil)
}

func (t *TradeContract) approveLOC(stub shim.ChaincodeStubInterface, tradeId string) peer.Response {
	tcBytes, _ := stub.GetState(tradeId)
	tc := TradeContract{}
	json.Unmarshal(tcBytes, &tc)

	if (tc.status == "LOC created") {
		tc.exportBankId = "BNK_E_1"
		tc.status = "LOC approved"
	} else {
		fmt.Printf("LOC not found")
	}

	tcBytes, _ = json.Marshal(tc)
	stub.PutState(tc.tradeId, tcBytes)

	return shim.Success(nil)
}
	
func (t *TradeContract) initiateShipment(stub shim.ChaincodeStubInterface, tradeId string) peer.Response {
	tcBytes, _ := stub.GetState(tradeId)
	tc := TradeContract{}
	json.Unmarshal(tcBytes, &tc)

	if (tc.status == "LOC approved") {
		tc.deliveryDate = "2017-10-31"
		tc.status = "Shipment initiated"
	} else {
		fmt.Printf("LOC not found")
	}

	tcBytes, _ = json.Marshal(tc)
	stub.PutState(tc.tradeId, tcBytes)
	
	return shim.Success(nil)
}

func (t *TradeContract) deliverGoods(stub shim.ChaincodeStubInterface, tradeId string) peer.Response {
	tcBytes, _ := stub.GetState(tradeId)
	tc := TradeContract{}
	json.Unmarshal(tcBytes, &tc)

	if (tc.status == "Shipment initiated") {
		tc.shipperId = "SHP_1"
		tc.status = "BOL created"
		fmt.Printf("Trade complete")
	} else {
		fmt.Printf("Shipment not initiated yet")
	}

	tcBytes, _ = json.Marshal(tc)
	stub.PutState(tc.tradeId, tcBytes)
	
	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *TradeContract) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]
	tcBytes, _ := stub.GetState(A)
	tc := TradeContract{}
	json.Unmarshal(tcBytes, &tc)


	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil trade for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Trade\":\"" + A + "\",\"Trade\":\"" + string(tc.status) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {

	err := shim.Start(new(TradeContract))
	if err != nil {
		fmt.Printf("Error creating new Trade Contract: %s", err)
	}
}
	
