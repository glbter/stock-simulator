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

insert into stock_daily
(id, ticker_id, date, high, low, open, close)
values
    ('60c9b9f6-453d-4035-86da-3f16a01e6361','aad17418-6764-4ecd-90ed-bb1d7091edcc', '20-Jan-2009', 3, 2, 4, 5),
    ('fa160058-3012-46dc-8837-0cd264aa26da','aad17418-6764-4ecd-90ed-bb1d7091edcc', '20-Jan-2010', 3, 2, 4, 5),
    ('a57cb59e-d6ec-4f85-9a8a-883b8b960018', '3439d561-b4db-4455-aff9-da2119573574', '20-Jan-2010', 4, 3, 5, 6),
    ('39c0ec64-c07b-48ff-afc0-c3a7f6e83bd3', '3439d561-b4db-4455-aff9-da2119573574', '21-Jan-2010', 5, 4, 6, 7),
    ('0ed01c72-cd61-4a31-9c1b-ebd5fa829820', '3439d561-b4db-4455-aff9-da2119573574', '22-Jan-2010', 6, 5, 7, 8),
    ('d6469956-4169-4ddd-a534-4d0a2cb15e9c', '3439d561-b4db-4455-aff9-da2119573574', '23-Jan-2010', 7, 6, 8, 9),
    ('8df9d2cd-3686-472e-83fd-a1826e97fce9', '3439d561-b4db-4455-aff9-da2119573574', '24-Jan-2010', 8, 7, 9, 10),
    ('349dc9ed-51b3-496d-aaf8-42ee508656df', '3439d561-b4db-4455-aff9-da2119573574', '25-Jan-2010', 9, 8, 10, 11);

