-- Adminer 4.8.0 PostgreSQL 11.5 (Debian 11.5-3.pgdg90+1) dump

DROP TABLE IF EXISTS "followers";
CREATE TABLE "public"."followers" (
    "following_user_id" character varying(80) NOT NULL,
    "followed_user_id" character varying(80) NOT NULL,
    CONSTRAINT "follower_following_user_id_followed_user_id" PRIMARY KEY ("following_user_id", "followed_user_id")
) WITH (oids = false);

COMMENT ON COLUMN "public"."followers"."following_user_id" IS 'フォローユーザID';

COMMENT ON COLUMN "public"."followers"."followed_user_id" IS '被フォローユーザID';


DROP TABLE IF EXISTS "replys";
CREATE TABLE "public"."replys" (
    "replying_tweet_id" character varying(36) NOT NULL,
    "replyed_tweet_id" character varying(36) NOT NULL,
    CONSTRAINT "reply_replying_tweet_id_replyed_tweet_id" PRIMARY KEY ("replying_tweet_id", "replyed_tweet_id")
) WITH (oids = false);

COMMENT ON COLUMN "public"."replys"."replying_tweet_id" IS 'リプライツイートID';

COMMENT ON COLUMN "public"."replys"."replyed_tweet_id" IS '被リプライツイートID';


DROP TABLE IF EXISTS "tweets";
CREATE TABLE "public"."tweets" (
    "tweet_id" character varying(36) NOT NULL,
    "tweet_user_id" character varying(80) NOT NULL,
    "tweeted_time" timestamp NOT NULL,
    "contents" character varying(140) NOT NULL,
    "is_deleted" boolean NOT NULL,
    CONSTRAINT "tweet_tweet_id" PRIMARY KEY ("tweet_id"),
    CONSTRAINT "tweet_tweet_user_id_tweeted_time" UNIQUE ("tweet_user_id", "tweeted_time")
) WITH (oids = false);

COMMENT ON COLUMN "public"."tweets"."tweet_id" IS '代理キー（UUID）';

COMMENT ON COLUMN "public"."tweets"."tweet_user_id" IS 'ツイートユーザID';

COMMENT ON COLUMN "public"."tweets"."tweeted_time" IS 'ツイート日時';

COMMENT ON COLUMN "public"."tweets"."contents" IS 'ツイート内容';

COMMENT ON COLUMN "public"."tweets"."is_deleted" IS '削除されたらtrue';


DROP TABLE IF EXISTS "users";
CREATE TABLE "public"."users" (
    "user_id" character varying(80) NOT NULL,
    "name" character varying(80) NOT NULL,
    "password_hash" character varying(128) NOT NULL,
    CONSTRAINT "user_user_id" PRIMARY KEY ("user_id")
) WITH (oids = false);

COMMENT ON COLUMN "public"."users"."user_id" IS 'ユーザID';

COMMENT ON COLUMN "public"."users"."name" IS 'ユーザ表示名';

COMMENT ON COLUMN "public"."users"."password_hash" IS 'SHA512ハッシュ値';


ALTER TABLE ONLY "public"."followers" ADD CONSTRAINT "follower_followed_user_id_fkey" FOREIGN KEY (followed_user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY "public"."followers" ADD CONSTRAINT "follower_following_user_id_fkey" FOREIGN KEY (following_user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

ALTER TABLE ONLY "public"."replys" ADD CONSTRAINT "reply_replyed_tweet_id_fkey" FOREIGN KEY (replyed_tweet_id) REFERENCES tweets(tweet_id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY "public"."replys" ADD CONSTRAINT "reply_replying_tweet_id_fkey" FOREIGN KEY (replying_tweet_id) REFERENCES tweets(tweet_id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

ALTER TABLE ONLY "public"."tweets" ADD CONSTRAINT "tweet_tweet_user_id_fkey" FOREIGN KEY (tweet_user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

-- 2021-05-06 22:11:04.253153+00
