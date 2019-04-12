package base

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/gobuffalo/envy"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gofrs/uuid"
	"gopkg.in/gomail.v2"
)

// EmailServerT 邮件服务器信息
type EmailServerT struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	UserName    string `json:"userName"`
	Password    string `json:"password"`
	DisplayName string `jsonL:"displayName"`
}

// 发邮件：可以定时作业，也可以手工触发
// 无论哪种模式，都是从 email 表中，读取所有 未发邮件，然后批量发送

// getEmailServer 每次发邮件之前，都重新从文件获取一次服务器信息
func getEmailServer() (EmailServerT, error) {
	a := EmailServerT{}

	configFile, err := os.Open("./config/email.json")
	defer configFile.Close()
	if err != nil {
		log.Fatalf("config failed: %+v", err)
		return a, err
	}
	jsonParser := json.NewDecoder(configFile)

	err = jsonParser.Decode(&a)

	// 文件中存的密文
	// var data []byte
	// data, _ = base64.StdEncoding.DecodeString(a.Password)
	// a.Password = string(data)

	return a, err
}

// EmailNotice 每一个邮件通知
type EmailNotice struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	SendTo    string    `json:"send_to" db:"send_to"`
	UserName  string    `json:"user_name" db:"user_name"`
	Subject   string    `json:"subject" db:"subject"`
	Content   string    `json:"content" db:"content"`
	Status    string    `json:"status" db:"status"`
	SentDate  string    `json:"sent_date" db:"sent_date"`
}

const (
	// WAITING 邮件未发送
	WAITING = "WAITING"

	// SENT 邮件已发送
	SENT = "SENT"
)

// Insert 当前对象，插入到数据库中
func (t EmailNotice) Insert() error {
	var err error
	db := DB()
	kv := KvCache()
	// sugar := Sugar()

	// 创建时间为当前时刻
	now := time.Now()
	t.ID, _ = uuid.NewV4()
	t.UpdatedAt = now
	t.CreatedAt = now
	t.Status = WAITING
	t.SentDate = ""

	// 根据 struct 中的 DB tag 进行自动 named parameter
	insertCmd := kv.GetCommand("email.insert", nil)
	_, err = db.NamedExec(insertCmd, t)
	return err
}

// EmailLoadByID 加载一条记录
func EmailLoadByID(id uuid.UUID) EmailNotice {
	var err error

	db := DB()
	kv := KvCache()

	sqlCmd := kv.GetCommand("email.loadByID", nil)
	e := EmailNotice{}
	err = db.Get(&e, sqlCmd, id)
	if err != nil {
		sugar.Warnw("load email failed", "id", id, "err", err)
	}
	return e
}

// markSent 标记已发送
func (t EmailNotice) markSent() error {
	kv := KvCache()
	db := DB()

	// 设置已发送
	now := time.Now()
	t.UpdatedAt = now
	t.Status = SENT
	t.SentDate = now.Format("2006-01-02 15:04:05")

	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := kv.GetCommand("email.markSent", nil)
	_, err := db.NamedExec(sqlCmd, t)

	return err
}

// Send 发送通知，并且标记已发送
func (t EmailNotice) Send(server EmailServerT) error {
	surveyDir := envy.Get("SurveyDir", "./config/surveys")

	// 发送邮件
	m := gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", server.UserName, server.DisplayName)

	// 收件人
	m.SetHeader("To", m.FormatAddress(t.SendTo, t.UserName))

	// 抄送
	m.SetHeader("Cc", m.FormatAddress("ibmpartment@powerdekor.com.cn", "IBM项目组"))

	// 主题
	m.SetHeader("Subject", t.Subject)

	// 正文
	m.SetBody("text/html", t.Content)

	// 附件
	// m.Attach("./survey/ABCD.jpeg")
	m.Embed(surveyDir + "/ABCD.png")

	// 发送邮件服务器、端口、发件人账号、发件人密码
	d := gomail.NewDialer(server.Host, server.Port, server.UserName, server.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	// 标记已发送
	err = t.markSent()
	return err
}

// SendAllEmails 发送所有邮件
func SendAllEmails() error {
	var err error

	// 获取所需的环境，同时设置程序运行目录是项目的根目录
	db := DB()
	kv := KvCache()
	sugar := Sugar()

	// 每次发邮件前，都重新获取一次邮件服务器信息
	server, err := getEmailServer()
	if err != nil {
		sugar.Errorw("can not get email server", "err", err.Error())
		return err
	}
	sugar.Debugw("mail server", "server", server)

	// 取得所有未发送的邮件
	var items []EmailNotice
	cmdSQL := kv.GetCommand("email.getByStatus", nil)
	err = db.Select(&items, cmdSQL, WAITING)
	if err != nil {
		sugar.Errorw("find emails", "err", err, "sql", cmdSQL)
		return err
	}
	sugar.Debugw("sending mails", "total", len(items))

	// 发送邮件
	for _, i := range items {
		err = i.Send(server)
		if err != nil {
			// 发送失败后，打印错误信息，然后继续发下一个，而不是终止
			sugar.Errorw("send mail failed", "id", i.ID, "err", err.Error())
			return nil
		}
		time.Sleep(3000)
	}

	return nil
}

// GetEmailNotices 获取邮件通知
func GetEmailNotices() ([]EmailNotice, error) {
	kv := KvCache()
	db := DB()

	cmdSQL := kv.GetCommand("email.getNotices", nil)
	var items []EmailNotice
	err := db.Select(&items, cmdSQL)
	return items, err
}

// GenerateInvitation 根据 excel ，生成待发送邮件
// 生成的邮件统一会保存在 emailNotice 表中，然后统一发送
func GenerateInvitation() error {
	sugar := Sugar()

	surveyDir := envy.Get("SurveyDir", "./config/surveys")

	type userT struct {
		Title     string
		Email     string
		SurveyNum string
		ID        string
	}

	file, err := os.Open(surveyDir + "/invitation.xlsx")
	if err != nil {
		return err
	}

	// 从文件中读取
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return err
	}

	//  上载的文件必须有这个 tab
	sheetName := "Sheet1"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 500

	results := []userT{}
	for r := 2; r < MaxLoop; r++ {
		// 从第二行开始
		rowData := RowReader(readHelper, r)

		// 第一列：, 第二列：
		if len(rowData(0)) == 0 && len(rowData(1)) == 0 {
			break
		}

		entry := userT{
			Title:     rowData(0),
			Email:     rowData(1),
			SurveyNum: rowData(2),
			ID:        rowData(3),
		}

		results = append(results, entry)
	}
	sugar.Infow("users", "users", results)

	invFile := surveyDir + "/invitation.tmpl.html"
	sugar.Infow("template file", "file", invFile)
	t := template.Must(template.New("").ParseFiles(invFile))
	for _, i := range results {
		var buf bytes.Buffer
		err := t.ExecuteTemplate(&buf, "invitation", H{
			"Name":      i.Title,
			"SurveyNum": i.SurveyNum,
			"ID":        i.ID,
		})
		if err != nil {
			sugar.Fatal(err)
			return err
		}
		sugar.Debugw(buf.String())

		// 	生成待发送邮件
		notice := EmailNotice{
			UserName: i.Title,
			SendTo:   i.Email,
			Subject:  "企业认知离散度调研",
			Content:  buf.String(),
		}
		err = notice.Insert()
		if err != nil {
			sugar.Errorw("insert notice", "notice", notice, "err", err)
			continue
		}
	}

	return nil
}
