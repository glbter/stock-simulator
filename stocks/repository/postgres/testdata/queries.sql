
-- get portfolio state without split (not tested)

select
    t.name,
--     daily.high * amount,
--     daily.low * amount,
--     (select t.name from ticker as t where t.id = buy.ticker_id) as ticker_name,
    buy.ticker_id as ticker_id,
    (buy.amount - sell.amount) as amount
from (
     select ticker_id,
            sum(amount) as amount
     from portfolio_record
     where investor_id = '4ffdaa1c-9c25-4a79-a3a6-cf47ba361728' and action = 'BUY'
     group by ticker_id, investor_id, action
) as buy
join (
    select ticker_id,
           sum(amount) as amount
    from portfolio_record
    where investor_id = '4ffdaa1c-9c25-4a79-a3a6-cf47ba361728' and action = 'SELL'
    group by ticker_id, investor_id, action
) as sell
    on buy.ticker_id = sell.ticker_id
join ticker as t
    on t.id = buy.ticker_id;
-- join stock_daily as daily
--     on daily.ticker_id = buy.ticker_id



-- get portfolio state without split (not tested)
--
-- select
--     t.name,
--     daily.high * amount,
--     daily.low * amount,
-- --     (select t.name from ticker as t where t.id = buy.ticker_id) as ticker_name,
--     buy.ticker_id as ticker_id,
--     (buy.amount - sell.amount) as amount,
-- from (
--          select ticker_id,
--                 sum(amount) as amount
--          from portfolio_record
--          where investor_id = $1 and action = "BUY"
--          group by ticker_id, investor_id, action
--      ) as buy
--          join (
--     select ticker_id,
--            sum(amount) as amount
--     from portfolio_record
--     where investor_id = $1 and action = "SELL"
--     group by ticker_id, investor_id, action
-- ) as sell
--               on buy.ticker_id = on sell.ticker_id
--     join ticker as t
--     on t.id = buy.ticker_id
--     join stock_daily as daily
--     on daily.ticker_id = buy.ticker_id

--
-- select ticker_id,
--        sum(amount) as amount
-- from portfolio_record
-- where investor_id = $1 and action = "SELL"
-- group by ticker_id, investor_id, action

-- CREATE AGGREGATE mul(bigint) ( SFUNC = int8mul, STYPE=float8 );

CREATE FUNCTION mul_accum (numeric, numeric, numeric)
    RETURNS numeric AS
$$
SELECT $1*$2*$3;
$$ LANGUAGE 'sql' STRICT;

CREATE FUNCTION mul_final(numeric)
    RETURNS numeric AS
$$
SELECT $1;
$$ LANGUAGE 'sql' STRICT;

CREATE AGGREGATE mul(bigint) ( INITCOND = 1, SFUNC = int8mul, STYPE=bigint );

CREATE AGGREGATE mul(numeric, numeric) (
    INITCOND = 1,
        STYPE = numeric,
        SFUNC = mul_accum
--         FINALFUNC = public.mul_final
);

CREATE AGGREGATE mul(numeric) (
    INITCOND = 1,
    STYPE = numeric,
    SFUNC = mul_accum
--         FINALFUNC = public.mul_final
    );

CREATE AGGREGATE mul(numeric) (
    INITCOND = 1,
    STYPE = numeric,
    SFUNC = numeric_mul
    --     SFUNC = mul_accum
--         FINALFUNC = public.mul_final
    );


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
    select ticker_id, date, high, low, open, close
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
    select  ticker_id, high, low, open, close
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
          -- count for each record the multiply value
            select
               pr.id as check_id,
--                pr.action,
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
      on pr.id = t.check_id
      group by t.ticker_id
) as t
on t.ticker_id = tt.id;





-- collect portfolio state (current amount of all tickers) and current net worth
with portfolio_state as (
    select
           t.ticker_id,
           t.total_amount               as amount,
           tt.name,
           tt.description,
           t.total_amount * stock.high  as high,
           t.total_amount * stock.low   as low,
           t.total_amount * stock.open  as open,
           t.total_amount * stock.close as close
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
            -- count for each record the multiply value
            select
                   pr.id as check_id,
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
            on pr.id = t.check_id
        group by t.ticker_id
    ) as t
        on t.ticker_id = tt.id
)

select
    ticker_id,
    amount,
    name,
    description,
    high,
    low,
    open,
    close
from (
    select
            ticker_id,
            amount,
            name,
            description,
            high,
            low,
            open,
            close,
            1 as ordering
    from portfolio_state
    union
    select
            uuid_nil()  as ticker_id,
            sum(amount) as amount,
            'TOTAL'     as name,
            'TOTAL'     as description,
            sum(high)   as high,
            sum(low)    as low,
            sum(open)   as open,
            sum(close)  as close,
            2           as ordering
    from portfolio_state
    order by ordering, name
) t;





with portfolio_state as (
    select
        t.ticker_id,
        t.total_amount               as amount,
        tt.name,
        tt.description,
        t.total_amount * stock.high  as high,
        t.total_amount * stock.low   as low,
        t.total_amount * stock.open  as open,
        t.total_amount * stock.close as close
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
            -- count for each record the multiply value
            select
                pr.id as check_id,
                pr.ticker_id,
                (case
                     when pr.action = 'BUY' then mul(coalesce(s.after / s.before, 1))
                     when pr.action = 'SELL' then -mul(coalesce(s.after / s.before, 1))
                    end) as mul
            from split as s
                     right join portfolio_record as pr
                                on s.ticker_id = pr.ticker_id
                                    and s.date > pr.date
            where pr.investor_id = $1
            group by pr.id, pr.action, pr.ticker_id
        ) as t
                      on pr.id = t.check_id
        group by t.ticker_id
    ) as t
                  on t.ticker_id = tt.id
)

select
    ticker_id,
    coalesce(amount, 0) as amount,
    name,
    coalesce(description, '') as description,
    coalesce(high, 0) as high,
    coalesce(low, 0) as low,
    coalesce(open, 0) as open,
    coalesce(close, 0) as close
from (
         select
             ticker_id,
             amount,
             name,
             description,
             high,
             low,
             open,
             close,
             1 as ordering
         from portfolio_state
         union all
         select
             uuid_nil()  as ticker_id,
             sum(amount) as amount,
             'TOTAL'     as name,
             'TOTAL'     as description,
             sum(high)   as high,
             sum(low)    as low,
             sum(open)   as open,
             sum(close)  as close,
             2           as ordering
         from portfolio_state
         order by ordering, name
     ) t;
















select
    ticker.id,
    ticker.name,
    coalesce(ticker.description, ''),
    sd.date,
    coalesce(sd.high),
    coalesce(sd.low),
    coalesce(sd.open),
    coalesce(sd.close),
    coalesce(sd.volume)
from ticker
join (
    select ticker_id, max(date) as date
    from stock_daily
    group by ticker_id
) latest_daily
    on ticker.id = latest_daily.ticker_id
join stock_daily sd
    on ticker.id = sd.ticker_id and latest_daily.date = sd.date


select * from stock_daily
where ticker_id
where date between
1


