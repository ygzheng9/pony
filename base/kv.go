package base

import (
	"html/template"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"sync"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

// KeyValue 目前是内存模式，以后可以扩展使用 _kv-_db
type KeyValue interface {
	LoadCommands() error
	GetCommand(key string, data interface{}) string
}

// 全局变量，仅此一份，必须是 pointer
var _kv KeyValue
var kvOnce sync.Once

// KvCache 全局
func KvCache() KeyValue {
	kvOnce.Do(func() {
		// 这里是策略，目前是 MemKV，可以换成别的 KV
		_kv = &MemKV{}

		err := _kv.LoadCommands()
		if err != nil {
			// err 在下一层已经处理过
			return
		}
	})

	return _kv
}

// 读取指定目录下的文件，扩展名有限制
// 读取每个文件中的 command
// 把每个 command 的 content 都编译成 text/template
// 使用时，通过 key 并传入 模板参数 data

// fileInfo 每一个文件信息
type fileInfo struct {
	// 去除扩展名后的文件名
	Name string
	// 文件的完整路径，包括目录和扩展名
	FullPath string
	// 扩展名，包含 .
	Extension string
}

// getAllFiles 从指定目录下读取文件
func getAllFiles(folder string) ([]fileInfo, error) {
	sugar := Sugar()

	sugar.Debugw("scan files", "folder", folder)

	fileInfos, err := ioutil.ReadDir(folder)
	if err != nil {
		sugar.Errorw("can not get file list", "folder", folder, "err", err.Error())
		return nil, err
	}

	// 处理的文件后缀名
	allowed := map[string]bool{
		".txt":  true,
		".sql":  true,
		".cmd":  true,
		".tmpl": true,
		".html": true,
	}

	var files []fileInfo
	for _, i := range fileInfos {
		n := i.Name()

		extension := path.Ext(n)
		goon, ok := allowed[extension]
		if !ok || !goon {
			sugar.Debugw("ignore file", "file", n, "ext", extension)
			continue
		}

		// 去掉后缀
		basename := path.Base(n)
		var name = basename[0 : len(basename)-len(extension)]

		// 检查是否有相同的文件名
		for _, c := range files {
			if name == c.Name {
				sugar.Errorw("duplicate file.", "name", name, "ext1", extension, "ext2", c.Extension)
			}
		}

		// 新文件，加入到列表中
		entry := fileInfo{
			Name:      name,
			FullPath:  folder + "/" + n,
			Extension: path.Ext(n),
		}

		files = append(files, entry)
	}

	return files, nil
}

// commandInfo 每一个命令
type commandInfo struct {
	// 唯一的 key
	Key string
	// 原始内容
	Content string
	// 编译后的信息
	Tmpl *template.Template
}

// 1. 读取每个文件，按照约定好的格式，拆分成不同的 command
// 2. 对 command 进行编译

// parseBatch 批量处理文件
func parseBatch(files []fileInfo) []commandInfo {
	sugar := Sugar()

	var commands []commandInfo
	for _, f := range files {
		cmds, err := parseFile(f)
		if err != nil {
			sugar.Errorw("batch err", "file", f.FullPath, "err", err)
			continue
		}

		// 两个 slice 的合并
		commands = append(commands, cmds...)
	}

	return commands
}

// parseFile 读取一个文件
func parseFile(info fileInfo) ([]commandInfo, error) {
	sugar := Sugar()
	// sugar.Debugw("parseFile", "file", info)

	filename := info.FullPath
	// 一次性读取文件中所有行
	lines, err := ReadLines(filename)
	if err != nil {
		sugar.Errorw("ReadLines failed:", "file", filename, "err", err)

		return nil, err
	}

	// 逐行遍历每一行，拆分成不同的 command；
	var commands []commandInfo

	// key 是 commandName
	markers := make(map[string]bool)
	cmdName := ""
	var content []string

	// 生成一条命令
	tickCommand := func() {
		if len(cmdName) > 0 {
			if _, ok := markers[cmdName]; ok {
				sugar.Errorw("command duplicate", "file", filename, "command", cmdName)
			} else {
				markers[cmdName] = true
				// sugar.Debugw("command got", "file", filename, "command", cmdName)

				// 生成新命令
				entry := commandInfo{
					Key:     info.Name + "." + cmdName,
					Content: strings.Join(content, "\n"),
				}
				entry.Tmpl = template.Must(template.New(entry.Key).Parse(entry.Content))
				commands = append(commands, entry)
			}
		}
	}

	// 格式固定，格式为：
	// cmd: getAllUser
	r := regexp.MustCompile(`cmd: ([[:alnum:]]+)`)

	last := len(lines) - 1
	for idx, l := range lines {
		// 最后一行
		if idx == last {
			tickCommand()
			break
		}

		// 不是最后一行，但是匹配到命令行，也即：发现了新的命令
		if r.MatchString(l) {
			// 有之前的命令
			tickCommand()

			// 新的命令开始, 设置命令名，清空内容
			match := r.FindStringSubmatch(l)
			// match[0] 是字符串本身，match[1] 是匹配的内容
			cmdName = match[1]
			content = []string{}

			continue
		}

		// 具体的内容
		// -- // ## 开头的是注释，忽略掉；空行也忽略掉
		if len(l) == 0 || strings.HasPrefix(l, "--") ||
			strings.HasPrefix(l, "//") || strings.HasPrefix(l, "##") {
			continue
		}

		// 每个 command 的内容有多行
		content = append(content, l)
	}

	return commands, nil
}
