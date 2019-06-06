package actions

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"

	"pony/base"
	"pony/models"
)

// gamesIndex 决策模型
func gamesIndex(c buffalo.Context) error {
	return c.Render(200, r.HTML("games/tailwind_index.html", "games/tailwind_layout.html"))
}

// gamesByID find by id
func gamesByID(c buffalo.Context) error {
	var err error
	sugar := base.Sugar()

	type paramT struct {
		ID string `db:"id" json:"id" form:"id"`
	}
	param := paramT{}
	err = c.Bind(&param)
	if err != nil {
		sugar.Errorw("bind", "err", err)
		return errors.WithStack(err)
	}
	sugar.Debugw("param", "param", param)

	id, err := uuid.FromString(param.ID)
	if err != nil {
		sugar.Errorw("uuid failed", "param", param, "err", err)
		return errors.WithStack(err)
	}

	// 根据 matrix，period，company 找到之前提交的数据
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	entry := models.Game{}
	err = tx.Find(&entry, id)
	if err != nil {
		sugar.Errorw("no entry", "param", param, "err", "err")
		return errors.WithStack(err)
	}

	result := pairCompareT{}
	err = json.Unmarshal([]byte(entry.Weights), &result)
	if err != nil {
		sugar.Errorw("unmarshal", "entry", entry)
	}

	sugar.Debugw("got entry", "entry", entry)
	return c.Render(200, r.JSON(H{
		"entry":  entry,
		"result": result,
	}))
}

// gamesCreate creat new game
func gamesCreate(c buffalo.Context) error {
	var err error
	sugar := base.Sugar()

	type paramT struct {
		Name string `db:"name" json:"name"`
	}
	param := paramT{}
	err = c.Bind(&param)
	if err != nil {
		return errors.WithStack(err)
	}
	sugar.Infow("param", "param", param)

	entry := models.Game{
		Name: param.Name,
	}

	// 根据 matrix，period，company 找到之前提交的数据
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// find existing, if no, insert.
	games := models.Games{}
	query := tx.Where("name = ?", entry.Name)
	err = query.All(&games)
	if err != nil {
		sugar.Errorw("sql error", "err", err)
		return errors.WithStack(err)
	}
	if len(games) >= 1 {
		// find existing
		return c.Render(200, r.JSON(games[0]))
	}

	// insert
	_, err = tx.ValidateAndSave(&entry)
	if err != nil {
		sugar.Errorw("create game", "entry", entry, "err", err)
		return errors.WithStack(err)
	}
	return c.Render(200, r.JSON(entry))
}

// gamesSaveCriterion save criterion in raw format
func gamesSaveCriterion(c buffalo.Context) error {
	var err error
	sugar := base.Sugar()

	type paramT struct {
		ID        string `db:"id" json:"id"`
		Criterion string `db:"criterion" json:"criterion"`
	}
	param := paramT{}
	err = c.Bind(&param)
	if err != nil {
		return errors.WithStack(err)
	}
	sugar.Infow("param", "param", param)

	id, err := uuid.FromString(param.ID)
	if err != nil {
		sugar.Errorw("uuid failed", "param", param, "err", err)
		return errors.WithStack(err)
	}

	// 根据 matrix，period，company 找到之前提交的数据
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	entry := models.Game{}
	err = tx.Find(&entry, id)
	if err != nil {
		sugar.Errorw("no entry", "param", param, "err", "err")
		return errors.WithStack(err)
	}

	entry.Criterion = param.Criterion
	_, err = tx.ValidateAndSave(&entry)
	if err != nil {
		sugar.Errorw("save criterion", "entry", entry, "err", err)
		return errors.WithStack(err)
	}

	items := strings.Split(entry.Criterion, "\n")
	var results []string
	for _, i := range items {
		if i != "" {
			results = append(results, i)
		}
	}
	sugar.Debugw("criterion", "raw", entry.Criterion, "items", items, "results", results)

	return c.Render(200, r.JSON(H{
		"items": results,
		"raw":   entry.Criterion,
	}))
}

