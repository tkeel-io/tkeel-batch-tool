package http

import (
	"fmt"
	//"time"

	//"crypto/hmac"
	//"crypto/sha256"
	//"encoding/base64"
	//"net/url"
	//"strconv"
	//"strings"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	NorthUrl        string
	IotUrl          string
	Token           string
	AccessKey       string
	SecretAccessKey string
)

func DoCreate(Url string, url string, verb string, queryMap map[string]interface{}, bodyData []byte) (map[string]interface{}, error) {

	client := &http.Client{}

	fmt.Printf("verb = %s \n", verb)
	fmt.Printf("url = %s \n", Url+url)
	fmt.Printf("\n")
	fmt.Printf("body = %s \n", string(bodyData))
	fmt.Printf("\n")

	reqest, err := http.NewRequest(verb, Url+url, bytes.NewReader(bodyData))
	if err != nil {
		fmt.Println("http error")
		return nil, err
	}
	//query
	if len(queryMap) > 0 {
		q := reqest.URL.Query()
		for k, v := range queryMap {
			q.Add(k, v.(string))
		}
		reqest.URL.RawQuery = q.Encode()
		/*q1 := reqest.URL.Query()
		  a,ok := q1["ids"]
		  if !ok {
		      fmt.Printf("param a does not exist\n");
		  } else {
		      fmt.Printf("param a value is [%s]\n", a);
		  }*/
	}
	//fmt.Println(reqest.URL)
	//times := strconv.FormatInt(time.Now().Unix(), 10)
	//reqest.Header.Add("Authorization", "QC"+AccessKey+":"+createSignature(verb, url, times, SecretAccessKey, string(reqest.URL.RawQuery), string(bodyData)))
	//reqest.Header.Add("Authorization", "QC" + "admin" + ":" + "admin")
	//reqest.Header.Add("x-iot-timestamp", times)
	//reqest.Header.Add("x-iot-signature-method", "HmacSHA256")
	//reqest.Header.Add("x-iot-signature-version", "1")

	//add header option
	reqest.Header.Add("Content-Type", "application/json")
	reqest.Header.Add("Authorization", "Bearer "+Token)

	resp, err := client.Do(reqest)
	if err != nil {
		fmt.Println("do err\n")
		return nil, err
	}

	resultMap, err := ParseResponse(resp)
	if err != nil {
		fmt.Println("response error \n")
	}
	resp.Body.Close()
	return resultMap, err
}

/*
func createSignature(verb string, Url string, times string, sk string, query string, bodyData string) string {

	//common
	stringToSign := verb + "\n" + Url + "\n" + "x-iot-signature-method:HmacSHA256" + "\n" + "x-iot-signature-version:1" + "\n" + "x-iot-timestamp:" + times

	//Canonicalized Query
	if len(query) != 0 {
		stringToSign += "\n" + query
	}

	//Canonicalized Body
	if len(bodyData) != 0 {
		stringToSign += "\n" + bodyData
	}

	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(stringToSign))

	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	signature = url.QueryEscape(signature)
	return signature
}
*/

func ParseResponse(response *http.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	/*fmt.Printf("all response = %s\n",response)
	  fmt.Printf("\n")
	  fmt.Printf("response.body = %s\n",string(body))*/
	if err == nil {
		err = json.Unmarshal(body, &result)
	}
	return result, err
}
