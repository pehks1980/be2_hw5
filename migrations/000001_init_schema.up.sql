CREATE TABLE "users" (
                         "user_id" INT,
                         "name" VARCHAR,
                         "age" INT,
                         "spouse" INT
);
CREATE UNIQUE INDEX "users_user_id" ON "users" ("user_id");

CREATE TABLE "activities" (
                              "user_id" INT,
                              "name" VARCHAR,
                              "date" TIMESTAMP
);
CREATE UNIQUE INDEX "activities_user_id" ON "activities" ("user_id");
