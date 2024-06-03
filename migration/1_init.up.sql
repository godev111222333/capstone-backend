create table roles
(
    "id"         serial primary key,
    "role_name"  varchar(255) not null default '',
    "role_code"  varchar(255) not null default '',
    "created_at" timestamptz           DEFAULT (now()),
    "updated_at" timestamptz           DEFAULT (now())
);
insert into roles(role_name, role_code)
values ('admin', 'AD');
insert into roles(role_name, role_code)
values ('customer', 'CS');
insert into roles(role_name, role_code)
values ('partner', 'PN');

create table accounts
(
    "id"                         serial primary key,
    "role_id"                    bigint references roles (id),
    "first_name"                 varchar(255) not null default '',
    "last_name"                  varchar(255) not null default '',
    "phone_number"               varchar(255) not null default '',
    "email"                      varchar(255) not null unique,
    "identification_card_number" varchar(255) not null default '',
    "password"                   varchar(255) not null default '',
    "avatar_url"                 varchar(255) not null default '',
    "driving_license"            varchar(255) not null default '',
    "status"                     varchar(255) not null default '',
    "date_of_birth"              timestamptz,
    "created_at"                 timestamptz           DEFAULT (now()),
    "updated_at"                 timestamptz           DEFAULT (now())
);

create unique index unique_phone_number on accounts (phone_number) where phone_number != '';
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

create table cars
(
    "id"            serial primary key,
    "partner_id"    bigint references accounts (id),
    "car_model_id"  bigint references car_models (id),
    "license_plate" varchar(255)  not null default '',
    "parking_lot"   varchar(255)  not null default '',
    "description"   varchar(1023) not null default '',
    "fuel"          varchar(255)  not null default '',
    "motion"        varchar(255)  not null default '',
    "price"         bigint        not null default 0,
    "status"        varchar(255)  not null default '',
    "created_at"    timestamptz            DEFAULT (now()),
    "updated_at"    timestamptz            DEFAULT (now())
);

create table documents
(
    "id"         serial primary key,
    "account_id" bigint references accounts (id),
    "url"        varchar(1023) not null default '',
    "extension"  varchar(255)  not null default '',
    "category"   varchar(255)  not null default '',
    "status"     varchar(255)  not null default '',
    "created_at" timestamptz            DEFAULT (now()),
    "updated_at" timestamptz            DEFAULT (now())
);

create table otps
(
    "id"            serial primary key,
    "account_email" varchar(255) references accounts (email),
    "otp"           varchar(20)  not null default '',
    "status"        varchar(255) not null default '',
    "otp_type"      varchar(255) not null default '',
    "expires_at"    timestamptz           DEFAULT (now()),
    "created_at"    timestamptz           DEFAULT (now()),
    "updated_at"    timestamptz           DEFAULT (now())
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

create table "payment_informations"
(
    "id"          serial primary key,
    "account_id"  bigint references accounts (id),
    "bank_number" varchar(255)  not null default '',
    "bank_owner"  varchar(255)  not null default '',
    "bank_name"   varchar(255)  not null default '',
    "qr_code_url" varchar(1023) not null default '',
    "created_at"  timestamptz            DEFAULT (now()),
    "updated_at"  timestamptz            DEFAULT (now())
);

create table "car_documents"
(
    "id"          serial primary key,
    "car_id"      bigint references cars (id),
    "document_id" bigint references documents (id),
    "created_at"  timestamptz DEFAULT (now()),
    "updated_at"  timestamptz DEFAULT (now())
);

create table "partner_contracts"
(
    "id"         serial primary key,
    "car_id"     bigint references cars (id),
    "partner_id" bigint references accounts (id),
    "start_date" timestamptz DEFAULT (now()),
    "end_date"   timestamptz DEFAULT (now()),
    "created_at" timestamptz DEFAULT (now()),
    "updated_at" timestamptz DEFAULT (now())
);

create table "partner_payment_histories"
(
    "id"         serial primary key,
    "partner_id" bigint references accounts (id),
    "from"       timestamptz           DEFAULT (now()),
    "to"         timestamptz           DEFAULT (now()),
    "amount"     bigint       not null default 0,
    "status"     varchar(255) not null default '',
    "created_at" timestamptz           DEFAULT (now()),
    "updated_at" timestamptz           DEFAULT (now())
);

create table "trips"
(
    "id"               serial primary key,
    "customer_id"      bigint references accounts (id),
    "car_id"           bigint references cars (id),
    "start_date"       timestamptz            default (now()),
    "end_date"         timestamptz            default (now()),
    "status"           varchar(255)  not null default '',
    "reason"           varchar(1023) not null default '',
    "insurance_amount" bigint        not null default 0,
    "created_at"       timestamptz            DEFAULT (now()),
    "updated_at"       timestamptz            DEFAULT (now())
);

create table "trip_payments"
(
    "id"           serial primary key,
    "trip_id"      bigint references trips (id),
    "payment_type" varchar(255)  not null default '',
    "amount"       bigint        not null default 0,
    "note"         varchar(1023) not null default '',
    "status"       varchar(255)  not null default '',
    "created_at"   timestamptz            DEFAULT (now()),
    "updated_at"   timestamptz            DEFAULT (now())
);

create table "trip_payment_documents"
(
    "id"              serial primary key,
    "trip_payment_id" bigint references trip_payments (id),
    "document_id"     bigint references documents (id),
    "created_at"      timestamptz DEFAULT (now()),
    "updated_at"      timestamptz DEFAULT (now())
);

create table "trip_contracts"
(
    "id"                         serial primary key,
    "trip_id"                    bigint references trips (id),
    "collateral_type"            varchar(255) not null default '',
    "status"                     varchar(255) not null default '',
    "is_return_collateral_asset" boolean               default false,
    "created_at"                 timestamptz           DEFAULT (now()),
    "updated_at"                 timestamptz           DEFAULT (now())
);

create table "trip_feedbacks"
(
    "id"         serial primary key,
    "trip_id"    bigint references trips (id),
    "content"    varchar(1023) not null default '',
    "rating"     bigint        not null default 0,
    "status"     bigint        not null default 0,
    "created_at" timestamptz            DEFAULT (now()),
    "updated_at" timestamptz            DEFAULT (now())
);

create table "trip_documents"
(
    "id"          serial primary key,
    "trip_id"     bigint references trips (id),
    "document_id" bigint references documents (id),
    "created_at"  timestamptz DEFAULT (now()),
    "updated_at"  timestamptz DEFAULT (now())
);

create table "sessions"
(
    "id"            uuid primary key,
    "email"         varchar(255)  not null default '',
    "refresh_token" varchar(1023) not null default '',
    "user_agent"    varchar(255)  not null default '',
    "client_ip"     varchar(255)  not null default '',
    "expires_at"    timestamptz   not null,
    "created_at"    timestamptz            DEFAULT (now()),
    "updated_at"    timestamptz            DEFAULT (now())
);

create table garage_configs
(
    "id"         serial primary key,
    "type"       varchar(255) not null default '',
    "maximum"    bigint       not null default 0,
    "created_at" timestamptz           DEFAULT (now()),
    "updated_at" timestamptz           DEFAULT (now())
);
