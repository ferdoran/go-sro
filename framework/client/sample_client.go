package client

func StartSampleClient() {
	client := NewClient("127.0.0.1", 15779, "SR_Client")
	client.Connect()
}
