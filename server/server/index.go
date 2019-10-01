package server

import (
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
)

// 请求截图服务
func Screenshot(data map[string]interface{}) {
	bytesData, err := json.Marshal(data)
	if err != nil {
			fmt.Println(err.Error() )
			return
	}

	reader := bytes.NewReader(bytesData)

	resp, err := http.Post("http://localhost:5000/screenshot", "application/json", reader)
	if err != nil {
	  fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	// var resBody ResBody
	// if err := json.Unmarshal([]byte(string(body)), &resBody); err == nil {
	// 	fmt.Println("================json str 转struct==")
	// 	fmt.Println(resBody)
	// 	fmt.Println(resBody.Url)
	// }
}