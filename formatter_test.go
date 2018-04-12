package formatter

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestFormatter(t *testing.T) {
	var log = logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
	log.Formatter = &Formatter{Debug: true}

	log.Info("This is a info")
	log.Error("This is a error")
	log.Warnf("This is a warnf")
	log.Warningf("This is a warningf")
}
