create table roles
(
    "id"         serial primary key,
    "role_name"  varchar(255) not null default '',
    "role_code"  varchar(255) not null default '',
    "created_at" timestamptz           DEFAULT (now()),
    "updated_at" timestamptz           DEFAULT (now())
);

create table accounts
(
    "id"                         serial primary key,
    "role_id"                    bigint references roles (id),
    "first_name"                 varchar(255)  not null default '',
    "last_name"                  varchar(255)  not null default '',
    "phone_number"               varchar(255)  not null unique,
    "email"                      varchar(255)  not null default '',
    "identification_card_number" varchar(255)  not null default '',
    "password"                   varchar(255)  not null default '',
    "avatar_url"                 varchar(255)  not null default '',
    "driving_license"            varchar(255)  not null default '',
    "status"                     varchar(255)  not null default '',
    "date_of_birth"              timestamptz,
    "bank_number"                varchar(255)  not null default '',
    "bank_owner"                 varchar(255)  not null default '',
    "bank_name"                  varchar(255)  not null default '',
    "qr_code_url"                varchar(1023) not null default '',
    "created_at"                 timestamptz            DEFAULT (now()),
    "updated_at"                 timestamptz            DEFAULT (now())
);


create unique index unique_email on accounts (email) where email != '';
create unique index unique_identification_card_number on accounts (identification_card_number) where identification_card_number != '';

create table car_models
(
    "id"              serial primary key,
    "brand"           varchar(255) not null default '',
    "model"           varchar(255) not null default '',
    "year"            bigint       not null default 0,
    "number_of_seats" bigint       not null default 0,
    "based_price"     bigint       not null default 0,
    "created_at"      timestamptz           DEFAULT (now()),
    "updated_at"      timestamptz           DEFAULT (now())
);

create table partner_contract_rules
(
    "id"                      serial primary key,
    "revenue_sharing_percent" numeric(3, 1) not null default 0.0,
    "max_warning_count"       bigint        not null default 0,
    "created_at"              timestamptz            DEFAULT (now()),
    "updated_at"              timestamptz            DEFAULT (now())
);

create table cars
(
    "id"                       serial primary key,
    "partner_id"               bigint references accounts (id),
    "car_model_id"             bigint references car_models (id),
    "license_plate"            varchar(255)  not null default '' unique,
    "parking_lot"              varchar(255)  not null default '',
    "description"              varchar(1023) not null default '',
    "fuel"                     varchar(255)  not null default '',
    "motion"                   varchar(255)  not null default '',
    "price"                    bigint        not null default 0,
    "status"                   varchar(255)  not null default '',
    "period"                   bigint        not null default 0,
    "partner_contract_rule_id" bigint references partner_contract_rules (id),
    "bank_name"                varchar(255)  not null default '',
    "bank_number"              varchar(255)  not null default '',
    "bank_owner"               varchar(255)  not null default '',
    "start_date"               timestamptz            DEFAULT (now()),
    "end_date"                 timestamptz            DEFAULT (now()),
    "partner_contract_url"     varchar(1023) not null default '',
    "partner_contract_status"  varchar(256)  not null default '',
    "warning_count"            bigint        not null default 0,
    "created_at"               timestamptz            DEFAULT (now()),
    "updated_at"               timestamptz            DEFAULT (now())
);

create table "notifications"
(
    "id"         serial primary key,
    "account_id" bigint references accounts (id),
    "content"    varchar(1023) not null default '',
    "url"        varchar(1023) not null default '',
    "title"      varchar(255)  not null default '',
    "status"     varchar(255)  not null default '',
    "created_at" timestamptz            DEFAULT (now()),
    "updated_at" timestamptz            DEFAULT (now())
);

create table "car_images"
(
    "id"         serial primary key,
    "car_id"     bigint references cars (id),
    "url"        varchar(1023) not null default '',
    "status"     varchar(255)  not null default '',
    "category"   varchar(255)  not null default '',
    "created_at" timestamptz            DEFAULT (now()),
    "updated_at" timestamptz            DEFAULT (now())
);

create table "partner_payment_histories"
(
    "id"          serial primary key,
    "partner_id"  bigint references accounts (id),
    "start_date"  timestamptz            DEFAULT (now()),
    "end_date"    timestamptz            DEFAULT (now()),
    "amount"      bigint        not null default 0,
    "status"      varchar(255)  not null default '',
    "payment_url" varchar(1023) not null default '',
    "created_at"  timestamptz            DEFAULT (now()),
    "updated_at"  timestamptz            DEFAULT (now())
);

create table customer_contract_rules
(
    "id"                     serial primary key,
    "insurance_percent"      numeric(3, 1) not null default 0.0,
    "prepay_percent"         numeric(3, 1) not null default 0.0,
    "collateral_cash_amount" bigint        not null default 0,
    "created_at"             timestamptz            DEFAULT (now()),
    "updated_at"             timestamptz            DEFAULT (now())
);

