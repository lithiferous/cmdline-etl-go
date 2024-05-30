CREATE TABLE "snapshots" (
  "user_name" varchar NOT NULL,
  "store_name" varchar NOT NULL,
  "credit_limit" decimal NOT NULL,
  "snapshot_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);