// gamesSaveCriterion save options in raw format
func gamesSaveOptions(c buffalo.Context) error {
	var err error
	sugar := base.Sugar()

	type paramT struct {
		ID      string `db:"id" json:"id"`
		Options string `db:"options" json:"options"`
	}
	param := paramT{}
	err = c.Bind(&param)
	if err != nil {
		return errors.WithStack(err)
	}
	sugar.Infow("param", "param", param)

	id, err := uuid.FromString(param.ID)
	if err != nil {
		sugar.Errorw("uuid failed", "param", param, "err", err)
		return errors.WithStack(err)
	}

	// 根据 matrix，period，company 找到之前提交的数据
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	entry := models.Game{}
	err = tx.Find(&entry, id)
	if err != nil {
		sugar.Errorw("no entry", "param", param, "err", "err")
		return errors.WithStack(err)
	}

	entry.Options = param.Options
	_, err = tx.ValidateAndSave(&entry)
	if err != nil {
		sugar.Errorw("save options", "entry", entry, "err", err)
		return errors.WithStack(err)
	}

	items := strings.Split(entry.Options, "\n")
	var results []string
	for _, i := range items {
		if i != "" {
			results = append(results, i)
		}
	}
	sugar.Debugw("criterion", "raw", entry.Options, "items", items, "results", results)

	return c.Render(200, r.JSON(H{
		"items": results,
		"raw":   entry.Options,
	}))
}

// gamesSaveCriterionPairs save pairs for  criterion in raw format
func gamesSaveCriterionPairs(c buffalo.Context) error {
	var err error
	sugar := base.Sugar()

	type paramT struct {
		ID    string `db:"id" json:"id"`
		Pairs string `db:"pairs" json:"pairs"`
	}
	param := paramT{}
	err = c.Bind(&param)
	if err != nil {
		return errors.WithStack(err)
	}
	sugar.Infow("param", "param", param)

	id, err := uuid.FromString(param.ID)
	if err != nil {
		sugar.Errorw("uuid failed", "param", param, "err", err)
		return errors.WithStack(err)
	}

	// 根据 matrix，period，company 找到之前提交的数据
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	entry := models.Game{}
	err = tx.Find(&entry, id)
	if err != nil {
		sugar.Errorw("no entry", "param", param, "err", "err")
		return errors.WithStack(err)
	}

	result, s, err := gamesCalcMaxEig(param.Pairs)
	if err != nil {
		sugar.Errorw("calc eig", "err", err)
		return errors.WithStack(err)
	}
	entry.Weights = s
	entry.Pairs = param.Pairs
	_, err = tx.ValidateAndSave(&entry)
	if err != nil {
		sugar.Errorw("save pairs", "entry", entry, "err", err)
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(H{
		"item":   entry,
		"result": result,
	}))
}

// pairCompareT 两两比较的结果
type pairCompareT struct {
	LambdaMax float64   `json:"lambdaMax"`
	CR        float64   `json:"cr"`
	Weights   []float64 `json:"weights"`
}

// gamesCalcMaxEig calc max eig value and vector
func gamesCalcMaxEig(s string) (pairCompareT, string, error) {
	var err error
	sugar := base.Sugar()

	var raw [][]string
	// s, err := strconv.Unquote(string(param.Pairs))
	// if err != nil {
	// 	sugar.Errorw("unquote",  "err", err)
	// 	return errors.WithStack(err)
	// }

	err = json.Unmarshal([]byte(s), &raw)
	if err != nil {
		sugar.Errorw("unmarshal", "s", s, "err", err)
		return pairCompareT{}, "", err
	}
	sugar.Debugw("pairs", "input", s, "pairs", raw)

	count := len(raw)
	data := make([][]float64, count)
	for i := range raw {
		data[i] = make([]float64, count)
	}
	for idxi, i := range raw {
		for idxj, j := range i {
			if j == "" {
				sugar.Warnw("parse", "zero", j)
				data[idxi][idxj] = 0.0
			} else if len(j) == 1 {
				d, err := strconv.ParseFloat(j, 64)
				if err != nil {
					sugar.Debugw("parse1", "raw", j, "err", err)
					data[idxi][idxj] = 0.0
				} else {
					data[idxi][idxj] = d
				}
			} else if len(j) == 3 {
				s := strings.Split(j, "/")
				a, err := strconv.ParseFloat(s[0], 64)
				if err != nil {
					sugar.Debugw("parse2a", "raw", j, "err", err)
					data[idxi][idxj] = 0
					continue
				}
				b, err := strconv.ParseFloat(s[1], 64)
				if err != nil {
					sugar.Debugw("parse2b", "raw", j, "err", err)
					data[idxi][idxj] = 0
					continue
				}

				data[idxi][idxj] = a / b
			} else {
				sugar.Warnw("unknown", "raw", j)
				data[idxi][idxj] = 0.0
			}
		}
	}
	sugar.Debugw("data", "[][]float", data)

	var input []float64
	for _, i := range data {
		input = append(input, i...)
	}
	sugar.Debugw("input", "[]float", input)

	a := mat.NewDense(count, count, input)
	sugar.Debugf("A = \n%v\n\n", mat.Formatted(a, mat.Prefix("")))

	var eig mat.Eigen
	ok := eig.Factorize(a, mat.EigenRight)
	if !ok {
		sugar.Fatal("Eigen decomposition failed")
	}
	sugar.Debugf("Eigenvalues of A:\n%v\n", eig.Values(nil))

	values := eig.Values(nil)
	v := real(values[0])
	sugar.Debugf("max eigValue: %.2f ", v)

	sugar.Debugf("Eigenvectors of A:\n%v\n", eig.VectorsTo(nil))
	vectors := eig.VectorsTo(nil)
	vect := []float64{}
	for i := 0; i < count; i++ {
		vect = append(vect, real(vectors.At(i, 0)))
	}
	sugar.Debugf("max vectors :\n%v\n", vect)

	ri := []float64{0, 0, 0, 0.58, 0.96, 1.12, 1.24, 1.32, 1.41, 1.45}
	dm := 2.0
	if count < len(ri) {
		dm = ri[count]
	}
	cr := (v - float64(count)) / dm

	result := pairCompareT{
		LambdaMax: v,
		CR:        cr,
		Weights:   vect,
	}

	buf, err := json.Marshal(result)
	if err != nil {
		sugar.Errorw("marshal", "result", result)
	}

	return result, string(buf), nil
}

