//19-01-16
//1. start data, end date 추가
//2. 쿼리 결과를 json으로반환하기
package main
  
import(
        "encoding/json"
        "fmt"
        //"reflect"
        "strconv"
        "bytes"
        "github.com/hyperledger/fabric/core/chaincode/shim"
        pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {}

type user struct {
        ObjectType string `json:"DocType"`
        ID string `json:"ID"`
        Name string `json:"Name"`
        Password string `json:"Password"`
        SocialNumber string `json:"SocialNumber"`
        Location string `json:"Location"`
} // user information

type uservoteresult struct {
        ID string `json:"ID"`
        Result []int `json:"Result"`
} // user's vote result

type vote struct {
        ObjectType string `json:"DocType"`
        Votename string `json:"Votename"`
        Question []string `json:"Question"`
        UserResult []uservoteresult `json:"UserResult"`
} //have question and its result per user

type tempvote struct {
        ObjectType string `json:"DocType"`
        Question []string `json:"Question"`
        UserResult []uservoteresult `json:UserReslult`
} // temp structur for append and Unmarshal()

func main () {
        err := shim.Start(new(SimpleChaincode))
        if err != nil {
                fmt.Printf("Error starting : %s", err)
        }
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
        fmt.Printf("Initialize\n")      
        return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        function, args := stub.GetFunctionAndParameters()

        if function == "insertUser" {
                return t.insertUser(stub, args)
        } else if function == "insertVote"{
                return t.insertVote(stub, args)
        } else if function == "insertVoteResult" {
                return t.insertVoteResult(stub, args)
        } else if function == "query" {
                return t.query(stub, args)
        } else if function == "delete" {
                return t.delete(stub, args)
        }

        fmt.Printf("not function");
        return shim.Error("not funtion");
}

func (t *SimpleChaincode) insertUser(stub shim.ChaincodeStubInterface, args []string) pb.Response{
        var err error

        fmt.Printf("insert...\n")
        
        //inser user information
        objectType := "user"
        id := args[0]
        name := args[1]
        password := args[2]
        SSN := args[3]
        location := args[4]

        //make structure and marshaling
        user := &user{objectType,id, name, password, SSN, location}
        userJSONasBytes, err := json.Marshal(user)
        if err != nil {
                fmt.Printf("error")
                return shim.Error(err.Error())
        }
        fmt.Println(string(userJSONasBytes))
        if err != nil {
                fmt.Printf("insert error 1\n")
                return shim.Error(err.Error())
        }

        //insert DB
        err = stub.PutState(id, userJSONasBytes)
        if err != nil {
                fmt.Printf("insert error 2\n")
                return shim.Error(err.Error())
        }

        return shim.Success(nil)
}

func (t *SimpleChaincode) insertVote (stub shim.ChaincodeStubInterface, args []string) pb.Response{
        //Admin insert vote question
        var err error
        votename := args[0]
        questionNum := len(args) - 1
        voteAsByte, err := stub.GetState(votename)
        if voteAsByte != nil {
                fmt.Println("vote already exist")
                return shim.Error("vote already exist")
        }

        //inser vote qustion
        votedata := vote{}
        votedata.ObjectType = "vote"
        votedata.Votename = args[0]
        votedata.Question = make([]string, questionNum)
        for i := 0; i < questionNum; i++ {
                votedata.Question[i] = args[i + 1]
        }
        fmt.Println(votedata)

        // marshaling
        voteAsJSONBytes, err := json.Marshal(votedata)
        if err != nil {
                fmt.Println("Marshal error")
                return shim.Error("Marshal error")
        }

        //insert DB
        err = stub.PutState(votename, voteAsJSONBytes)
        if err != nil {
                fmt.Println("DB insert error")
                shim.Error("DB insert error")
        }

        fmt.Println("vote question insert success")
        return shim.Success(nil)

}

func (t *SimpleChaincode) insertVoteResult (stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var queryString string
        var err error
        var temp tempvote

        //query key that check for vote question have existed
        queryString = "{\"selector\":{\"DocType\":\"vote\", \"Votename\":\"" + args[0] + "\"}}"
        existflag, err := stub.GetState(args[0])
        if existflag == nil {
                fmt.Println("this vote not exist")
                return shim.Error("this vote not exsit")
        }

        // prev query get key-vaule form
        resultIterator, err := stub.GetQueryResult(queryString)
        if err != nil {
                fmt.Println("get query error")
                return shim.Error("get query error")
        }

        response, err := resultIterator.Next()
        if err != nil {
                fmt.Println("not next")
                return shim.Error("not next")
        }

        //value unmarshaling to temp structure
        err = json.Unmarshal(response.Value, &temp)
        if err != nil {
                fmt.Println("unmarshal error")
                return shim.Error("unmarshal error")
        }
        fmt.Println("prev struct : ", temp)

        //make new vote structure and append present data
        votedata := vote{}
        questionNum := len(temp.Question)
        userNum := len(temp.UserResult)
        resultNum := len(args) - 2
        fmt.Println("questionNum : ", questionNum, "userNum : ", userNum, "restultNum : ", resultNum)
        votedata.ObjectType = temp.ObjectType
        votedata.Votename = args[0]
        fmt.Println("nomraml success")
        //votedata.Question = temp.Question
        votedata.Question = make([]string, questionNum)
        for i := 0; i < questionNum; i++ {
                votedata.Question[i] = temp.Question[i]
        }
        fmt.Println("insert Question suucess")
        votedata.UserResult = make([]uservoteresult, userNum + 1)
        //votedata.UserResult = temp.UserResult
        for i := 0; i < userNum; i++ {
                votedata.UserResult[i] = temp.UserResult[i]
        }
        fmt.Println("prev user data insert success")
        votedata.UserResult[userNum].ID = args[1]
        votedata.UserResult[userNum].Result = make([]int, resultNum)
        for i:=0; i < resultNum; i++ {
                votedata.UserResult[userNum].Result[i], _ = strconv.Atoi(args[i+2])
        }

        //new vote structure marshalring
        fmt.Println("now struct : ", votedata)
        userResultAsJSONBytes, err := json.Marshal(votedata)
        if err != nil {
                fmt.Println("marshal error")
                return shim.Error("marshal error")
        }

        //insert DB
        err = stub.PutState(args[0], userResultAsJSONBytes)
        if err != nil {
                fmt.Println("insert DB error")
                shim.Error("insert DB error")
        }

        fmt.Println("new data make and inser success")
        return shim.Success(nil)
}

func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var queryString string
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
                fmt.Printf("quert error 1\n")
                return nil, err
        }
        defer resultsIterator.Close()

        //make buffer -> string form
        buffer, err := constructQueryResponseFromIterator(resultsIterator)
        if err != nil {
                fmt.Printf("query error2\n")
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

func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var err error
        userName := args[0]
        fmt.Printf("delete user named : %s\n", userName)
        // maybe i have to connect with user database
        existflag, err := stub.GetState(userName)
        if existflag == nil {
                fmt.Printf("this user not exist")
                return shim.Error(err.Error())
        }
        if err != nil {
                fmt.Printf("GetState() error")
                return shim.Error(err.Error())
        }
        err = stub.DelState(userName)
        if err != nil {
                fmt.Printf("DelState() error")
                return shim.Error(err.Error())
        }
        fmt.Printf("deletion successed.\n")
        return shim.Success(nil)
}
