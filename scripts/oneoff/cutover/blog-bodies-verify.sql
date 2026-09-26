-- Published posts still without a body. Expected: 0.
SELECT count(*) AS published_without_body
FROM blog_posts
WHERE state = 'published' AND rich_body = '';
