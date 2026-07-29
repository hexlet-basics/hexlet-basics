-- Add the trusted editor HTML stored directly on blog posts. The five
-- production bodies are populated manually before the Go cutover, so this
-- migration intentionally contains no ActionText data copy.
ALTER TABLE "blog_posts" ADD COLUMN "rich_body" text NOT NULL DEFAULT '';
