CREATE TABLE IF NOT EXISTS timestamps {
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL, 
    stamp_type varchar(255) NOT NULL,
    time timestamp(0) with time zone NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version int NOT NULL DEFAULT 0
}