CREATE TABLE "message" (
                           "message_id" bigserial PRIMARY KEY,
                           "from_number" varchar NOT NULL,
                           "to_number" varchar NOT NULL,
                           "message_text" varchar NOT NULL,
                           "sent_datetime" timestamptz NOT NULL DEFAULT (now()),
                           "contact_id" bigserial NOT NULL
);

CREATE TABLE "contact" (
                           "contact_id" bigserial PRIMARY KEY,
                           "first_name" varchar NOT NULL,
                           "last_name" varchar NOT NULL,
                           "profile_photo" bytea NOT NULL,
                           "phone_number" varchar UNIQUE NOT NULL,
                           "username" varchar UNIQUE NOT NULL,
                           "hashed_password" varchar NOT NULL

);

CREATE TABLE "message_group" (
                                 "group_id" bigserial PRIMARY KEY,
                                 "group_name" varchar NOT NULL
);

CREATE TABLE "group_member" (
                                "contact_id" bigserial NOT NULL,
                                "group_id" bigserial NOT NULL,
                                "joined_datetime" timestamptz NOT NULL DEFAULT (now()),
                                "left_datetime" timestamptz
);

ALTER TABLE "message" ADD FOREIGN KEY ("contact_id") REFERENCES "contact" ("contact_id");

ALTER TABLE "group_member" ADD FOREIGN KEY ("contact_id") REFERENCES "contact" ("contact_id");

ALTER TABLE "group_member" ADD FOREIGN KEY ("group_id") REFERENCES "message_group" ("group_id");
