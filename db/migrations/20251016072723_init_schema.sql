-- +goose Up
-- +goose StatementBegin

create table if not exists users (
    id serial primary key,
    username varchar(100) unique not null,
    email varchar(100) unique not null,
    password varchar(100) not null,
    created_at timestamptz default now()
);

create table if not exists posts (
    id serial primary key,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    text varchar(1000) not null,
    user_id int not null,
    
    constraint fk_post_user foreign key (user_id) references users(id) on delete cascade
);

create table if not exists likes (
    id serial primary key,
    created_at timestamptz default now(),
    post_id int not null,
    user_id int not null,
    
    constraint fk_like_post foreign key (post_id) references posts(id) on delete cascade,
    constraint fk_like_user foreign key (user_id) references users(id) on delete cascade,
    constraint likes_unique unique (post_id, user_id)
);

-- indexes 
create index if not exists idx_posts_user_id on posts(user_id);
create index if not exists idx_likes_post_id on likes(post_id);
create index if not exists idx_likes_user_id on likes(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists users;
drop table if exists posts;
drop table if exists likes;

-- +goose StatementEnd
