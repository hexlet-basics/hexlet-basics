import { useTranslation } from "react-i18next";
import type { BlogPost, BlogPostInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";

// Blog-post form values ARE the generated request body (BlogPostInput).
// richBody is trusted editor HTML behind the engine's richtext field (Mantine
// RTE); coverAttachmentId is deferred with the attachment wave. The related
// courses live on their own endpoint and get a separate panel on the edit page.

export const emptyBlogPost: BlogPostInput = {
  name: "",
  slug: "",
  description: "",
  state: "draft",
  richBody: "",
  coverAttachmentId: null,
};

export function blogPostToForm(post: BlogPost): BlogPostInput {
  return {
    name: post.name ?? "",
    slug: post.slug ?? "",
    description: post.description ?? "",
    state: post.state,
    richBody: post.richBodyHtml,
    coverAttachmentId: null,
  };
}

export function useBlogPostFields(): CrudFieldSpec<BlogPostInput>[] {
  const { t } = useTranslation();
  return [
    {
      name: "name",
      label: t(($) => $.models.attributes.base.name),
      required: true,
    },
    {
      name: "slug",
      label: t(($) => $.models.attributes.base.slug),
      required: true,
    },
    {
      name: "state",
      label: t(($) => $.models.attributes.blog_post.state),
      type: "select",
      required: true,
      options: [
        {
          value: "draft",
          label: t(($) => $.models.attributes.blog_post["state/values"].draft),
        },
        {
          value: "published",
          label: t(($) => $.models.attributes.blog_post["state/values"].published),
        },
        {
          value: "archived",
          label: t(($) => $.models.attributes.blog_post["state/values"].archived),
        },
      ],
    },
    {
      name: "description",
      label: t(($) => $.models.attributes.base.description),
      type: "textarea",
    },
    {
      name: "richBody",
      label: t(($) => $.models.attributes.blog_post.rich_body),
      type: "richtext",
    },
  ];
}
