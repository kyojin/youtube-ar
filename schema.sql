begin;
create table urls(
    id serial primary key not null,
    url text not null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    status text not null default 'pending'
);

-- update updated_at on update
CREATE FUNCTION urls_update() RETURNS trigger AS $urls_update$
    BEGIN
        NEW.updated_at := current_timestamp;
        RETURN NEW;
    END;
$urls_update$ LANGUAGE plpgsql;

CREATE TRIGGER urls_update BEFORE UPDATE ON urls
    FOR EACH ROW EXECUTE FUNCTION urls_update();
commit