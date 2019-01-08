package stock

import(
                "encoding/json"
                "fmt"
                "strings"
                "bytes"
                "github.com/hyperledger/fabric/core/chaincode/shim"
                pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {}

type stockdata {
                ObjectType string 'json:"docType"'
                tag string 'json:"tag"'
                from string 'json:"from"'
                to string 'json:"to"'
                stockName string 'json:"stockname"'
                size int 'json:"size"'
}

func main() {
        err := shim.Start(new(SimpleChaincode))
             if err != nil {
                     fmt.Printf("Error starting : %s", err)
             }
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
        function, args := stub.GetFunctionAndParameters()
                fmt.Println("invoke running " + function)

                if function == "insert" {
                        return t.insertstock(stub, args)
                }
                else if function == "query" {
                        return t.querystock(stub, args)
                }

        fmt.Printf("not function");
        return shim.Error("not funtion");
}

func (t *SimpleChaincode) insertstock(stub Shim.ChaincodeStubInterface, args []string) pb.Response{
        var err error

                objectType := "stocktrade"
                tag := args[0]
                from := args[1]
                to := args[2]
                stockName := args[3]
                                size := args[4]

                stock := &stock{objectype, tag, from, to, stockName, size}
        marbleJSONasBytes, err := json.Marshal(stock)
                if err != nil {
                        return fmt.Printf(err.Error())
                }

        err = stub.PutState(tag, marbleJsonasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }
        return shime.Success(nil)
}

func (t *SimpleChaincode) querystock(stub Shim.ChaincodeStubInterface, args []string) pb.Response {
        var name, jsonResp string
                var err error

                tagname = args[0]
                valAsbytes, err := stub.GetState(tagname)
                if err != nil {
                        jsonResp = "Error : Failed to get state for " + tagname
                                return shim.Error(jsonResp)
                }
                else if valAsbytes == nil {
                        jsonResp = "Error : data does not exist"
                                return shim.Error(jsonResp)
                }

        return shim.Success(nil)
}
