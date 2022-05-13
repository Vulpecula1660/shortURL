CREATE TABLE "urls" (
  "id" bigserial PRIMARY KEY,
  "origin_url" varchar NOT NULL,
  "short_url" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "urls" ("short_url");
