package main

import (
    // "crypto/rand"
    // "time"
    "fmt"
    "bytes"
    // "io"
    // "strings"
    "encoding/json"
    "strconv"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleTrans struct {
}


// Define the record structure using encoding/json library
type Transaction struct {
    Id string `json:"id"`
    Sender string `json:"sender"`
    Receiver string `json:"receiver"`
    Sender_currency string `josn:"sender_currency"`
    Receiver_currency string `josn:"receiver_currency"`
    Timestamp string `json:"timestamp"`
    Forex float64 `json:"forex"`
    Balance float64 `json:"balance"`

}
type User struct{
    Email string `json:"email"`
    INR float64 `json:"inr"`
    USD float64 `json:"usd"`
    AUD float64 `json:"aud"`
    DTX float64 `json:"dtx"`
}

func (t *SimpleTrans) Init(stub shim.ChaincodeStubInterface) pb.Response {
    return shim.Success(nil)
}

func (t *SimpleTrans) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    if function == "newTrans" {
        return t.newTrans(stub,args)
    } else if function == "queryAll" {
        return t.queryAll(stub)
    } else if function == "queryUser" {
        return t.queryUser(stub,args)
    } else if function == "newUser" {
        return t.newUser(stub,args)
    } else if function == "queryTrans"{
        return t.queryTrans(stub,args)
    }
    return shim.Error("Invalid invoke function name. Expecting \"newTrans\" or \"queryTrans\"  or \"queryUser\" or \"newUser\"")
}



// func to create new transaction record
func (t *SimpleTrans) newTrans(stub shim.ChaincodeStubInterface,args []string) pb.Response {
    var sender, receiver string    // Entities
    var senderVal, receiverVal float64 // Asset holdings
    var amount, forex float64          // Transaction value
    var err error

    if len(args) != 6 {
        return shim.Error("Expecting 6 arguments")
    }
    sender = args[0]
    receiver = args[1]
    // sender_currency = strings.ToUpper(args[2])
    // receiver_currency = strings.ToUpper(args[3])
    forex, err = strconv.ParseFloat(args[5], 64)

    // Get the state from the ledger
    senderAsBytes, err := stub.GetState(sender)
    if err != nil {
        return shim.Error("Failed to get state of senderUser")
    }
    if senderAsBytes == nil {
        return shim.Error("Entity not found")
    }
    senderUser := User{}
    json.Unmarshal(senderAsBytes, &senderUser)
    senderVal = senderUser.AUD
    if senderVal < amount{
        return shim.Error("You do not have sufficent Balance")
    }
    receiverAsBytes, err := stub.GetState(receiver)
    if err != nil {
        return shim.Error("Failed to get state of receiverUser")
    }
    if receiverAsBytes == nil {
        return shim.Error("Entity not found")
    }
    receiverUser := User{}
    json.Unmarshal(receiverAsBytes, &receiverUser)
    receiverVal = receiverUser.INR

    // Perform the execution
    amount, err = strconv.ParseFloat(args[4], 64)
    if err != nil {
        return shim.Error("Invalid transaction amount, expecting a Float value")
    }
    senderVal = senderVal - amount
    receiverVal = receiverVal + amount*forex
    fmt.Printf("senderVal = %d, receiverVal = %d\n", senderVal, receiverVal)

    // Write the state back to the ledger
    senderUser.AUD = senderVal
    senderAsBytes, _ = json.Marshal(senderUser)
    err = stub.PutState(sender, senderAsBytes)
    if err != nil {
        return shim.Error(err.Error())
    }
    receiverUser.INR = receiverVal
    receiverAsBytes, _ = json.Marshal(receiverUser)
    err = stub.PutState(receiver, receiverAsBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    //update daala ledger of sushmit
    daalaAudAsBytes, err := stub.GetState("sushmit@sgit.io")
    daalaAudUser := User{}
    json.Unmarshal(daalaAudAsBytes, &daalaAudUser)
    audVal := daalaAudUser.AUD
    audVal = audVal + amount
    dtxVal := daalaAudUser.DTX
    dtxVal -= amount   //1 DTX = 1 AUD
    daalaAudUser.DTX = dtxVal
    daalaAudUser.AUD = audVal
    daalaAudAsBytes, _ = json.Marshal(daalaAudUser)
    err = stub.PutState("sushmit@sgit.io", daalaAudAsBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    //update daala ledger of nitin
    daalaInrAsBytes, err := stub.GetState("nitin@sgit.io")
    daalaInrUser := User{}
    json.Unmarshal(daalaInrAsBytes, &daalaInrUser)
    inrVal := daalaInrUser.INR
    inrVal = inrVal - amount*forex
    dtxVal = daalaInrUser.DTX
    dtxVal = dtxVal + amount   //1 DTX = 1 AUD
    daalaInrUser.DTX = dtxVal
    daalaInrUser.INR = inrVal
    daalaInrAsBytes, _ = json.Marshal(daalaInrUser)
    err = stub.PutState("nitin@sgit.io", daalaInrAsBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    // write final transation ledger
    // id := make([]byte, 16)
    // n, err := io.ReadFull(rand.Reader, id)
    // if n != len(id) || err != nil {
    //     return shim.Error(err.Error())
    // }
    // timestamp := time.Now()
    // uuid := fmt.Sprintf("%s", id)
    // var transation = Transaction{Id: uuid, Sender:sender, Receiver: receiver, Sender_currency: sender_currency, Receiver_currency: receiver_currency, Timestamp: timestamp.String(), Forex: forex, Balance:amount}

    // // Recording data to the ledger
    // transationAsBytes, _ := json.Marshal(transation)
    // err = stub.PutState(sender, transationAsBytes)
    // if err != nil {
    //     return shim.Error(fmt.Sprintf("Failed to create new transation: %s", args[0]))
    // }

    return shim.Success(nil)


}

// func to create new transaction record
func (t *SimpleTrans) newUser(stub shim.ChaincodeStubInterface,args []string) pb.Response {
    if len(args) != 5 {
        return shim.Error("Expecting 5 arguments")
    }
    inr, _ := strconv.ParseFloat(args[1], 64)
    aud, _ := strconv.ParseFloat(args[2], 64)
    usd, _ := strconv.ParseFloat(args[3], 64)
    dtx, _ := strconv.ParseFloat(args[4], 64)
    var user = User{Email: args[0], INR: inr, AUD: aud, USD: usd, DTX: dtx}

    // Recording to the ledger
    userAsBytes, _ := json.Marshal(user)
    err := stub.PutState(args[0], userAsBytes)
    if err != nil {
        return shim.Error(fmt.Sprintf("Failed to create new user: %s", args[0]))
    }
    return shim.Success(nil)
}

// func to query All the existing value in chain
func (t *SimpleTrans) queryAll(stub shim.ChaincodeStubInterface) pb.Response {
    startKey := "000"
    endKey := "999"

    resultsIterator, err := stub.GetStateByRange(startKey, endKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }

        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("------->")
        buffer.WriteString("{Key : ")
        buffer.WriteString(queryResponse.Key)
        buffer.WriteString(" | ")

        buffer.WriteString("Record : ")
        buffer.WriteString(string(queryResponse.Value))
        buffer.WriteString("}")

        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")

    fmt.Printf("- queryAll:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

func (t *SimpleTrans) queryUser(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    userAsBytes, _ := APIstub.GetState(args[0])
    return shim.Success(userAsBytes)
}
func (t *SimpleTrans) queryTrans(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    transAsBytes, _ := APIstub.GetState(args[0])
    return shim.Success(transAsBytes)
}
func main() {
    err := shim.Start(new(SimpleTrans))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}