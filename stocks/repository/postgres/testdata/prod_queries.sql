-- collect portfolio state (current amount of all tickers) and current net worth
select
    t.ticker_id,
    t.sum as amount,
    tt.name,
    tt.description,
    t.sum*stock.high as high,
    t.sum*stock.low as low,
    t.sum*stock.open as open,
    t.sum*stock.close as close
from (
         select ticker_id, sum(total_amount)
         from (
                  select action,
                         ticker_id,
                         (case
                              when t.action = 'BUY' then t.total_amount
                              when t.action = 'SELL' then -t.total_amount
                             end) total_amount
                  from (
                           select sum(t.mul * pr.amount) as total_amount,
                                  t.action,
                                  t.ticker_id
                           from (
                                    select t.id,
                                           t.action,
                                           t.ticker_id,
                                           mul(t.split) as mul
                                    from (
                                             select pr.ticker_id,
                                                    pr.id,
                                                    coalesce(s.after / s.before, 1) as split,
                                                    pr.action
                                             from portfolio_record as pr
                                                      left join split as s
                                                                on s.ticker_id = pr.ticker_id and s.date > pr.date
                                             where investor_id = '4ffdaa1c-9c25-4a79-a3a6-cf47ba361728'

                                         ) as t
                                    group by t.id, t.action, t.ticker_id
                                ) as t
                                    join portfolio_record as pr on pr.id = t.id
                           group by t.ticker_id, t.action
                       ) as t
              ) as t
         group by t.ticker_id
     ) as t
         join ticker tt on t.ticker_id = tt.id
         join (
    select *
    from stock_daily
    where id in (
        select id
        from (
                 select id, max(date)
                 from stock_daily
                 group by id) t
    )
) stock
              on t.ticker_id = stock.ticker_id;





-- IMPROVED AND REFACTORED
-- collect portfolio state (current amount of all tickers) and current net worth
select
    t.ticker_id,
    t.total_amount as amount,
    tt.name,
    tt.description,
    t.total_amount*stock.high as high,
    t.total_amount*stock.low as low,
    t.total_amount*stock.open as open,
    t.total_amount*stock.close as close
from ticker tt
join (
    select ticker_id, high, low, open, close
    from stock_daily
    where id in (
        select id
        from (
             select id, max(date)
             from stock_daily
             group by id
        ) t
    )
) stock
    on tt.id = stock.ticker_id
join (
    select
        t.ticker_id,
        sum(t.mul * pr.amount) as total_amount
    from portfolio_record as pr
    join (
        select
            pr.id,
            pr.action,
            pr.ticker_id,
            (case
                 when pr.action = 'BUY' then mul(coalesce(s.after / s.before, 1))
                 when pr.action = 'SELL' then -mul(coalesce(s.after / s.before, 1))
                end) as mul
        from split as s
        right join portfolio_record as pr
            on s.ticker_id = pr.ticker_id
            and s.date > pr.date
        where pr.investor_id = '4ffdaa1c-9c25-4a79-a3a6-cf47ba361728'
        group by pr.id, pr.action, pr.ticker_id
    ) as t
      on pr.id = t.id
    group by t.ticker_id
) as t
    on t.ticker_id = tt.id;

