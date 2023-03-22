create table socialDB.friends
(
    user_id   int null,
    friend_id int null,
    id        int auto_increment
        primary key,
    constraint user_id
        foreign key (user_id) references socialDB.users (user_id)
            on delete cascade
);

