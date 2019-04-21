package base

import (
	"html/template"
	"path"
	"regexp"
	"strings"
)

// packr version

// getAllFiles 从指定目录下读取文件
func getAllFilesBox(folder string) []fileInfo {
	sugar := Sugar()

	sugar.Debugw("scan files", "folder", folder)

	all := Box.List()
	filered := []string{}
	for _, a := range all {
		if strings.HasPrefix(a, folder) {
			filered = append(filered, a)
		}
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
	for _, i := range filered {
		extension := path.Ext(i)
		goon, ok := allowed[extension]
		if !ok || !goon {
			sugar.Debugw("ignore file", "file", i, "ext", extension)
			continue
		}

		// 去掉后缀 commands/matrix.sql  --> {matrix, commands/matrix.sql, sql}
		basename := path.Base(i)
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
			FullPath:  i,
			Extension: extension,
		}

		files = append(files, entry)
	}

	return files
}

// 1. 读取每个文件，按照约定好的格式，拆分成不同的 command
// 2. 对 command 进行编译

// parseBatch 批量处理文件
func parseBatchBox(files []fileInfo) []commandInfo {
	sugar := Sugar()

	var commands []commandInfo
	for _, f := range files {
		cmds, err := parseFileBox(f)
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
func parseFileBox(info fileInfo) ([]commandInfo, error) {
	var err error
	sugar := Sugar()

	filename := info.FullPath
	s, err := Box.FindString(filename)
	if err != nil {
		sugar.Errorw("ReadLines failed:", "file", filename, "err", err)

		return nil, err
	}
	lines := strings.Split(s, "\n")

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
