CREATE TABLE "users" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password" varchar NOT NULL,
  "telephone_number" varchar NOT NULL,
  "university" varchar NOT NULL,
  "department" varchar NOT NULL,
  "date_of_birth" timestamp NOT NULL,
  "role" varchar    NOT NULL,
  "active" boolean  NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(NOW()),
  "updated_at" timestamp NOT NULL DEFAULT(NOW()),
  "deleted_at" timestamp
);

CREATE TABLE "teams" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(NOW()),
  "updated_at" timestamp NOT NULL DEFAULT(NOW()),
  "deleted_at" timestamp
);

CREATE TABLE "projects" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(NOW()),
  "updated_at" timestamp NOT NULL DEFAULT(NOW()),
  "deleted_at" timestamp
);

CREATE TABLE "announcements" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar NOT NULL,
  "body" varchar NOT NULL,
  "authorId" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(NOW()),
  "updated_at" timestamp NOT NULL DEFAULT(NOW()),
  "deleted_at" timestamp
);

CREATE TABLE "project_users" (
"id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "project_id" int NOT NULL,
  "user_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(NOW()),
  "updated_at" timestamp NOT NULL DEFAULT(NOW()),
  "deleted_at" timestamp
);

CREATE TABLE "team_users" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "team_id" int NOT NULL,
  "user_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(NOW()),
  "updated_at" timestamp NOT NULL DEFAULT(NOW()),
  "deleted_at" timestamp
);

CREATE TABLE "team_leads" (
"id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "team_id" int NOT NULL,
  "user_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(NOW()),
  "updated_at" timestamp NOT NULL DEFAULT(NOW()),
  "deleted_at" timestamp
);

CREATE TABLE "team_projects" (
"id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "team_id" int NOT NULL,
  "project_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(NOW()),
  "updated_at" timestamp NOT NULL DEFAULT(NOW()),
  "deleted_at" timestamp
);

CREATE TABLE "project_leads" (
"id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "project_id" int NOT NULL,
  "user_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT(NOW()),
  "updated_at" timestamp NOT NULL DEFAULT(NOW()),
  "deleted_at" timestamp
);

ALTER TABLE "announcements" ADD FOREIGN KEY ("authorId") REFERENCES "users" ("id");

ALTER TABLE "project_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "project_users" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "team_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "team_users" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("id");

ALTER TABLE "team_leads" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("id");

ALTER TABLE "team_leads" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "team_projects" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("id");

ALTER TABLE "team_projects" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "project_leads" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "project_leads" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
