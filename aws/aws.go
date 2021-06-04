package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

var Sesion *session.Session

func SetSession() error {
	var err error
	if Sesion != nil {
		return nil
	}
	Sesion, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	return err
}
