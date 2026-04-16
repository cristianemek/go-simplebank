package mail

import (
	"testing"

	"github.com/cristianemek/go-simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test Email"
	content := `
	<h1>Hello World</h1>
	<p>This is a test message</p>
	`
	to := []string{"emecuatroceroochoca@gmail.com"}
	attachFiles := []string{"../README.MD"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)

	require.NoError(t, err)

}
