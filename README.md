GCM
===

Google Cloud Message implementation in Go

	go get github.com/bashtian/gcm

Sending a message

	s := gcm.NewSender("API_KEY")

	m := gcm.NewMessage("DEVICE_REGISTRAION_ID")
	m.Add("text", "Hello World")

	res, err := s.SendNoRetry(m)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v", res)
	}
	
### Documentation
http://godoc.org/github.com/bashtian/gcm

http://developer.android.com/google/gcm/
