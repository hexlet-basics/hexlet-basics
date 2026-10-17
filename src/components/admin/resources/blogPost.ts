import type { BlogPost, BlogPostInput, BlogPostState } from "@/client/types.gen";
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
  name: "",
  slug: "",
  description: "",
  richBody: "",
  coverAttachmentId: null,
};

// The transition operations and the states each may start from: a mirror of the
// server's table (internal/handlers/blog_posts.go), used only to offer the legal
// buttons. The server decides and refuses the rest with a 409.
export type BlogPostTransition = "publish" | "unpublish" | "archive";

const transitionSources: Record<BlogPostTransition, BlogPostState[]> = {
  publish: ["draft", "archived"],
  unpublish: ["published"],
  archive: ["draft", "published"],
};

export function transitionsFrom(state: BlogPostState): BlogPostTransition[] {
  return (Object.keys(transitionSources) as BlogPostTransition[]).filter((name) =>
    transitionSources[name].includes(state),
  );
}

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
// only what the schema cannot say. The state is not a field: BlogPostInput has
// none, see the note above.
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
