package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "strconv"

    "github.com/hyperledger/fabric/core/chaincode/shim"
)

//var txNo int = 0
var cpNo int = 0

type SimpleChaincode struct {

}

func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting SimpleChaincode: %s", err)
    }
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, funcName string, args []string) ([]byte, error) {
    if len(args) % 2 != 0 {
        return nil, errors.New("Incorrect number of arguments. Expecting EVEN number.")
    }
    if funcName != "init" {
        return nil, errors.New("Incorrect function name. Expecting init")
    }
    for i := 0; i < len(args); i = i+2 {
        name, asset := args[i], args[i+1]
        err := stub.PutState(name, asset)
        if err != nil {
            return nil, errors.New("PutState Error: %s, %s", args[i], args[i+1])
        }
    }
    return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, funcName string, args []string) ([]byte, error) {
    if len(args) != 3 {
        return nil, errors.New("Incorrect number of arguments. Expecting 3.")
    }
    from, to := args[0], args[1]
    val := strconv.Atoi(args[2])
    var (
        fromByte  []byte
        toByte    []byte
        fromAsset int
        toAsset   int
        err       error
    )
    fromByte, err = getAsset(stub, from)
    if err != nil {
        return nil, err
    }
    fromAsset = strconv.Atoi(string(fromByte))
    toByte, err = getAsset(stub, to)
    if err != nil {
        return nil, err
    }
    toAsset = strconv.Atoi(string(toByte))

    fromAsset, toAsset = fromAsset-val, toAsset+val

    err = stub.PutState(from, []byte(strconv.Itoa(fromAsset)))
    if err != nil {
        return nil, errors.New("PutState Error: %s", args[0])
    }
    err = stub.PutState(to, []byte(strconv.Itoa(toAsset)))
    if err != nil {
        return nil, errors.New("PutState Error: %s", args[1])
    }
    return nil, nil
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, funcName string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }
    if funcName != "query" {
        return nil, errors.New("Incorrect function name. Expecting query")
    }
    return getAsset(stub, args[0])
}

func getAsset(stub *shim.ChaincodeStub, args string) ([]byte, error) {
    assetByte, err := stub.GetState(args)
    if err != nil {
        return nil, errors.New("GetState Error: "+err.Error())
    }
    if assetByte == nil {
        return nil, errors.New("GetState Error: Nil value for "+args)
    }
    return assetByte, nil
}
