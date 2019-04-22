-- doc word
-- cmd: wordCloud
select word, count, '' doc_name  from (
	select word, sum(count) count, string_agg(doc_name, ', ') doc_name from words group by word
) a
where a.count >= 5
order by a.count desc;

-- cmd: wordFreq
select a.doc_count "level",  count(1) "count" from (
	select word, count(distinct doc_name) doc_count from words group by word
) a
group by a.doc_count
order by a.doc_count;


-- cmd: wordDist
select a.word, a.wc_count from (
								   select word, sum(count) wc_count, string_agg(doc_name, ', ') doc_name, count(distinct doc_name) doc_count from words group by word
							   ) a
order by a.wc_count desc;