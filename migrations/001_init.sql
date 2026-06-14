create extension if not exists btree_gist;

create table if not exists event_types (
    id text primary key,
    owner_id text not null,
    name text not null,
    description text not null,
    duration_minutes integer not null check (duration_minutes > 0)
);

create table if not exists bookings (
    id text primary key,
    event_type_id text not null references event_types(id) on delete cascade,
    start_at timestamptz not null,
    end_at timestamptz not null,
    created_at timestamptz not null,
    constraint bookings_non_overlapping
        exclude using gist (tstzrange(start_at, end_at, '[)') with &&)
);

create index if not exists bookings_event_type_id_idx on bookings (event_type_id);
create index if not exists bookings_start_at_idx on bookings (start_at);
