CREATE VIEW vw_product
AS
SELECT p.id,
       p.sku,
       p.description,
       p.vendor,
       p.stock,
       p.price,
       p.times_clicked_update,
       p.id_store,
       p.last_update,
       p.first_update,
       p.num_updates,
       p.url,
       s.name AS store_name,
       s.country
FROM product p
         INNER JOIN store s
                    ON p.id_store = s.id;