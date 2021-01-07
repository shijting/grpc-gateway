**用户（users）**

|**字段名称**|**类型**|**约束**|**默认值**|**描述**|
|:----|:----|:----|:----|:----|
|id|int|pk|    |自增主键|
|phone_number|char(11)|unique,not null|    |手机号码|
|password|varchar(255)|not null|空|登录密码|
|nickname|varchar(20)|not null|空|昵称|
|avatar|varchar(255)|not null|空|头像|
|last_login_ip|varchar(15)|not  null|空|最后登录ip|
|last_login_at|timestamptz|    |    |最后登录时间|
|status|smallint|not null|1|状态(1正常，2禁用)|
|created_at|timestamptz|not null|current_timestamp|创建时间|
|updated_at|timestamptz|    |    |更新时间|

**设备(cameras)**

|**字段名称**|**类型**|**约束**|**默认值**|**描述**|
|:----|:----|:----|:----|:----|
|id|int|pk|    |自增主键|
|no|varchar(20)|not null,unique|    |设备编号|
|user_id|int|index|    |管理id|
|name|varchar(60)|not  null|空|设备名称|
|model|varchar(30)|not null|    |设备型号|
|mac|varchar(24)|not null|    |mac地址|
|ip|varchar(15)|not null|空|ip地址|
|port|smallint|not null|0|端口|
|password|varchar(255)|not null|空|设备密码|
|is_alarm|bool|not null|false|是否开启报警|
|status|smallint|not null|1|状态(1正常2禁用)|
|created_at|timestamptz|not null|current_timestamp|创建时间|
|updated_at|timestamptz|    |    |更新时间|

**我的设备(user_cameras)**

|**字段名称**|**类型**|**约束**|**默认值**|**描述**|
|:----|:----|:----|:----|:----|
|id|int|pk|    |自增主键|
|user_id|int|not null,index|    |用户id|
|camera_id|int|not nul,index|    |设备id(外键：关联cameras表)|
|permissions|int|    |    |权限|
|is_admin|bool|not null|false|    |
|created_at|timestamptz|not null|current_timestamp|创建时间|
|updated_at|timestamptz|    |    |更新时间|

**消息（camera_messages）**

|**字段名称**|**类型**|**约束**|**默认值**|**描述**|
|:----|:----|:----|:----|:----|
|id|int|pk|    |自增主键|
|camera_id|int|not null,index|    |设备id(外键：关联my_cameras表)|
|image_url|varchar(255)|not null|空|封面|
|title|varchar(255)|not null|空|标题|
|video_url|varchar(255)|not null|空|视频地址|
|is_read|bool|not null|false|    |
|created_at|timestamptz|not null|current_timestamp|创建时间|
|updated_at|timestamptz|    |    |更新时间|


**反馈（feedback）**

|**字段名称**|**类型**|**约束**|**默认值**|**描述**|
|:----|:----|:----|:----|:----|
|id|int|pk|    |自增主键|
|content|text|not null|    |内容|
|phone_number|char(11)|not null|    |联系手机|
|created_at|timestamp|not null|current_timestamp|创建时间|























