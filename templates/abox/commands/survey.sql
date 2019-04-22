-- 调研问卷

-- cmd: saveAnswer
-- 插入一条记录
INSERT INTO survey (
  surveyNum,
  submitUser,
  questionNum,
  answers,
  submitDate
) VALUES (
           :surveyNum,
           :submitUser,
           :questionNum,
           :answers,
           :submitDate
         );

-- cmd: deleteAnswer
-- 按照 问卷编号+提交人 删除
DELETE FROM survey
WHERE 1 = 1
  AND surveyNum = :surveyNum
  AND submitUser = :submitUser;


