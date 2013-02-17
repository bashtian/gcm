gcm
===

Google Cloud Message implementation in Go

	import "github.com/bashtian/gcm"

Sending a message

	s := gcm.NewSender("API_KEY")

	m := gcm.NewMessage([]string{"DEVICE_REGISTRAION_ID"})
	m.Add("text", "Hello World")

	res, err := s.SendNoRetry(m)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v", res)
	}