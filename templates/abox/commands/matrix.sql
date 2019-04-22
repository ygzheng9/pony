-- cmd: findBy
select id, company, period, matrix, code, value,
       submit_user, created_at, updated_at
  from matrices
 where company = :company
   and version = :version
   and period = :period
   and matrix = :matrix;

-- cmd: checkExist
select id
  from matrices
where company = :company
  and version = :version
  and period = :period
  and matrix = :matrix
  and code = :code;

-- cmd: updateValue
update matrices
  set value = :value,
      updated_at = :updated_at
where id = :id;

