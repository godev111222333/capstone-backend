alter table partner_payment_histories
    add column "payment_url" varchar(1023) not null default '';
