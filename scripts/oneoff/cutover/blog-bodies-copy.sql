-- Copy the ActionText bodies into blog_posts.rich_body (the column comes from
-- migrations/20260729130000_add_blog_post_rich_body.sql, which copies no data
-- on purpose). Only empty rich_body is filled, so a re-run never overwrites a
-- body that was hand-fixed after the first copy. Run with --single-transaction.
UPDATE blog_posts AS bp
SET rich_body = art.body
FROM action_text_rich_texts AS art
WHERE art.record_type = 'BlogPost'
  AND art.name = 'rich_body'
  AND art.record_id = bp.id
  AND art.body IS NOT NULL
  AND bp.rich_body = '';