// gamesSaveOptionPairs save pairs for  option by each criteria in raw format
func gamesSaveOptionPairs(c buffalo.Context) error {
	var err error
	sugar := base.Sugar()

	type paramT struct {
		ID    string `json:"id"`
		Seq   string `json:"seq"`
		Pairs string `json:"pairs"`
	}
	param := paramT{}
	err = c.Bind(&param)
	if err != nil {
		return errors.WithStack(err)
	}
	sugar.Infow("param", "param", param)

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	var entries models.GameOptions
	cond := fmt.Sprintf("game_id = '%s' and seq = %s", param.ID, param.Seq)
	err = tx.Where(cond).All(&entries)
	if err != nil {
		sugar.Errorw("sql error", "param", param, "err", "err")
		return errors.WithStack(err)
	}

	var entry models.GameOption
	if len(entries) == 0 {
		entry = models.GameOption{}

		id, err := uuid.FromString(param.ID)
		if err != nil {
			sugar.Errorw("uuid failed", "param", param, "err", err)
			return errors.WithStack(err)
		}
		entry.GameID = id
		i, err := strconv.Atoi(param.Seq)
		if err != nil {
			sugar.Errorw("parse seq", "param", param)
			return errors.WithStack(err)
		}
		entry.Seq = i
		sugar.Debugw("new", "entry", entry)
	} else {
		entry = entries[0]
		sugar.Debugw("exist", "entry", entry)
	}

	entry.Pairs = param.Pairs
	sugar.Debugw("save", "entry", entry)

	result, s, err := gamesCalcMaxEig(param.Pairs)
	if err != nil {
		sugar.Errorw("calc eig", "err", err)
		return errors.WithStack(err)
	}
	entry.Weights = s

	verrs, err := tx.ValidateAndSave(&entry)
	if err != nil {
		sugar.Errorw("save pairs", "entry", entry, "err", err)
		return errors.WithStack(err)
	}
	if verrs.Count() > 0 {
		sugar.Errorw("validate", "verrs", verrs)
		return errors.WithStack(verrs)
	}

	return c.Render(200, r.JSON(H{
		"item":   entry,
		"result": result,
	}))
}

// gamesLoadOptionPairs load pairs for option by each criteria in raw format
func gamesLoadOptionPairs(c buffalo.Context) error {
	var err error
	sugar := base.Sugar()

	type paramT struct {
		ID  string `json:"id"`
		Seq string `json:"seq"`
	}
	param := paramT{}
	err = c.Bind(&param)
	if err != nil {
		return errors.WithStack(err)
	}
	sugar.Infow("param", "param", param)

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	var entries models.GameOptions
	cond := fmt.Sprintf("game_id = '%s' and seq = %s", param.ID, param.Seq)
	sugar.Debugw("option cond", "cond", cond)

	err = tx.Where(cond).All(&entries)
	if err != nil {
		sugar.Errorw("sql error", "param", param, "err", err)
		return errors.WithStack(err)
	}

	result := pairCompareT{}
	var entry models.GameOption
	if len(entries) == 0 {
		entry = models.GameOption{}
	} else {
		entry = entries[0]

		err = json.Unmarshal([]byte(entry.Weights), &result)
		if err != nil {
			sugar.Errorw("unmarshal", "entry", entry)
		}
	}

	return c.Render(200, r.JSON(H{
		"item":   entry,
		"result": result,
	}))
}