create table "customer_contracts"
(
    "id"                         serial primary key,
    "customer_id"                bigint references accounts (id),
    "car_id"                     bigint references cars (id),
    "start_date"                 timestamptz            default (now()),
    "end_date"                   timestamptz            default (now()),
    "status"                     varchar(255)  not null default '',
    "reason"                     varchar(1023) not null default '',
    "rent_price"                 bigint        not null default 0,
    "insurance_amount"           bigint        not null default 0,
    "collateral_type"            varchar(255)  not null default '',
    "is_return_collateral_asset" boolean                default false,
    "url"                        varchar(1023) not null default '',
    "bank_name"                  varchar(255)  not null default '',
    "bank_number"                varchar(255)  not null default '',
    "bank_owner"                 varchar(255)  not null default '',
    "customer_contract_rule_id"  bigint references customer_contract_rules (id),
    "feedback_content"           varchar(1023) not null default '',
    "feedback_rating"            bigint        not null default 0,
    "feedback_status"            varchar(255)  not null default '',
    "created_at"                 timestamptz            DEFAULT (now()),
    "updated_at"                 timestamptz            DEFAULT (now())
);

create table "customer_contract_images"
(
    "id"                   serial primary key,
    "customer_contract_id" bigint references customer_contracts (id),
    "url"                  varchar(1023) not null default '',
    "category"             varchar(255)  not null default '',
    "status"               varchar(255)  not null default '',
    "created_at"           timestamptz            DEFAULT (now()),
    "updated_at"           timestamptz            DEFAULT (now())
);

create table "customer_payments"
(
    "id"                   serial primary key,
    "customer_contract_id" bigint references customer_contracts (id),
    "payment_type"         varchar(255)  not null default '',
    "payment_url"          varchar(1023) not null default '',
    "amount"               bigint        not null default 0,
    "note"                 varchar(1023) not null default '',
    "status"               varchar(255)  not null default '',
    "created_at"           timestamptz            DEFAULT (now()),
    "updated_at"           timestamptz            DEFAULT (now())
);

create table garage_configs
(
    "id"         serial primary key,
    "type"       varchar(255) not null default '',
    "maximum"    bigint       not null default 0,
    "created_at" timestamptz           DEFAULT (now()),
    "updated_at" timestamptz           DEFAULT (now())
);

create table conversations
(
    "id"         serial primary key,
    "account_id" bigint references accounts (id),
    "status"     varchar(255) not null default '',
    "created_at" timestamptz           DEFAULT (now()),
    "updated_at" timestamptz           DEFAULT (now())
);

create table messages
(
    "id"              serial primary key,
    "conversation_id" bigint references conversations (id),
    "sender"          bigint references accounts (id),
    "content"         varchar(1023) not null default '',
    "created_at"      timestamptz            DEFAULT (now()),
    "updated_at"      timestamptz            DEFAULT (now())
);

create table driving_license_images
(
    "id"         serial primary key,
    "account_id" bigint references accounts (id),
    "url"        varchar(1023) not null default '',
    "status"     varchar(255)  not null default '',
    "created_at" timestamptz            DEFAULT (now()),
    "updated_at" timestamptz            DEFAULT (now())
);

create table partner_payment_customer_contracts
(
    "id"                         serial primary key,
    "partner_payment_history_id" bigint references partner_payment_histories (id),
    "customer_contract_id"       bigint references customer_contracts (id),
    "created_at"                 timestamptz DEFAULT (now()),
    "updated_at"                 timestamptz DEFAULT (now())
);

insert into roles(role_name, role_code)
values ('admin', 'AD');
insert into roles(role_name, role_code)
values ('customer', 'CS');
insert into roles(role_name, role_code)
values ('partner', 'PN');
insert into roles(role_name, role_code)
values ('technician', 'TN');

-- insert into accounts(role_id, phone_number, password, status)
-- values (1, 'admin', 'JDJhJDA0JHNrSmNTRmdpQmVGaXp0SVE1SnVUcHU5ZC5UQ0VkeWRQRmx2VHFPUkF5NzRTRnVrcFVXeWd1', 'active');
--
-- insert into accounts(role_id, phone_number, password, status)
-- values (4, 'tech', 'JDJhJDA0JEVOY211WnhIbnVWeHIua3l2YzVkNy53eUI1LnVHdm9uaWFGVERWbWJFUjlQVWNGb2FDdFBl', 'active');
--
-- insert into customer_contract_rules(insurance_percent, prepay_percent, collateral_cash_amount)
-- values (10.0, 30.0, 15000000);
--
-- insert into partner_contract_rules(revenue_sharing_percent, max_warning_count)
-- values (5, 3);

insert into garage_configs(type, maximum)
values ('MAX_4_SEATS', 10);
insert into garage_configs(type, maximum)
values ('MAX_7_SEATS', 5);
insert into garage_configs(type, maximum)
values ('MAX_15_SEATS', 3);

