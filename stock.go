package donation

import(
                "encoding/json"
                "fmt"                
                "strconv"
                "github.com/hyperledger/fabric/core/chaincode/shim"
                pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {}

type donationlist struct {
         ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	       tag       string `json:"tag"`    //the fieldtags are needed to keep case from bouncing around
	       name      string `json:"name"`
	       money       int    `json:"money"`
         date      string `json:"date"`
}

func main() {
        err := shim.Start(new(SimpleChaincode))
             if err != nil {
                     fmt.Printf("Error starting : %s", err)
             }
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	      return shim.Success(nil)
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
        function, args := stub.GetFunctionAndParameters()
                fmt.Println("invoke running " + function)

                if function == "insert" {
                        return t.insert(stub, args)
                } else if function == "query" {
                        return t.query(stub, args)
                }

        fmt.Printf("not function");
        return shim.Error("not funtion");
}

func (t *SimpleChaincode) insert(stub shim.ChaincodeStubInterface, args []string) pb.Response{
        var err error

                objectType := "receive"
                tag := args[0]
                name := args[1]
                money, _ := strconv.Atoi(args[2])
                date := args[3]
            

                donationlist := &donationlist{objectype, tag, name, money, date}
                marbleJSONasBytes, err := json.Marshal(stock)
                if err != nil {
                       return shim.Error(err.Error())
                }

        err = stub.PutState(tag, marbleJSONasBytes)
                if err != nil {
                        return shim.Error(err.Error())
                }
        return shim.Success(nil)
}

func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var tagname, jsonResp string
                var err error

                tagname = args[0]
                valAsbytes, err := stub.GetState(tagname)
                if err != nil {
                        jsonResp = "Error : Failed to get state for " + tagname
                                return shim.Error(jsonResp)
                } else if valAsbytes == nil {
                        jsonResp = "Error : data does not exist"
                                return shim.Error(jsonResp)
                }

        return shim.Success(nil)
}
