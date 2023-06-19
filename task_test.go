package cronjob

import (
	"cronjob/mysql"
	"testing"
)

func TestDb(t *testing.T) {

	mysql.NewClient()
}
