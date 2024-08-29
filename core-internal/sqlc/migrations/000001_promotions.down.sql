begin;

drop table if exists promotion.promotions;
drop index if exists promotion.idx_promotions_id;
drop index if exists promotion.idx_promotions_uuid;


drop table if exists promotion.products;
drop index if exists promotion.idx_products;

drop table if exists promotion.product_variants;
drop index if exists promotion.idx_variants;

commit;