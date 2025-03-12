package mail

import (
	"testing"

	"github.com/flexGURU/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {

	config, err := utils.LoadConfig("..")
	require.NoError(t, err )

	mailSender := NewGmailSender(config.EmailSendName, config.From_Email, config.EamilPassword)

	subject := "TEST EMAIL"
	content := `
	<h1>TEST EMAIL</h1>
	<p>hello there</p>
	`
	to := []string{"mukunajohn329@gmail.com"}

	err = mailSender.SendEmail(subject, content, to, nil, nil)
	require.NoError(t, err )
	


}