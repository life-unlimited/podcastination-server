create sequence podcastination_key_seq
    as integer;

create table owners
(
    id        serial
        constraint owners_pk
            primary key,
    name      varchar not null,
    email     varchar not null,
    copyright varchar not null
);

create unique index owners_id_uindex
    on owners (id);

create table podcasts
(
    id             serial
        constraint podcasts_pk
            primary key,
    title          varchar not null,
    subtitle       varchar,
    language       varchar,
    owner_id       integer not null
        constraint podcasts_owners_id_fk
            references owners,
    description    varchar,
    keywords       varchar,
    link           varchar,
    image_location varchar,
    type           varchar,
    key            varchar,
    feed_link      varchar not null
);

create unique index podcasts_id_uindex
    on podcasts (id);

create unique index podcasts_key_uindex
    on podcasts (key);

create table seasons
(
    id             serial
        constraint seasons_pk
            primary key,
    title          varchar not null,
    subtitle       varchar,
    image_location varchar,
    podcast_id     integer not null
        constraint seasons_podcasts_id_fk
            references podcasts,
    num            integer,
    description    varchar,
    key            varchar
);

create unique index seasons_id_uindex
    on seasons (id);

create unique index seasons_key_podcast_id_uindex
    on seasons (key, podcast_id);

create unique index seasons_num_podcast_id_uindex
    on seasons (num, podcast_id);

create table episodes
(
    id             serial
        constraint episodes_pk
            primary key,
    title          varchar               not null,
    subtitle       varchar,
    date           date                  not null,
    author         varchar,
    description    varchar,
    mp3_location   varchar,
    season_id      integer               not null
        constraint episodes_seasons_id_fk
            references seasons,
    num            integer               not null,
    image_location varchar,
    yt_url         varchar,
    mp3_length     integer               not null,
    is_available   boolean default false not null,
    pdf_location   varchar
);

create unique index episodes_id_uindex
    on episodes (id);

create unique index episodes_num_season_id_uindex
    on episodes (num, season_id);

create table podcastination
(
    key   varchar not null
        constraint podcastination_pk
            primary key,
    value varchar
);

create unique index podcastination_key_uindex
    on podcastination (key);

