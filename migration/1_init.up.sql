create table roles
(
    "id"         serial primary key,
    "role_name"  varchar(255) not null default '',
    "created_at" TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);
insert into roles(role_name)
values ('admin');
insert into roles(role_name)
values ('customer');
insert into roles(role_name)
values ('partner');

create table accounts
(
    "id"                         serial primary key,
    "role_id"                    bigint references roles (id),
    "first_name"                 varchar(255) not null default '',
    "last_name"                  varchar(255) not null default '',
    "phone_number"               varchar(255) not null default '' unique,
    "email"                      varchar(255) not null default '' unique,
    "identification_card_number" varchar(255) not null default '' unique,
    "password"                   varchar(255) not null default '',
    "avatar_url"                 varchar(255) not null default '',
    "status"                     varchar(255) not null default '',
    "created_at"                 TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "updated_at"                 TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

create table partners
(
    "id"         serial primary key,
    "account_id" bigint references accounts (id),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table customers
(
    "id"              serial primary key,
    "account_id"      bigint references accounts (id),
    "driving_license" varchar(255) not null default '',
    "created_at"      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "updated_at"      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

create table car_models
(
    "id"              serial primary key,
    "brand"           varchar(255) not null default '',
    "model"           varchar(255) not null default '',
    "number_of_seats" bigint       not null default 0,
    "created_at"      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "updated_at"      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

create table cars
(
    "id"            serial primary key,
    "partner_id"    bigint references partners (id),
    "car_model_id"  bigint references car_models (id),
    "license_plate" varchar(255)  not null default '',
    "description"   varchar(1023) not null default '',
    "fuel"          varchar(255)  not null default '',
    "motion"        varchar(255)  not null default '',
    "price"         bigint        not null default 0,
    "status"        varchar(255)  not null default '',
    "created_at"    TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
    "updated_at"    TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
);

create table documents
(
    "id"         serial primary key,
    "account_id" bigint references accounts (id),
    "url"        varchar(1023) not null default '',
    "extension"  varchar(255)  not null default '',
    "category"   varchar(255)  not null default '',
    "status"     varchar(255)  not null default '',
    "created_at" TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
);

create table otps
(
    "id"            serial primary key,
    "account_email" varchar(255) references accounts (email),
    "otp"           varchar(20)  not null default '',
    "status"        varchar(255) not null default '',
    "otp_type"      varchar(255) not null default '',
    "expired_at"    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "created_at"    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "updated_at"    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

create table "notifications"
(
    "id"         serial primary key,
    "account_id" bigint references accounts (id),
    "content"    varchar(1023) not null default '',
    "url"        varchar(1023) not null default '',
    "title"      varchar(255)  not null default '',
    "status"     varchar(255)  not null default '',
    "created_at" TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
);

create table "payment_informations"
(
    "id"          serial primary key,
    "account_id"  bigint references accounts (id),
    "bank_number" varchar(255)  not null default '',
    "bank_owner"  varchar(255)  not null default '',
    "bank_name"   varchar(255)  not null default '',
    "qr_code_url" varchar(1023) not null default '',
    "created_at"  TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
);

create table "car_documents"
(
    "id"          serial primary key,
    "car_id"      bigint references cars (id),
    "document_id" bigint references documents (id),
    "created_at"  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table "partner_contracts"
(
    "id"         serial primary key,
    "car_id"     bigint references cars (id),
    "partner_id" bigint references partners (id),
    "start_date" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "end_date"   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table "partner_payment_histories"
(
    "id"         serial primary key,
    "partner_id" bigint references partners (id),
    "from"       TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "to"         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "amount"     bigint       not null default 0,
    "status"     varchar(255) not null default '',
    "created_at" TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

create table "trips"
(
    "id"          serial primary key,
    "customer_id" bigint references customers (id),
    "car_id"      bigint references cars (id),
    "start_date"  TIMESTAMP              default CURRENT_TIMESTAMP,
    "end_date"    TIMESTAMP              default CURRENT_TIMESTAMP,
    "status"      varchar(255)  not null default '',
    "reason"      varchar(1023) not null default '',
    "created_at"  TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
);

create table "trip_payments"
(
    "id"           serial primary key,
    "trip_id"      bigint references trips (id),
    "payment_type" varchar(255)  not null default '',
    "amount"       bigint        not null default 0,
    "note"         varchar(1023) not null default '',
    "status"       varchar(255)  not null default '',
    "created_at"   TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
    "updated_at"   TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
);

create table "trip_payment_documents"
(
    "id"              serial primary key,
    "trip_payment_id" bigint references trip_payments (id),
    "document_id"     bigint references documents (id),
    "created_at"      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at"      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table "trip_contracts"
(
    "id"                         serial primary key,
    "trip_id"                    bigint references trips (id),
    "collateral_type"            varchar(255) not null default '',
    "status"                     varchar(255) not null default '',
    "is_return_collateral_asset" boolean               default false,
    "created_at"                 TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "updated_at"                 TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

create table "trip_feedbacks"
(
    "id"         serial primary key,
    "trip_id"    bigint references trips (id),
    "content"    varchar(1023) not null default '',
    "rating"     bigint        not null default 0,
    "status"     bigint        not null default 0,
    "created_at" TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP              DEFAULT CURRENT_TIMESTAMP
);

create table "trip_documents"
(
    "id"          serial primary key,
    "trip_id"     bigint references trips (id),
    "document_id" bigint references documents (id),
    "created_at"  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at"  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
