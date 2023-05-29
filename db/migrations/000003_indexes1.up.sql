-- Can only be done with OWNER of table privileges
-- O

CREATE INDEX lower_sku_idx ON product ((lower(sku)));
CREATE INDEX lower_description_idx ON product ((lower(description)));