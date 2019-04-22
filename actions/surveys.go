package actions

import (
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"pony/base"
	"pony/models"
)

// SurveysOpen default implementation.
func SurveysOpen(c buffalo.Context) error {
	// 从 querystring 中获取 问卷编号
	var err error
	// 解析 querystring
	info := struct {
		SurveyNum  string `form:"num" json:"num" db:"surveyNum"`
		SurveyUser string `form:"id" json:"id" db:"submitUser"`
	}{}
	err = c.Bind(&info)
	if err != nil {
		return errors.WithStack(err)
	}

	s, err := loadSurvey(info.SurveyNum)
	if err != nil {
		return errors.WithStack(err)
	}

	p := struct {
		SurveyNum string
		ID        string
		Title     string
		Sections  []sectionT
	}{
		SurveyNum: info.SurveyNum,
		ID:        info.SurveyUser,
		Title:     s.Title,
		Sections:  s.Sections,
	}
	c.Set("p", p)

	return c.Render(200, r.HTML("surveys/open.html", "surveys/simple.html"))
}

type surveyT struct {
	Title    string     `yaml:"title"`
	Sections []sectionT `yaml:"sections"`
}

type sectionT struct {
	Section   string      `yaml:"section"`
	Questions []questionT `yaml:"questions"`
}

type questionT struct {
	Question string   `yaml:"question"`
	SeqNum   string   `yaml:"seq"`
	Type     string   `yaml:"type"`
	Options  []string `yaml:"options"`
}

func loadSurvey(surveyNum string) (surveyT, error) {
	var s surveyT

	// dir := envy.Get("SurveyDir", "")
	// fileName := dir + "/" + surveyNum + ".yaml"
	// source, err := ioutil.ReadFile(fileName)
	// if err != nil {
	// 	return s, err
	// }

	fileName := "surveys/" + surveyNum + ".yaml"
	source, err := base.ABox.Find(fileName)

	err = yaml.Unmarshal(source, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

// SurveysSubmit default implementation.
func SurveysSubmit(c buffalo.Context) error {
	// 从 querystring 中获取 问卷编号
	var err error
	// 解析 querystring
	info := struct {
		SurveyNum  string `form:"num" json:"num" db:"surveyNum"`
		SurveyUser string `form:"id" json:"id" db:"submitUser"`
	}{}
	err = c.Bind(&info)
	if err != nil {
		return errors.WithStack(err)
	}

	survey, err := loadSurvey(info.SurveyNum)
	if err != nil {
		return errors.WithStack(err)
	}

	// 获取 form data
	err = c.Request().ParseForm()
	if err != nil {
		return errors.WithStack(err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	var answers models.Surveys

	// 获取每一题的答案，都转换成字符串
	for _, s := range survey.Sections {
		for _, q := range s.Questions {
			a := models.Survey{
				SurveyNo:   info.SurveyNum,
				SubmitUser: info.SurveyUser,
				QuestionNo: q.SeqNum,
				Answers:    strings.Join(c.Request().Form[q.SeqNum], " "),
				SubmitDate: now,
			}
			answers = append(answers, a)
		}
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	if len(answers) > 0 {
		var err error

		sugar := base.Sugar()
		// db := tx.TX

		// 按照 (问卷编号 + 提交人) 删除
		// cmdSql = kv.GetCommand("survey.deleteAnswer", nil)
		// _, err = db.NamedExec(cmdSql, answers[0])
		// if err != nil {
		// 	sugar.Errorw("failed to delete answer", "answer", answers[0], "err", err)
		// 	return errors.WithStack(err)
		// }

		for i := range answers {
			var a = answers[i]
			_, err = tx.ValidateAndCreate(&a)
			if err != nil {
				sugar.Errorw("failed to save answer", "answer", a, "err", err)
				continue
			}
		}
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "survey.submit.success"))

	return c.Render(200, r.HTML("surveys/submit.html", "surveys/simple.html"))
}
