CREATE TABLE IF NOT EXISTS "servers" (
  "id" VARCHAR(32) NOT NULL PRIMARY KEY,
  "tag" VARCHAR(255) NOT NULL,
  "host" VARCHAR(255) NOT NULL,
  "address" VARCHAR(255) NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  -- healthy, unhealthy, unknown
  "status" VARCHAR(12) NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "servers_tag" ON "servers" ("tag");

CREATE INDEX IF NOT EXISTS "servers_host" ON "servers" ("host");

CREATE INDEX IF NOT EXISTS "servers_address" ON "servers" ("address");