create table partner_payment_customer_contracts
(
    "id"                         serial primary key,
    "partner_payment_history_id" bigint references partner_payment_histories (id),
    "customer_contract_id"       bigint references customer_contracts (id),
    "created_at"                 timestamptz           DEFAULT (now()),
    "updated_at"                 timestamptz           DEFAULT (now())
);
