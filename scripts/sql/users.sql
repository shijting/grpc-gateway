-- 用户表
create table users (
	id serial not null  primary key,
	phone_number char(11) not null ,
	avatar varchar(255) not null default '',
	last_login_date timestamptz null,
	last_login_ip varchar(18) not null default '',
	status smallint not null default 1,
	created_at timestamptz not null default CURRENT_TIMESTAMP,
	updated_at timestamptz null,
	unique("phone_number")
);
-- 反馈表
create table feedback (
	id serial not null primary key,
	content text not null,
	phone_number char(11) not null,
	created_at timestamptz not null default current_timestamp
);
-- 设备
create table cameras (
	id serial not null primary key,
	camera_name varchar(60) not null default '',
	camera_id varchar(50) not null,
	mac_addr varchar(100) not null,
	status smallint not null default 1,
	created_at timestamptz not null default current_timestamp,
	updated_at timestamptz  null
);
create index "idx_cameras_machine_id" on cameras("machine_id");

-- 我的设备
create table my_cameras (
	id serial not null primary key,
	camera_name varchar(60) not null default '',
	camera_id varchar(50) not null,
	user_id int not null,
	sharer_phone char(11) not null default '',
	status smallint not null default 1,
	created_at timestamptz not null default current_timestamp,
	updated_at timestamptz  null
);
create index "idx_my_cameras_user_id" on my_cameras("user_id");





-- 消息
create table camera_messages (
	id serial not null primary key,
	camera_id int not null,
	to_user_id int[] null,
	image_url varchar(255) not null default '',
	video_url varchar(255) not null default '',
	title varchar(60) not null default '',
	created_at timestamptz not null default CURRENT_TIMESTAMP,
	updated_at timestamptz null
);

create index "idx_camera_messages_camera_id" on camera_messages("camera_id")
