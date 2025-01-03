begin;

create index if not exists idx_promotions_id on promotions (id, deleted_at);
create index if not exists idx_promotions_uuid on promotions (uuid, deleted_at);
create index if not exists idx_promotions_shop_id on promotions (id, deleted_at, shop_id);

create index if not exists idx_products on products (promotion_id);

create index if not exists idx_variants on product_variants (product_id);

commit;