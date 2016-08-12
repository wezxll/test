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

type Company struct {
    Name        string
    Balance     int
    Id          int
}
/*
type Transaction struct {
    FromName        string
    FromId          int
    ToName          string
    ToId            int
    Number          int
    Time            int64
    txId            int
}
*/

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

    var cpName string
    var cpBal  int
    var cpId   int
    var cp     Company
    for i := 0; i < len(args); i+=2 {
        cpName = args[i]
        cpId   = cpNo
        cpBal, err := strconv.Atoi(args[i+1])
        if err != nil {
            return nil, errors.New("Expecting integer value for company balance.")
        }
        cpNo++
        cp = Company{Name: cpName, Balance: cpBal, id: cpId}
        err = writeCompany(stub, cp)
        if err != nil {
            return nil, errors.New("writeCompany Error" + err.Error())
        }
    }
    return nil, nil
}

func writeCompany(stub *shim.ChaincodeStub, cp Company) (error) {
    cpBytes, err := json.Marshal(&cp)
    if err != nil {
        return err
    }
    err = stub.PutState("company"+cp.Name, cpBytes)
    if err != nil {
        return errors.New("PutState Error" + err.Error())
    }
    return nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, funcName string, args []string) ([]byte, error) {
    if len(args) != 3 {
        return nil, errors.New("Incorrect number of arguments. Expecting 3.")
    }
    if funcName == "transfer" {
        var from Company
        var to   Company
        var x    int
        from, err := getCompanyByName(stub, args[0])
        if err != nil {
            return nil, err
        }
        to, err = getCompanyByName(stub, args[1])
        if err != nil {
            return nil, err
        }
        from.balance -= x
        to.balance   += x
        err = writeCompany(stub, from)
        if err != nil {
            return nil, err
        }
        err = writeCompany(stub, to)
        if err != nil {
            return nil, err
        }
    } else {
        return nil, errors.New("Incorrect function name.")
    }
    return nil, nil
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, funcName string, args []string) ([]byte, error) {
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }
    if funcName == "company" {
        var cpBytes []byte
        cpBytes, err := getCompanyByName(stub, "company"+args[0])
        if err != nil {
            return nil, errors.New("Query company Error"+err.Error())
        }
        return cpBytes, nil
    } else {
        return nil, errors.New("Incorrect function name")
    }
}

func getCompanyByName(stub *shim.ChaincodeStub, name string) ([]byte, error) {
    cpBytes, err := stub.GetState("company"+name)
    if err != nil {
        return nil, errors.New("GetState Error"+err.Error())
    }
    if cpBytes == nil {
        return nil, errors.New("Nil for "+name)
    }
    return cpBytes, nil
}
