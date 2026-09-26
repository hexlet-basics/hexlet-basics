-- Collapse duplicate blog likes (#799): keep the lowest id per
-- (blog_post_id, user_id). Likes have no dependents and no counter column, so
-- deleting is the whole collapse. The equality join skips NULL user_id rows,
-- which the unique index treats as distinct anyway. Run with
-- --single-transaction; it is idempotent.
DELETE FROM blog_post_likes AS loser
USING blog_post_likes AS winner
WHERE loser.blog_post_id = winner.blog_post_id
  AND loser.user_id = winner.user_id
  AND loser.id > winner.id;
