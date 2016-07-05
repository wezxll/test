/*
	author:swb
	time:16/06/30
	MIT License
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var bankNo int = 0
var cpNo int = 0
var transactionNo int = 0

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type CenterBank struct {
	Name        string
	TotalNumber int
	RestNubmer  int
}

type Bank struct {
	Name        string
	TotalNumber int
	RestNubmer  int
	ID          int
}

type Company struct {
	Name   string
	Number int
	ID     int
}

type Transaction struct {
	FromType int //CenterBank 0 Bank 1  Company 1
	FromID   int
	ToType   int //Bank 1 Company 2
	ToID     int
	Time     int64
	Number   int
	ID       int
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	var totalNumber int
	var centerBank CenterBank
	totalNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	centerBank = CenterBank{Name: args[0], TotalNumber: totalNumber, RestNubmer: 0}
	centerBankBytes, err := json.Marshal(&centerBank)
	if err != nil {
		return nil, err
	}
	err = stub.PutState("centerBank", centerBankBytes)
	if err != nil {
		return nil, errors.New("PutState Error" + err.Error())
	}
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "createBank" {
		return t.createBank(stub, args)
	} else if function == "createCompany" {
		return t.createCompany(stub, args)
	} else if function == "issueCoin" {
		return t.issueCoin(stub, args)
	} else if function == "issueCoinToBank" {
		return t.issueCoinToBank(stub, args)
	} else if function == "issueCoinToCp" {
		return t.issueCoinToCp(stub, args)
	} else if function =="transfer"{
		return t.transfer(stub,args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) createBank(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var bank Bank
	bank = Bank{Name:args[0],TotalNumber:0,RestNubmer:0,ID:bankNo}
	bankBytes, err := json.Marshal(&bank)
	if err != nil {
		return nil, err
	}
	err = stub.PutState("bank"+strconv.ItoA(bankNo), bankBytes)
	if err != nil {
		return nil, errors.New("PutState Error" + err.Error())
	}
	bankNo = bankNo +1
	return bankBytes, nil
}

func (t *SimpleChaincode) createCompany(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var company Company
	company = Company{Name:args[0],Number:0,ID:cpNo}
	cpBytes,err := json.Marshal(&company)
	if(err!=nil){
		return nil,err
	}
	err = stub.PutState("company"+strconv.ItoA(cpNo),cpBytes)
	if err!= nil{
		return nil, errors.New("PutState Error" + err.Error())
	}
	cpNo = cpNo +1
	return cpBytes, nil
}

func (t *SimpleChaincode) issueCoin(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var centerBank CenterBank
	var centerBankBytes []byte
	var tsBytes []byte

	issueNumber := strconv.Atoi(args[0])
	cbBytes,err := stub.GetState("centerBank")
	if err != nil {
		fmt.Println("Error retrieving cbBytes")
	}
	err = json.Unmarshal(cbBytes, &centerBank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}

	centerBank.TotalNumber = centerBank.TotalNumber + issueNumber
	centerBank.RestNubmer = centerBank.RestNubmer + issueNumber
	centerBankBytes, err = json.Marshal(&centerBank)
	if err != nil {
		return nil, err
	}
	err = stub.PutState("centerBank", centerBankBytes)
	if err != nil {
		return nil, errors.New("PutState Error" + err.Error())
	}

	transaction := Transaction{FromType:0,FromID:0,ToType:0,ToId:0,Time:time.Now().Unix(),Number:issueNumber,ID:transactionId}
	tsBytes,err = json.Marshal(&transaction)
	if(err!=nil){
		return nil,err
	}
	err = stub.PutState("transaction"+strconv.ItoA(transactionNo),tsBytes)
	if err!= nil{
		return nil, errors.New("PutState Error" + err.Error())
	}
	transactionNo = transactionNo +1
	return tsBytes, nil
}

func (t *SimpleChaincode) issueCoinToBank(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	var centerBank CenterBank
	var bank Bank
	var bankId string
	var issueNumber int
	var bankBytes []byte


	bankId = args[0]
	issueNumber = strconv.Itoa(args[1])

	cbBytes,err := stub.GetState("centerBank")
	if err != nil {
		fmt.Println("Error retrieving cbBytes")
	}
	err = json.Unmarshal(cbBytes, &centerBank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	if centerBank.RestNubmer<issueNumber{
		return nil,errors.New("Not enough money")
	}

	bankBytes,err = stub.GetState("bank"+bankId)
	if err != nil {
		fmt.Println("Error retrieving cbBytes")
	}
	err = json.Unmarshal(bankBytes, &bank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	bank.Number = bank.Number + issueNumber
	centerBankBytes.RestNubmer = centerBank.RestNubmer - issueNumber


	cbBytes, err = json.Marshal(&centerBank)
	if err != nil {
		bank.Number = bank.Number - issueNumber
		centerBankBytes.RestNubmer = centerBank.RestNubmer + issueNumber
		return nil, err
	}
	err = stub.PutState("centerBank", cbBytes)
	if err != nil {
		bank.Number = bank.Number - issueNumber
		centerBankBytes.RestNubmer = centerBank.RestNubmer + issueNumber
		return nil, errors.New("PutState Error" + err.Error())
	}

	bankBytes, err = json.Marshal(&bank)
	if err != nil {
		//这里的代码逻辑有问题，需要撤销对centerBank的修改
		bank.Number = bank.Number - issueNumber
		centerBankBytes.RestNubmer = centerBank.RestNubmer + issueNumber
		return nil, err
	}
	err = stub.PutState("bank"+bankId, bankBytes)
	if err != nil {
		bank.Number = bank.Number - issueNumber
		centerBankBytes.RestNubmer = centerBank.RestNubmer + issueNumber
		return nil, errors.New("PutState Error" + err.Error())
	}

	transaction := Transaction{FromType:0,FromID:0,ToType:1,ToId:args[0],Time:time.Now().Unix(),Number:issueNumber,ID:transactionId}
	tsBytes,err = json.Marshal(&transaction)
	if(err!=nil){
		return nil,err
	}
	err = stub.PutState("transaction"+strconv.ItoA(transactionNo),tsBytes)
	if err!= nil{
		return nil, errors.New("PutState Error" + err.Error())
	}
	transactionNo = transactionNo +1
	return tsBytes, nil
}

func (t *SimpleChaincode) issueCoinToCp(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var company Company
	var bank Bank
	var bankId string
	var companyId string
	var issueNumber int
	var cpBytes []byte

	bankId = args[0]
	companyId = args[1]
	issueNumber = strconv.Itoa(args[2])

	bankBytes,err := stub.GetState("bank"+bankId)
	if err != nil {
		fmt.Println("Error retrieving cbBytes")
	}
	err = json.Unmarshal(bankBytes, &bank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	if bank.RestNubmer<issueNumber{
		return nil,errors.New("Not enough money")	
	}

	cpBytes,err = stub.GetState("company"+companyId)
	if err != nil {
		fmt.Println("Error retrieving cpBytes")
	}
	err = json.Unmarshal(cpBytes, &company)
	if err != nil {
		fmt.Println("Error unmarshalling company")
	}

	bank.RestNubmer = bank.RestNubmer - issueNumber
	company.Number = company.Number + issueNumber

	bankBytes, err = json.Marshal(&bank)
	if err != nil {
		bank.Number = bank.Number + issueNumber
		company.Number = company.Number - issueNumber
		return nil, err
	}
	err = stub.PutState("bank"+bankId, cbBytes)
	if err != nil {
		bank.Number = bank.Number + issueNumber
		company.Number = company.Number - issueNumber
		return nil, errors.New("PutState Error" + err.Error())
	}

	cpBytes, err = json.Marshal(&company)
	if err != nil {
		//这里的代码逻辑有问题，需要撤销对centerBank的修改
		bank.Number = bank.Number + issueNumber
		company.Number = company.Number - issueNumber
		return nil, err
	}
	err = stub.PutState("company"+companyId, cpBytes)
	if err != nil {
		bank.Number = bank.Number + issueNumber
		company.Number = company.Number - issueNumber
		return nil, errors.New("PutState Error" + err.Error())
	}

	transaction := Transaction{FromType:1,FromID:strconv.Atoi(args[0]),ToType:1,ToId:args[0],Time:time.Now().Unix(),Number:issueNumber,ID:transactionId}
	tsBytes,err = json.Marshal(&transaction)
	if(err!=nil){
		return nil,err
	}
	err = stub.PutState("transaction"+strconv.ItoA(transactionNo),tsBytes)
	if err!= nil{
		return nil, errors.New("PutState Error" + err.Error())
	}
	transactionNo = transactionNo +1
	return tsBytes, nil
}

func (t *SimpleChaincode) transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var cpFrom Company
	var cpTo Company
	var cpFromId string
	var cpToId string
	var issueNumber int
	var cpFromBytes []byte
	var cpToBytes []byte

	cpFromId = args[0]
	cpToId = args[1]
	issueNumber,err := strconv.Atoi(args[2])

	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	cpFromBytes,err := stub.GetState("company"+cpFromId)
	if err != nil {
		fmt.Println("Error retrieving cbBytes")
	}
	err = json.Unmarshal(cpFromBytes, &cpFrom)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	if cpFrom.RestNubmer<issueNumber{
		return nil,errors.New("Not enough money")	
	}

	cpToBytes,err = stub.GetState("company"+cpToId)
	if err != nil {
		fmt.Println("Error retrieving cpBytes")
	}
	err = json.Unmarshal(cpToBytes, &cpTo)
	if err != nil {
		fmt.Println("Error unmarshalling company")
	}

	cpFrom.RestNubmer = cpFrom.RestNubmer - issueNumber
	cpTo.Number = cpTo.Number + issueNumber

	cpFromBytes, err = json.Marshal(&cpFrom)
	if err != nil {
		cpFrom.RestNubmer = cpFrom.RestNubmer - issueNumber
		cpTo.Number = cpTo.Number + issueNumber
		return nil, err
	}
	err = stub.PutState("company"+cpFromId, cpFromBytes)
	if err != nil {
		cpFrom.RestNubmer = cpFrom.RestNubmer - issueNumber
		cpTo.Number = cpTo.Number + issueNumber
		return nil, errors.New("PutState Error" + err.Error())
	}

	cpToBytes, err = json.Marshal(&cpTo)
	if err != nil {
		//这里的代码逻辑有问题，需要撤销对centerBank的修改
		cpFrom.RestNubmer = cpFrom.RestNubmer - issueNumber
		cpTo.Number = cpTo.Number + issueNumber
		return nil, err
	}
	err = stub.PutState("company"+cpToId, cpToBytes)
	if err != nil {
		cpFrom.RestNubmer = cpFrom.RestNubmer - issueNumber
		cpTo.Number = cpTo.Number + issueNumber
		return nil, errors.New("PutState Error" + err.Error())
	}

	transaction := Transaction{FromType:2,FromID:cpFromId,ToType:2,ToId:cpToId,Time:time.Now().Unix(),Number:issueNumber,ID:transactionId}
	tsBytes,err = json.Marshal(&transaction)
	if(err!=nil){
		return nil,err
	}
	err = stub.PutState("transaction"+strconv.ItoA(transactionNo),tsBytes)
	if err!= nil{
		return nil, errors.New("PutState Error" + err.Error())
	}
	transactionNo = transactionNo +1
	return tsBytes, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if len(args) == 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. name of the key and value to set")
	}

	if function == "getCenterBank" {
		cbBytes, err := getCenterBank(stub, args)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return cbBytes, nil
	} else if function == "getBankById" {
		bankBytes, err := getBankById(stub, args)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return bankBytes, nil
	} else if function == "getCompanyById" {
		cpBytes, err := getCompanyById(stub, args)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return cpBytes, nil
	} else if function == "getTransactionById" {
		tsBytes, err := getTransactionById(stub, args)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return tsBytes, nil
	} else if function == "getBanks" {
		bankBytes, err := getBanks(stub)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return bankBytes, nil
	} else if function == "getCompanys" {
		cpBytes, err := getCompanys(stub)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return cpBytes, nil
	} else if function == "getTransactions" {
		tsBytes, err := getTransactions(stub)
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return tsBytes, nil
	}
}

func getCenterBank(stub *shim.ChaincodeStub) (CenterBank, []byte,error) {
	var centerBank CenterBank
	cbBytes, err := stub.GetState("centerBank")
	if err != nil {
		fmt.Println("Error retrieving cbBytes")
	}
	err = json.Unmarshal(cbBytes, &centerBank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	return centerBank,cbBytes, nil
}

func getCompanyById(stub *shim.ChaincodeStub, id string) (Company,[]byte, error) {
	var company Company
	cpBytes,err := stub.GetState("company"+id)
	if err != nil {
		fmt.Println("Error retrieving cpBytes")
	}
	err = json.Unmarshal(cpBytes, &company)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	return company,cpBytes, nil
}

func getBankById(stub *shim.ChaincodeStub, id string) (Bank, []byte,error) {
	var bank Bank
	cbBytes,err := stub.GetState("bank"+id)
	if err != nil {
		fmt.Println("Error retrieving cpBytes")
	}
	err = json.Unmarshal(cbBytes, &bank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	return bank,cbBytes, nil
}

func getTransactionById(stub *shim.ChaincodeStub, id string) (Transaction,[]byte, error) {
	var transaction Transaction
	tsBytes,err := stub.GetState("transaction"+id)
	if err != nil {
		fmt.Println("Error retrieving cpBytes")
	}
	err = json.Unmarshal(tsBytes, &transaction)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	return transaction,tsBytes, nil
}

func getBanks(stub *shim.ChaincodeStub) ([]Bank, error) {
	//需要看一下golang for循环与数组操作写法
	var banks []Bank

	if bankNo<=10 {
		i:=0
		for i<bankNo {
			bankBytes ,err = stub.GetState("bank"+strconv.Itoa(i))
			var bank Bank
			err = json.Unmarshal(bankBytes, &bank)
			if err != nil {
			return nil, errors.New("Error retrieving bank")
			banks = append(banks,bank)
			}
		}
	} else{
		i:=0
		for i<10{
			bankBytes ,err = stub.GetState("bank"+strconv.Itoa(i))
			var bank Bank
			err = json.Unmarshal(bankBytes, &bank)
			if err != nil {
			return nil, errors.New("Error retrieving bank")
			banks = append(banks,bank)
		}
	}

	return banks, nil
}

func getCompanys(stub *shim.ChaincodeStub) ([]Company, error) {
	var companys []Company

	if cpNo<=10 {
		i:=0
		for i<bankNo {
			cpBytes ,err := stub.GetState("company"+strconv.Itoa(i))
			var company Company
			err = json.Unmarshal(cpBytes, &company)
			if err != nil {
			return nil, errors.New("Error retrieving company")
			companys = append(companys,company)
			}
		}
	} else{
		i:=0
		for i<10{
			cpBytes ,err := stub.GetState("company"+strconv.Itoa(i))
			var company Company
			err = json.Unmarshal(cpBytes, &company)
			if err != nil {
			return nil, errors.New("Error retrieving company")
			companys = append(companys,company)
		}
	}

	return companys, nil
}

func getTransactions(stub *shim.ChaincodeStub) ([]Transaction, error) {
	var transactions []Transaction

	if transactionNo<=10 {
		i:=0
		for i<transactionNo {
			tsBytes ,err := stub.GetState("transaction"+strconv.Itoa(i))
			var transaction Transaction
			err = json.Unmarshal(tsBytes, &transaction)
			if err != nil {
			return nil, errors.New("Error retrieving transaction")
			transactions = append(transactions,transaction)
			}
		}
	} else{
		i:=0
		for i<10{
			tsBytes ,err := stub.GetState("transaction"+strconv.Itoa(i))
			var transaction Transaction
			err = json.Unmarshal(tsBytes, &transaction)
			if err != nil {
			return nil, errors.New("Error retrieving transaction")
			transactions = append(transactions,transaction)
		}
	}

	return transactions, nil
}

func writeCenterBank(stub *shim.ChaincodeStub,centerBank CenterBank) (error) {
	cbBytes, err := json.Marshal(&centerBank)
	if err != nil {
		return err
	}
	err = stub.PutState("centerBank", cbBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func writeBank(stub *shim.ChaincodeStub,bank Bank) (error) {
	var bankId string
	bankBytes, err := json.Marshal(&bank)
	if err != nil {
		return err
	}
	bankId,err = strconv.Itoa(bank.ID)
	if err!= nil{
		return errors.new("want Integer number")
	}
	err = stub.PutState("bank"+bankId, bankBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func writeCompany(stub *shim.ChaincodeStub,company Company) (error) {
	var companyId string
	cpBytes, err := json.Marshal(&company)
	if err != nil {
		return err
	}
	companyId,err = strconv.Itoa(company.ID)
	if err!= nil{
		return errors.new("want Integer number")
	}
	err = stub.PutState("company"+companyId, cpBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func writeTransaction(stub *shim.ChaincodeStub,transaction Transaction) (error) {
	var tsId string
	tsBytes, err := json.Marshal(&transaction)
	if err != nil {
		return err
	}
	tsId,err = strconv.Itoa(transaction.ID)
	if err!= nil{
		return errors.new("want Integer number")
	}
	err = stub.PutState("transaction"+tsId, tsBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}