//yeasy/hyperledger-fabric-couchdb
package main

import(
                "encoding/json"
                "fmt"                
                "strconv"
		"bytes"
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

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
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
            

                donationlist := &donationlist{objectType, tag, name, money, date}
                marbleJSONasBytes, err := json.Marshal(donationlist)
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
        var queryString, jsonResp string
                var err error

                queryString = args[0]
                queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
        } 

        return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}
