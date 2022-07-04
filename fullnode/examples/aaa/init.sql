CREATE TABLE IF NOT EXISTS "zta_node" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "peer_id" text NOT NULL DEFAULT '',
  "addr" TEXT NOT NULL DEFAULT '',
  "port" integer NOT NULL DEFAULT 0,
  "ip" text NOT NULL DEFAULT '',
  "loc" text NOT NULL DEFAULT '',
  "colo" text NOT NULL DEFAULT '',
  "gas_price" NUMBER NOT NULL DEFAULT 0,
  "type" text NOT NULL DEFAULT '',
  "created_at" integer(20) NOT NULL,
  "updated_at" integer(20) NOT NULL
);
