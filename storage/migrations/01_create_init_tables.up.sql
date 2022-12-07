BEGIN;

CREATE TABLE "user" (
	"id" CHAR(36) NOT NULL PRIMARY KEY,
	"username" VARCHAR NOT NULL, 
	"password" VARCHAR(255) NOT NULL,
	"user_type" VARCHAR NOT NULL,
	"created_at" TIMESTAMP DEFAULT now(),
	"updated_at" TIMESTAMP,
	"deleted_at" TIMESTAMP
);

COMMIT;