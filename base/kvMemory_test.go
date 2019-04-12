package base

import (
	"testing"
)

func TestKV(t *testing.T) {
	SetRoot()

	kv := KvCache()
	sql := kv.GetCommand("meta.mysqlColumns", nil)

	t.Log(sql)
}
