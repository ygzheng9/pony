package grifts

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"
	"unicode/utf8"

	"code.sajari.com/docconv"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
	"github.com/yanyiwu/gojieba"

	"pony/base"
	"pony/models"
)

var _ = grift.Namespace("words", func() {
	grift.Add("info", func(c *grift.Context) error {
		fmt.Println(c.Args)
		return nil
	})

	grift.Desc("skip", "Show skip words")
	grift.Add("skip", func(c *grift.Context) error {
		skips := readSkipWords()
		fmt.Print(skips)
		return nil
	})

	grift.Desc("merge", "Show merge words")
	grift.Add("merge", func(c *grift.Context) error {
		merge := readMergeWords()
		fmt.Print(merge)
		return nil
	})

	grift.Desc("jieba", "Extract words from files")
	grift.Add("jieba", func(c *grift.Context) error {
		sugar := base.Sugar()

		sugar.Info("here dump words")

		dumpWords()

		sugar.Info("completed dump words")

		return nil
	})

})

func isStrSliceContainStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func isMapKeyMatchStr(mapStr map[string]string, key string) bool {
	for k := range mapStr {
		if k == key {
			return true
		}
	}
	return false
}

// readSkipWords get all skip words
func readSkipWords() []string {
	var err error
	sugar := base.Sugar()

	dir := envy.Get("WordDir", "./config/words")
	file := dir + "/skip.txt"
	lines, err := base.ReadLines(file)
	if err != nil {
		sugar.Errorw("failed to get skipword", "file", file)
		return []string{}
	}
	return lines
}

// readMergeWords get all mapping a -> b
func readMergeWords() map[string]string {
	var err error
	sugar := base.Sugar()

	dir := envy.Get("WordDir", "./config/words")
	file := dir + "/merge.txt"
	lines, err := base.ReadLines(file)
	if err != nil {
		sugar.Errorw("failed to get merge word", "file", file)
		return map[string]string{}
	}

	m := make(map[string]string)
	for _, l := range lines {
		s := strings.Split(l, " ")
		if len(s) != 2 {
			sugar.Warnw("merge", "file", file, "line", l)
			continue
		}
		m[s[0]] = s[1]
	}

	return m
}

// fileInfo 每一个文件信息
type fileInfo struct {
	// 去除扩展名后的文件名
	Name string
	// 文件的完整路径，包括目录和扩展名
	FullPath string
	// 扩展名，包含 .
	Extension string
}

func dumpWords() {
	var err error
	sugar := base.Sugar()

	// folder := `/Users/ygzheng/Documents/project/Support/圣象/itsp2/20.战略理解与现状诊断/10.访谈/12.访谈纪要&会议纪要`
	folder := envy.Get("MemoDir", "")
	sugar.Infow("dump", "folder", folder)

	fileInfos, err := ioutil.ReadDir(folder)
	if err != nil {
		sugar.Errorw("can not get file list", "folder", folder, "err", err.Error())
		return
	}

	// 处理的文件后缀名
	allowed := map[string]bool{
		".docx": true,
		".doc":  true,
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

	// drop words to make name simple
	replaceWords := map[string]string{
		"圣象项目_": "",
		"会议记录":  "",
		"访谈纪要":  "",
		"会议纪要":  "",
	}
	for i, a := range files {
		b := a.Name
		for k, v := range replaceWords {
			b = strings.Replace(b, k, v, 1)
		}

		c := strings.Split(b, "_")
		files[i].Name = c[0]
	}

	// new instance for jieba
	x := gojieba.NewJieba()
	defer x.Free()

	// 在处理过程中过滤掉无关的高频词
	excludes := readSkipWords()
	// 替换的词
	mergeData := readMergeWords()

	// read from file and split and get word count
	var allWords []fileWords
	for _, i := range files {
		fmt.Println(i.FullPath)
		entry := fileWords{
			Name:  i.Name,
			Words: parseFile(i.FullPath, x, excludes, mergeData),
		}
		allWords = append(allWords, entry)
	}

	// save to DB
	db := models.DB
	db.Transaction(func(tx *pop.Connection) error {
		// delete all
		tx.Store.Exec("truncate table words; ")

		// each doc
		for _, i := range allWords {
			// each word
			for _, w := range i.Words {
				m := models.Word{
					DocName: i.Name,
					Word:    w.Key,
					Count:   w.Value,
				}
				verr, err := tx.ValidateAndSave(&m)
				if err != nil {
					sugar.Errorw("save failed", "doc", i.Name, "word", w, "err", err)
					return errors.Wrap(err, "save failed")
				}

				if len(verr.Errors) > 0 {
					sugar.Errorw("save err", "doc", i.Name, "word", w, "err", verr.Errors)
					return errors.Wrap(err, "save with error")
				}
			}

			sugar.Debugw("save complete", "doc", i.Name, "num", len(i.Words))
		}
		return nil
	})

	return
}

// kv for (word, count)
type kv struct {
	Key   string
	Value int
}

// fileWords for fileName with (word, count)
type fileWords struct {
	Name  string
	Words []kv
}

// getText from file
func getText(file string) string {
	sugar := base.Sugar()

	// get txt from docx
	res, err := docconv.ConvertPath(file)
	if err != nil {
		sugar.Errorw("extract text failed", "file", file, "err", err)
		return ""
	}
	// fmt.Println(res.Body)
	// fmt.Println(res.Meta)

	return res.Body
}

// parseFile get text from file and split by jieba
func parseFile(file string, eng *gojieba.Jieba, excludes []string, mergeData map[string]string) []kv {
	// get txt from docx
	input := getText(file)

	// jieba
	var words []string
	var ss []kv
	counts := make(map[string]int)

	useHmm := true
	words = eng.CutForSearch(input, useHmm)

	var rword string
	for _, word := range words {
		if utf8.RuneCountInString(word) == 1 {
			continue
		} else if isStrSliceContainStr(excludes, word) {
			continue
		} else if isMapKeyMatchStr(mergeData, word) {
			rword = mergeData[word]
		} else {
			rword = word
		}
		if i, ok := counts[rword]; ok == false {
			counts[rword] = 1
		} else {
			counts[rword] = i + 1
		}
	}

	// 把字典转换为列表
	for k, v := range counts {
		ss = append(ss, kv{k, v})
	}
	return ss
}

// rawSplit demo for jieba
func rawSplit(file string) []kv {
	// get txt from docx
	s := getText(file)

	// jieba
	var words []string
	var ss []kv
	counts := make(map[string]int)

	// 在处理过程中过滤掉无关的高频词
	excludes := []string{"什么", "一个", "我们", "那里"}

	// 替换的词
	mergeData := make(map[string]string)
	mergeData["老太太"] = "贾母"
	mergeData["老太"] = "贾母"

	x := gojieba.NewJieba()
	defer x.Free()

	useHmm := true
	words = x.CutForSearch(s, useHmm)

	var rword string
	for _, word := range words {
		if utf8.RuneCountInString(word) == 1 {
			continue
		} else if isStrSliceContainStr(excludes, word) {
			continue
		} else if isMapKeyMatchStr(mergeData, word) {
			rword = mergeData[word]
		} else {
			rword = word
		}
		if i, ok := counts[rword]; ok == false {
			counts[rword] = 1
		} else {
			counts[rword] = i + 1
		}
	}

	// 把字典转换为以词语出现次数DESC排序的列表
	for k, v := range counts {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// showCount := 30
	// for i := 0; i < showCount; i++ {
	// 	fmt.Printf("%s\t %d\n", string(ss[i].Key), ss[i].Value)
	// }

	return ss

	// fmt.Println(res.Meta)
}
