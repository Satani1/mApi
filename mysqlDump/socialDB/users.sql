create table socialDB.users
(
    user_id int auto_increment
        primary key,
    name    varchar(40) default '' null,
    age     int         default 0  null
);

