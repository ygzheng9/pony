-- 发送邮件

-- cmd: loadByID
SELECT
  id,
  send_to,
  user_name,
  subject,
  content,
  status,
  sent_date,
  created_at,
  updated_at
FROM notices
where id = $1;



-- cmd: insert
INSERT INTO notices (
  id,
  send_to,
  user_name,
  subject,
  content,
  status,
  sent_date,
  created_at,
  updated_at
) VALUES (
  :id,
  :send_to,
  :user_name,
  :subject,
  :content,
  :status,
  :sent_date,
  :created_at,
  :updated_at
);

-- cmd: markSent
update notices set
        status = :status,
        sent_date = :sent_date,
        updated_at = :updated_at
where id=:id;


-- cmd: getByStatus
-- 取得未发送的邮件
SELECT
  id,
  send_to,
  user_name,
  subject,
  content,
  status,
  sent_date,
  created_at,
  updated_at
FROM notices
where status = $1
order by created_at desc;


-- cmd: getNotices
-- 取得近 90 天内的邮件
-- 这个还没有修改
select ID, ifnull(CRE_DTE,'') CRE_DTE, ifnull(SEND_TO,'') SEND_TO, ifnull(USR_NME,'') USR_NME,
       ifnull(SUBJ,'') SUBJ, ifnull(CTENT,'') CTENT, ifnull(STS,'') STS, ifnull(SENT_DTE,'') SENT_DTE
from T_EMAIL
where 1 = 1
  and CRE_DTE >= date_format(date_add(now(), interval -90 day), '%Y-%m-%d')
order by CRE_DTE desc;


