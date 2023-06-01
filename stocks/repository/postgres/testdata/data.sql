insert into ticker
(id, name)
values
    ('aad17418-6764-4ecd-90ed-bb1d7091edcc', 'AAPL'),
    ('3439d561-b4db-4455-aff9-da2119573574', 'AAPL2');

insert into split
(date, ticker_id, before, after)
values
    ('20-Jan-2000', 'aad17418-6764-4ecd-90ed-bb1d7091edcc', 1, 2),
    ('20-Jan-2005', 'aad17418-6764-4ecd-90ed-bb1d7091edcc', 1, 2);

insert into portfolio_record
(id, investor_id, ticker_id, date, amount, price_usd, action)
values
    ('23c484bf-9332-4b81-a482-aba787f78f32', '4ffdaa1c-9c25-4a79-a3a6-cf47ba361728', 'aad17418-6764-4ecd-90ed-bb1d7091edcc', '20-Jan-1999', 1, 1, 'BUY'),
    ('9454d154-cc05-42f7-b0ee-8b7b78390ab3', '4ffdaa1c-9c25-4a79-a3a6-cf47ba361728', 'aad17418-6764-4ecd-90ed-bb1d7091edcc', '20-Jan-2001', 1, 1, 'BUY'),
    ('c028e70f-55ba-4b55-868f-f22940346a8a', '4ffdaa1c-9c25-4a79-a3a6-cf47ba361728', 'aad17418-6764-4ecd-90ed-bb1d7091edcc', '20-Jan-2006', 1, 1, 'BUY'),
    ('792b02de-04c8-41b8-9ec2-c88ab152a944', '4ffdaa1c-9c25-4a79-a3a6-cf47ba361728', '3439d561-b4db-4455-aff9-da2119573574', '20-Jan-2001', 1, 1, 'BUY'),
    ('1d4370ea-4f53-44f7-a4b7-be0793cb6b2f', '4ffdaa1c-9c25-4a79-a3a6-cf47ba361728', 'aad17418-6764-4ecd-90ed-bb1d7091edcc', '20-Jan-2007', 1, 1, 'SELL');

insert into stock_daily
(id, ticker_id, date, high, low, open, close)
values
    ('f14dc8c3-f320-402d-8353-d49a3ba35f70','aad17418-6764-4ecd-90ed-bb1d7091edcc', '19-Jan-2010', 3, 2, 4, 5),
    ('5460ab36-feeb-45de-ab2a-dca06e904971', '3439d561-b4db-4455-aff9-da2119573574', '19-Jan-2010', 3, 2, 4, 5),
    ('fa160058-3012-46dc-8837-0cd264aa26da','aad17418-6764-4ecd-90ed-bb1d7091edcc', '20-Jan-2010', 3, 2, 4, 5),
    ('a57cb59e-d6ec-4f85-9a8a-883b8b960018', '3439d561-b4db-4455-aff9-da2119573574', '20-Jan-2010', 3, 2, 4, 5);
