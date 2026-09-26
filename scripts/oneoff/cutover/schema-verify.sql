-- Objects the cutover migrations create, one row each. Expected: every
-- present = t.
SELECT 'index_language_members_on_user_id_and_language_id' AS object,
       to_regclass('public.index_language_members_on_user_id_and_language_id') IS NOT NULL AS present
UNION ALL
SELECT 'blog_posts.rich_body',
       EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name = 'blog_posts' AND column_name = 'rich_body')
UNION ALL
SELECT 'river_job', to_regclass('public.river_job') IS NOT NULL
UNION ALL
SELECT 'attachments', to_regclass('public.attachments') IS NOT NULL;
