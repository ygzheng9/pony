-- doc word
-- cmd: wordCloud
select word, count, '' doc_name  from (
	select word, sum(count) count, string_agg(doc_name, ', ') doc_name from words group by word
) a
where a.count >= 5
order by a.count desc;
