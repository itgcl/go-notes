package main

import (
	"fmt"
	"time"

	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
)

// js处理msg payload和元数据
func main() {
	config := rulego.NewConfig()

	metaData := types.NewMetadata()
	metaData.PutValue("productType", "test01")

	// js处理
	ruleEngine, err := rulego.New("rule01", []byte(chainJsonFile1), rulego.WithConfig(config))
	if err != nil {
		panic(err)
	}

	msg1 := types.NewMsg(0, "TEST_MSG_TYPE333", types.JSON, metaData, "{\"temperature\":41}")

	ruleEngine.OnMsgWithOptions(msg1, types.WithEndFunc(func(msg types.RuleMsg, err error) {
		fmt.Println("msg1处理结果=====")
		// 得到规则链处理结果
		fmt.Println(msg, err)
		fmt.Printf("%+v\n", msg)
	}))

	////js处理后，并调用http推送
	//ruleEngine2, err := rulego.New("rule02", []byte(chainJsonFile2), rulego.WithConfig(config))
	//if err != nil {
	//	panic(err)
	//}
	//
	//msg2 := types.NewMsg(0, "TEST_MSG_TYPE2", types.JSON, metaData, "{\"temperature\":30}")
	//ruleEngine2.OnMsgWithOptions(msg2, types.WithEndFunc(func(msg types.RuleMsg, err error) {
	//	fmt.Println("msg2处理结果=====")
	//	//得到规则链处理结果
	//	//因为推送的url:http://192.168.136.26:9099/api/msg 是无效url，所以会返回超时错误
	//	fmt.Println(msg, err)
	//}))
	//
	time.Sleep(time.Second * 3)
}

var chainJsonFile1 = `
{
  "ruleChain": {
	"id":"rule01",
    "name": "测试规则链",
    "root": true
  },
  "metadata": {
    "nodes": [
       {
        "id": "s1",
        "type": "jsTransform",
        "name": "转换",
        "debugMode": true,
        "configuration": {
          "jsScript": "metadata['name']='test01';\n metadata['index']=11;\n msg['addField']='addValue1';\n msg['ww']='ss'; return {'msg':msg,'metadata':metadata,'msgType':'222'};"
        }
      }
    ],
    "connections": [
    ],
    "ruleChainConnections": null
  }
}
`

var chainJsonFile2 = `
{
  "ruleChain": {
    "id":"rule02",
    "name": "测试规则链",
    "root": true
  },
  "metadata": {
    "nodes": [
       {
        "id": "s1",
        "type": "jsTransform",
        "name": "转换",
        "debugMode": true,
        "configuration": {
          "jsScript": "metadata['name']='test02';\n metadata['index']=22;\n msg['addField']='addValue2'; return {'msg':msg,'metadata':metadata,'msgType':msgType};"
        }
      },
      {
        "id": "s2",
        "type": "restApiCall",
        "name": "推送数据",
        "debugMode": true,
        "configuration": {
          "restEndpointUrlPattern": "http://192.168.136.26:9099/api/msg",
          "requestMethod": "POST",
          "maxParallelRequestsCount": 200
        }
      }
    ],
    "connections": [
      {
        "fromId": "s1",
        "toId": "s2",
        "type": "Success"
      }
    ],
    "ruleChainConnections": null
  }
}
`
