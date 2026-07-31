import { Stack, Title } from "@mantine/core";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import {
  adminCreateBlogPostMutation,
  adminListBlogPostsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zBlogPostInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import { emptyBlogPost, useBlogPostFields } from "@/components/admin/resources/blogPost";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/blog_posts/new")({
  component: NewBlogPost,
});

export function NewBlogPost() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useBlogPostFields();

  const backToList = () => navigate({ to: "/{-$locale}/admin/blog_posts" });

  const mutation = useResourceMutation({
    mutation: adminCreateBlogPostMutation(),
    invalidate: [adminListBlogPostsQueryKey()],
    successMessage: t(($) => $.admin.crud.created),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zBlogPostInput}
        defaultValues={emptyBlogPost}
        onSubmit={(values) => mutation.mutate({ body: values })}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
      />
    </Stack>
  );
}
