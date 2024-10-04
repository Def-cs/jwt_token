package mail

import "fmt"

var Connection *Mail

type Mail struct {
}

func NewMailConnection() {
	Connection = &Mail{}
}

func (m *Mail) SendWarning(ip, email string) error {
	fmt.Println("Отправлено сообщении о новом ip пользователя: " + ip + " на почту: " + email)
	return nil
}
