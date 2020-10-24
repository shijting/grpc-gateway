create table users (
	id serial not null  primary key,
	phone_number char(11) not null ,
	password varchar(255) not null,
	last_login_date timestamptz null,
	last_login_ip varchar(18) not null default '',
	status smallint not null default 1,
	created_at timestamptz not null default CURRENT_TIMESTAMP,
	updated_at timestamptz null,
	unique("phone_number")
);