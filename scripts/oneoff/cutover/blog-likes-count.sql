-- Duplicate blog likes (#799), read-only. The unique index
-- index_blog_post_likes_on_blog_post_id_and_user_id fails to build while any
-- (blog_post_id, user_id) pair repeats. NULL user_id rows (legacy guests) stay
-- distinct under PostgreSQL's NULL semantics, so they are not duplicates.
SELECT count(*)                  AS duplicate_pairs,
       coalesce(sum(c - 1), 0)   AS rows_to_delete
FROM (
  SELECT blog_post_id, user_id, count(*) AS c
  FROM blog_post_likes
  WHERE user_id IS NOT NULL
  GROUP BY blog_post_id, user_id
  HAVING count(*) > 1
) AS dup;
