-- 根据 tablename，从数据库中找到 columns
-- cmd: mysqlColumns
SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_KEY
  FROM information_schema.columns
 WHERE table_schema=? AND table_name=?;

-- cmd: postgresColumns
SELECT column_name, data_type, is_nullable
FROM information_schema.columns
WHERE 1 = 1
  AND table_catalog = $1
  AND table_name   = $2 ;