// gamesCalcFinal calc option sequence
func gamesCalcFinal(c buffalo.Context) error {
	var err error
	sugar := base.Sugar()

	type paramT struct {
		ID string `json:"id"`
	}
	param := paramT{}
	err = c.Bind(&param)
	if err != nil {
		sugar.Errorw("bind", "err", err)
		return c.Render(200, r.JSON(H{
			"status":  5,
			"message": "bind err",
		}))
	}
	sugar.Infow("param", "param", param)

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	id, err := uuid.FromString(param.ID)
	if err != nil {
		sugar.Errorw("uuid failed", "param", param, "err", err)
		return c.Render(200, r.JSON(H{
			"status":  6,
			"message": "uuid failed",
		}))
	}

	entry := models.Game{}
	err = tx.Find(&entry, id)
	if err != nil {
		sugar.Errorw("no entry", "param", param, "err", "err")
		return errors.WithStack(err)
	}

	// option count
	optItems := strings.Split(entry.Options, "\n")
	var optVals []string
	for _, i := range optItems {
		if i != "" {
			optVals = append(optVals, i)
		}
	}
	a := len(optVals)

	// criterion count
	criItems := strings.Split(entry.Criterion, "\n")
	var criVals []string
	for _, i := range criItems {
		if i != "" {
			criVals = append(criVals, i)
		}
	}
	b := len(criVals)

	result := pairCompareT{}
	err = json.Unmarshal([]byte(entry.Weights), &result)
	if err != nil {
		sugar.Errorw("unmarshal", "entry", entry)
		// return errors.WithStack(err)
		return c.Render(200, r.JSON(H{
			"status":  7,
			"message": "no weights",
		}))
	}

	if len(result.Weights) != b {
		sugar.Errorw("criterion pair mismatch", "entry", entry, "result", result)
		return c.Render(200, r.JSON(H{
			"status":  10,
			"message": "Should calc criterion",
		}))
	}

	options := models.GameOptions{}
	err = tx.Where("game_id = ?", entry.ID).Order("seq ").All(&options)
	if err != nil {
		sugar.Errorw("sql error", "err", err)
		return errors.WithStack(err)
	}

	// 每个条件都有比较矩阵
	if len(options) < b {
		sugar.Errorw("pair", "require", b, "actual", len(options))
		return c.Render(200, r.JSON(H{
			"status":  20,
			"message": "each criteria should have options pairs",
		}))
	}

	var m []float64
	for i := 0; i < b; i++ {
		for _, opt := range options {
			if opt.Seq == i {
				p := pairCompareT{}
				err = json.Unmarshal([]byte(opt.Weights), &p)
				if err != nil {
					sugar.Errorw("unmarshal", "opt")
					return errors.WithStack(err)
				}
				if len(p.Weights) != a {
					sugar.Errorw("option pair mismatch", "i", i, "pair", p)
					return c.Render(200, r.JSON(H{
						"status":  30,
						"message": "option pair criterion",
					}))
				}
				m = append(m, p.Weights...)
				break
			}
		}
	}

	total := a * b
	if len(m) != total {
		sugar.Errorw("wrong matrix", "a", a, "b", b, "actual", len(m), "m", m)
		return c.Render(200, r.JSON(H{
			"status":  40,
			"message": "matrix wrong",
		}))
	}

	sugar.Debugw("matrix", "a", a, "b", b, "m", m)
	d1 := mat.NewDense(b, a, m)
	d1t := d1.T()
	sugar.Debugf("M1 = \n%v\n\n", mat.Formatted(d1t, mat.Prefix(""), mat.Squeeze()))

	d2 := mat.NewVecDense(b, result.Weights)
	sugar.Debugf("M2 = \n%v\n\n", mat.Formatted(d2, mat.Prefix(""), mat.Squeeze()))

	var d3 mat.Dense
	d3.Mul(d1t, d2)
	sugar.Debugf("M3 = \n%v\n\n", mat.Formatted(&d3, mat.Prefix(""), mat.Squeeze()))

	score := mat.Col(nil, 0, &d3)

	return c.Render(200, r.JSON(H{
		"status":  0,
		"score":   score,
		"options": optItems,
	}))
}
