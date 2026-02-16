CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS products(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    product_id INT NOT NULL UNIQUE ,
    name VARCHAR NOT NULL ,
    product_code VARCHAR NOT NULL ,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS risk_types(
   id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
   risk_type_id INT NOT NULL UNIQUE ,
   name VARCHAR NOT NULL ,
   risk_category VARCHAR NOT NULL ,
   risk_type_code VARCHAR NOT NULL ,
   description TEXT,
   created_at TIMESTAMP DEFAULT NOW()
);