/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/op/go-logging"
)

var myLogger = logging.MustGetLogger("test_contract")

// TODO: review the explanation later.
// TestContractChainCode is simple chaincode implementing a basic Contract
// https://github.com/hyperledger/fabric/blob/master/docs/tech/application-ACL.md
type TestContractChainCode struct {
}

type Asset struct {
	ID string
	Owner   string
	FishName string
	Weight   uint64
	MinTemperature uint64
	MaxTemperature uint64
	Price uint64
	Location string
}

type Contract struct {
	AssetID string
	PreviousTxId string
}


// Init method will be called during deployment.
// args: contractId(string), seller(string), fishName(uint), price(uint), weight(uint)
// TODO: confirm what seller value should be for later tracking. ECA/ TCA etc.
func (t *TestContractChainCode) Init(stub shim.ChaincodeStubInterface, methodName string, args []string) ([]byte, error) {
	myLogger.Debug("Init Chaincode...done")
	return nil, nil
}

func (t *TestContractChainCode) create_supply_chain(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}

	assetId := args[0]
	owner := args[1]
	fishName := args[2]
	price, err := strconv.ParseUint(args[3], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse price.")
	}
	weight, err := strconv.ParseUint(args[4], 10, 64)
	if err != nil {
		return nil, errors.New("Failed to parse weight.")
	}

	location := args[5]

	asset := Asset{
		ID:       assetId,
		Owner:   owner,
		FishName: fishName,
		Price:    price,
		Weight:   weight,
		MinTemperature: 0,
		MaxTemperature: 0,
		Location: location,
	}
	b, err := json.Marshal(asset)

	stub.PutState("asset/"+assetId, b)

	return nil, nil
}

func (t *TestContractChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if(function == "create_supply_chain"){
		return t.create_supply_chain(stub, args)
	}

	return nil, nil
}

// args: transaction id

//result:
   //print all the previous transactions data structure

func (t *TestContractChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	myLogger.Debug("Query Chaincode...")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	contractId := args[0]
	state, err := stub.GetState("contracts/" + contractId)
	return state, err
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(TestContractChainCode))
	if err != nil {
		fmt.Printf("Error starting TestContractChainCode: %s", err)
	}
}
