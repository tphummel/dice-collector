delete from session; delete from turn; delete from throw;

select w.turnid as turn, w.sequence as seq, w.value as val, w.event_prop_comeout as co, 
r.name
from throw w, result r
where r.id = w.result
order by w.turnid, w.sequence;
