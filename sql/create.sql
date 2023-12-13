CREATE TABLE IF NOT EXISTS "services" (
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

CREATE INDEX IF NOT EXISTS "services_tag" ON "services" ("tag");

CREATE INDEX IF NOT EXISTS "services_host" ON "services" ("host");

CREATE INDEX IF NOT EXISTS "services_address" ON "services" ("address");