import type { BlogPost, BlogPostInput } from "@/client/types.gen";
import type { CrudFieldSpec } from "@/components/admin/CrudForm";
import { blogPostFields } from "@/components/admin/generated/blogPost";
import { useGeneratedFields } from "@/components/admin/useGeneratedFields";

// Blog-post form values ARE the generated request body (BlogPostInput).
// RichBody is trusted editor HTML behind the engine's richtext field (Mantine
// RTE); coverAttachmentId is deferred with the attachment wave. The related
// Courses live on their own endpoint and get a separate panel on the edit page,
// And so does the Publication State: it is changed by its own operations, never
// By the form.

export const emptyBlogPost: BlogPostInput = {
  name: "",
  slug: "",
  description: "",
  richBody: "",
  coverAttachmentId: null,
};

export function blogPostToForm(post: BlogPost): BlogPostInput {
  return {
    name: post.name ?? "",
    slug: post.slug ?? "",
    description: post.description ?? "",
    richBody: post.richBodyHtml,
    coverAttachmentId: null,
  };
}

// Fields come from the contract (generated/blogPost.ts); what is listed here is
// Only what the schema cannot say. The state is not a field: BlogPostInput has
// None, see the note above.
export function useBlogPostFields(): CrudFieldSpec<BlogPostInput>[] {
  return useGeneratedFields<BlogPostInput>(blogPostFields, {
    // The schema allows null for these, the form insists.
    order: ["name", "slug", "description", "richBody"],
    set: {
      name: { required: true },
      slug: { required: true },
      description: { type: "textarea" },
      richBody: { type: "richtext", required: false },
    },
  });
}
