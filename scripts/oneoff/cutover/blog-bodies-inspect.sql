-- The ActionText bodies that become blog_posts.rich_body, read-only. The Go
-- stack stores and serves plain trusted HTML, so an <action-text-attachment>
-- (an sgid only Rails can resolve) renders as nothing: those bodies need a hand
-- fix after the copy. img/iframe/table are the other tags legacy's own check
-- (legacy/app/models/blog_post.rb) singles out.
SELECT bp.slug,
       bp.state,
       bp.locale,
       length(art.body)                                   AS body_chars,
       art.body ~ '<action-text-attachment'               AS has_attachment,
       art.body ~ '<(img|iframe|table)'                   AS has_img_iframe_table,
       bp.rich_body <> ''                                 AS already_copied
FROM blog_posts AS bp
LEFT JOIN action_text_rich_texts AS art
  ON art.record_type = 'BlogPost'
 AND art.name = 'rich_body'
 AND art.record_id = bp.id
ORDER BY bp.state, bp.slug;
