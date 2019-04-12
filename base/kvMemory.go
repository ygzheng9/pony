package base

import (
	"strings"

	"github.com/gobuffalo/envy"
)

// MemKV 预编译模板
type MemKV struct {
	allCommands []commandInfo
}

// LoadCommands 从指定目录，加载文件
func (k *MemKV) LoadCommands() error {
	sugar := Sugar()

	// envy.Load(".env")
	tmplDir := envy.Get("CommandDir", "./config/commands")
	// sugar.Info(tmplDir)

	files, err := getAllFiles(tmplDir)
	if err != nil {
		sugar.Errorw("load command dir failed", "command dir", tmplDir, "err", err)
		return err
	}

	k.allCommands = parseBatch(files)

	return nil
}

// GetCommand 从全局的 map 中取得 key 对应的 value，都是 string
func (k *MemKV) GetCommand(key string, data interface{}) string {
	sugar := Sugar()
	env := envy.Get("GO_ENV", "development")

	if env == "development" {
		// 在 debug 下，每次都重新加载所有文件
		err := k.LoadCommands()
		if err != nil {
			sugar.Errorf("KV failed. %+v", err)
			return ""
		}
		sugar.Debugf("allCommands: %d", len(k.allCommands))
	}

	for _, c := range k.allCommands {
		if c.Key == key {
			var sb strings.Builder

			err := c.Tmpl.Execute(&sb, data)
			if err != nil {
				sugar.Errorw("execute err", "key", key, "data", data, "err", err.Error())
				return ""
			}

			result := sb.String()
			sugar.Debugw("template result", "key", key, "data", data)
			sugar.Debugf("content: \n%s", result)

			return result
		}
	}

	sugar.Warnf("can not find key: %s", key)
	return ""
}
