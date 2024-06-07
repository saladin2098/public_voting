CREATE TABLE party (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50),
    slogan VARCHAR(50),
    open_date TIMESTAMP NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE public (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    birthday TIMESTAMP NOT NULL,
    gender VARCHAR(1) NOT NULL,
    nation VARCHAR(30) NOT NULL,
    party_id UUID REFERENCES party(id)
);