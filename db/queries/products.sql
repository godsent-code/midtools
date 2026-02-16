-- name: CreateProducts :exec
INSERT INTO products(product_id, name, product_code, description)
 SELECT unnest(@product_id::int[]),unnest(@name::text[]),unnest(@product_code::text[]),
        unnest(@description::text[]) ON CONFLICT DO NOTHING ;

-- name: GetProducts :many
SELECT * FROM products;


-- name: CreateRiskType :exec
INSERT INTO risk_types(risk_type_id, name, risk_category, risk_type_code, description)
SELECT unnest(@risk_type_id::int[]),unnest(@name::text[]),unnest(@risk_category::text[]),
       unnest(@risk_type_code::text[]),unnest(@description::text[]) ON CONFLICT DO NOTHING;

-- name: GetRiskType :many
SELECT * FROM risk_types;