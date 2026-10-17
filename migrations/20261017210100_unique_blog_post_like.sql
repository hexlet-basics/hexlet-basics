-- A person likes a post once. The legacy app checked before it inserted, so two
-- simultaneous clicks could each insert and the pair may repeat. Collapse each
-- repeated (post, user) pair to its earliest like, then make the pair unique so
-- the like is an insert that ignores a conflict. Rows without a user are not
-- pairs: NULLs never collide in a unique index, so they are left alone.
DELETE FROM "blog_post_likes" WHERE "id" IN (
  SELECT "id" FROM (
    SELECT "id", ROW_NUMBER() OVER (
      PARTITION BY "blog_post_id", "user_id" ORDER BY "created_at", "id"
    ) AS "rank"
    FROM "blog_post_likes"
    WHERE "user_id" IS NOT NULL
  ) AS "ranked"
  WHERE "rank" > 1
);
CREATE UNIQUE INDEX "index_blog_post_likes_on_blog_post_id_and_user_id"
ON "blog_post_likes" ("blog_post_id", "user_id");
