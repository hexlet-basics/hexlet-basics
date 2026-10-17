import { useTranslation } from "react-i18next";
import {
  adminArchiveBlogPostMutation,
  adminGetBlogPostQueryKey,
  adminListBlogPostsQueryKey,
  adminPublishBlogPostMutation,
  adminUnpublishBlogPostMutation,
} from "@/client/@tanstack/react-query.gen";
import type { BlogPost } from "@/client/types.gen";
import { PublicationStatePanel } from "@/components/admin/PublicationStatePanel";
import { useResourceMutation } from "@/hooks/useResourceMutation";

// A blog post's Publication State and the buttons that change it (see
// PublicationStatePanel); this one binds the post's three operations.
export function BlogPostStatePanel({ post }: { post: BlogPost }) {
  const { t } = useTranslation();
  const options = {
    invalidate: [adminListBlogPostsQueryKey(), adminGetBlogPostQueryKey({ path: { id: post.id } })],
    successMessage: t(($) => $.admin.publicationState.changed),
    errorMessage: t(($) => $.admin.publicationState.changeError),
  };
  const publish = useResourceMutation({ mutation: adminPublishBlogPostMutation(), ...options });
  const unpublish = useResourceMutation({ mutation: adminUnpublishBlogPostMutation(), ...options });
  const archive = useResourceMutation({ mutation: adminArchiveBlogPostMutation(), ...options });
  const state = post.state ?? "draft";

  return (
    <PublicationStatePanel
      id={post.id}
      title={t(($) => $.models.attributes.blog_post.state)}
      state={state}
      stateLabel={t(($) => $.models.attributes.blog_post["state/values"][state])}
      mutations={{ publish, unpublish, archive }}
    />
  );
}
