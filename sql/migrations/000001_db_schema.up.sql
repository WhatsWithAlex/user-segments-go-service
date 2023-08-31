CREATE TYPE "op" AS ENUM (
  'remove',
  'add'
);

CREATE TABLE "segments" (
  "id" serial PRIMARY KEY,
  "slug" varchar UNIQUE NOT NULL,
  "auto_prob" real NOT NULL DEFAULT 0.0 CHECK (auto_prob BETWEEN 0.0 AND 1.0)
);

CREATE TABLE "user_segments" (
  "id" serial PRIMARY KEY,
  "user_id" int NOT NULL,
  "segment_id" int NOT NULL,
  "remove_at" timestamptz DEFAULT null
);

CREATE TABLE "operations" (
  "id" serial PRIMARY KEY,
  "user_id" int NOT NULL,
  "segment_slug" varchar NOT NULL,
  "op_type" op NOT NULL,
  "done_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "segments" ("slug");

CREATE INDEX ON "user_segments" ("user_id");

CREATE UNIQUE INDEX ON "user_segments" ("user_id", "segment_id");

CREATE INDEX ON "operations" ("user_id", "done_at");

COMMENT ON COLUMN "segments"."slug" IS 'Segment unique name.';

COMMENT ON COLUMN "segments"."auto_prob" IS 'Probability of user to be added automatically to the segment, [0..1].';

ALTER TABLE "user_segments" ADD FOREIGN KEY ("segment_id") REFERENCES "segments" ("id") ON DELETE CASCADE;
