import type { BlogPost, BlogPostInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { blogPostFields } from "@/components/admin/generated/blogPost";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

// Blog-post form values ARE the generated request body (BlogPostInput).
// richBody is trusted editor HTML behind the engine's richtext field (Mantine
// RTE); coverAttachmentId is deferred with the attachment wave. The related
// courses live on their own endpoint and get a separate panel on the edit page,
// and so does the Publication State: it is changed by its own operations, never
// by the form.

export const emptyBlogPost: BlogPostInput = {
  coverAttachmentId: null,
  description: "",
  name: "",
  richBody: "",
  slug: "",
};

export function blogPostToForm(post: BlogPost): BlogPostInput {
  return {
    coverAttachmentId: null,
    description: post.description ?? "",
    name: post.name ?? "",
    richBody: post.richBodyHtml,
    slug: post.slug ?? "",
  };
}

// Fields come from the contract (generated/blogPost.ts); what is listed here is
// only what the schema cannot say. The state is not a field: BlogPostInput has
// none, see the note above.
export function useBlogPostFields(): CrudFieldSpec<BlogPostInput>[] {
  return useGeneratedFields<BlogPostInput>(blogPostFields, {
    // The schema allows null for these, the form insists.
    order: ["name", "slug", "description", "richBody"],
    set: {
      description: { type: "textarea" },
      name: { required: true },
      richBody: { required: false, type: "richtext" },
      slug: { required: true },
    },
  });
}
