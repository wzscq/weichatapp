package bill

import (
	"encoding/json"
	"weichatapp/mqtt"
)

type BillSender struct {
	Topic string
	Broker string
	User string
	Password string
	ClientID string
}

type BillRec struct {
	Sessionid string `json:"sessionid"`
	Amount int `json:"amount"`
}

func (client *BillSender)BillingRecord(sessionid string,amount int){
	billRec:=&BillRec{
		Sessionid:sessionid,
		Amount:amount,
	}

	billJson,_:=json.Marshal(billRec)

	client.SendBill(sessionid,string(billJson))
}

func (client *BillSender)SendBill(sessionid,msg string){
	//初始化mqtt客户端
	mqttClient:=&mqtt.MQTTClient{
		Broker:client.Broker,
		User:client.User,
		Password:client.Password,
		ClientID:sessionid,
	}
	mqttClient.Init()
	mqttClient.Publish(client.Topic,msg)
	mqttClient.Close()
}