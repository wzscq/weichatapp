package mqtt

type MQTTRedirectClient struct {
	Topic string
	Broker string
	User string
	Password string
	ClientID string
	mqttClient *MQTTClient
}

func (client *MQTTRedirectClient)SendMessage(sessionid,msg string){
	//初始化mqtt客户端
	if client.mqttClient == nil {
		client.mqttClient=&MQTTClient{
			Broker:client.Broker,
			User:client.User,
			Password:client.Password,
			ClientID:client.ClientID,
		}
		client.mqttClient.Init()
	}
	client.mqttClient.Publish(client.Topic,msg)
}