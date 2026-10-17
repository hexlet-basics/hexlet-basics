import type { BlogPost, BlogPostInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { blogPostFields } from "@/components/admin/generated/blogPost";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

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

// Fields come from the contract (generated/blogPost.ts); what is listed here is
// only what the schema cannot say.
export function useBlogPostFields(): CrudFieldSpec<BlogPostInput>[] {
  return useGeneratedFields<BlogPostInput>(blogPostFields, {
    // The contract lists state after richBody; the form shows it third. The
    // schema allows null for these, the form insists.
    order: ["name", "slug", "state", "description", "richBody"],
    set: {
      name: { required: true },
      slug: { required: true },
      state: { required: true },
      description: { type: "textarea" },
      richBody: { type: "richtext", required: false },
    },
  });
}
