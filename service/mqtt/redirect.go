package mqtt

type MQTTRedirectClient struct {
	Topic string
	Broker string
	User string
	Password string
	ClientID string
}

func (client *MQTTRedirectClient)SendMessage(sessionid,msg string){
	//初始化mqtt客户端
	mqttClient:=&MQTTClient{
		Broker:client.Broker,
		User:client.User,
		Password:client.Password,
		ClientID:sessionid,
	}
	mqttClient.Init()
	mqttClient.Publish(client.Topic,msg)
	mqttClient.Close()
}