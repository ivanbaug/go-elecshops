CREATE TABLE "store"(
                        "id" BIGSERIAL NOT NULL,
                        "name" VARCHAR(255) NOT NULL,
                        "url" VARCHAR(255) NOT NULL,
                        "country" VARCHAR(255) NOT NULL DEFAULT 'Colombia',
                        "region" VARCHAR(255) NOT NULL DEFAULT 'LATAM',
                        "bad_ping_count" BIGINT NOT NULL DEFAULT 0
                    );
ALTER TABLE
    "store" ADD PRIMARY KEY("id");
CREATE TABLE "product"(
                          "id" BIGSERIAL NOT NULL,
                          "sku" VARCHAR(255) NOT NULL,
                          "description" varchar(255) DEFAULT ''::character varying NOT NULL,
                          "vendor" varchar(255) DEFAULT ''::character varying NOT NULL,
                          "stock" BIGINT DEFAULT 0 NOT NULL,
                          "price" BIGINT DEFAULT 0 NOT NULL,
                          "times_clicked_update" BIGINT NOT NULL DEFAULT 0,
                          "id_store" BIGINT NOT NULL,
                          "last_update" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          "first_update" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          "num_updates" BIGINT NOT NULL DEFAULT 0,
                          "url" varchar(255) DEFAULT ''::character varying NOT NULL
                      );
ALTER TABLE
    "product" ADD PRIMARY KEY("id");
ALTER TABLE
    "product" ADD CONSTRAINT "product_id_store_foreign" FOREIGN KEY("id_store") REFERENCES "store"("id");