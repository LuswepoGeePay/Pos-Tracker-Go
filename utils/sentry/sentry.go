package sentry

import (
	"fmt"
	"pos-master/utils"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func InitSentry() error {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://402b2ae744431da248fca4975f4f7d4a@o4509478949683201.ingest.de.sentry.io/4509484617039952",
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
		return utils.CapitalizeError(fmt.Sprintf("error: %v", err))
	}
	return nil
}

func SentryLogger(c *gin.Context, err error) {

	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.CaptureException(err)
	}

}
