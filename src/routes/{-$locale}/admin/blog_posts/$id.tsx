import { Button, Card, Group, MultiSelect, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminGetBlogPostOptions,
  adminGetBlogPostQueryKey,
  adminListBlogPostsQueryKey,
  adminListCourseOptionsOptions,
  adminSetBlogPostRelatedCoursesMutation,
  adminSuggestBlogPostRelatedCoursesMutation,
  adminUpdateBlogPostMutation,
} from "@/client/@tanstack/react-query.gen";
import { zBlogPostInput } from "@/client/zod.gen";
import { BlogPostStatePanel } from "@/components/admin/BlogPostStatePanel";
import FormSkeleton from "@/components/admin/FormSkeleton";
import { CrudForm } from "@/components/admin/CrudForm";
import { blogPostToForm, useBlogPostFields } from "@/components/admin/resources/blogPost";
import { useResourceMutation } from "@/hooks/useResourceMutation";

export const Route = createFileRoute("/{-$locale}/admin/blog_posts/$id")({
  component: EditBlogPost,
});

function EditBlogPost() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const fields = useBlogPostFields();
  const { id } = Route.useParams();
  const postId = Number(id);

  const { data, isLoading } = useQuery(adminGetBlogPostOptions({ path: { id: postId } }));

  const backToList = () => {
    void navigate({ to: "/{-$locale}/admin/blog_posts" });
  };

  const mutation = useResourceMutation({
    mutation: adminUpdateBlogPostMutation(),
    fieldErrors: true,
    invalidate: [adminListBlogPostsQueryKey(), adminGetBlogPostQueryKey({ path: { id: postId } })],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <FormSkeleton />
      ) : (
        <>
          <CrudForm
            key={data.id}
            fields={fields}
            schema={zBlogPostInput}
            defaultValues={blogPostToForm(data)}
            onSubmit={(values) => mutation.mutateAsync({ path: { id: postId }, body: values })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
          />
          <BlogPostStatePanel post={data} />
          <RelatedCoursesPanel
            key={`related-${data.id}-${data.relatedCourseIds.join(",")}`}
            postId={postId}
            initial={data.relatedCourseIds}
          />
        </>
      )}
    </Stack>
  );
}

// The promoted-courses editor: an explicit set, plus the legacy AI pick
// (`related_courses` member action). Selection order is the display order the
// Backend persists. The pick runs in a background job, so its result shows up
// On a later load of the post — the same manual reload legacy needed.
function RelatedCoursesPanel({ postId, initial }: { postId: number; initial: number[] }) {
  const { t } = useTranslation();
  const [selected, setSelected] = useState(initial.map(String));

  // Every course, from the unpaginated options list: a page would silently drop
  // The courses past its end from the picker.
  const courses = useQuery(adminListCourseOptionsOptions());
  const options = (courses.data ?? []).map((course) => ({
    value: String(course.id),
    label: course.name,
  }));

  const mutation = useResourceMutation({
    mutation: adminSetBlogPostRelatedCoursesMutation(),
    invalidate: [adminGetBlogPostQueryKey({ path: { id: postId } })],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
  });

  const suggestMutation = useResourceMutation({
    mutation: adminSuggestBlogPostRelatedCoursesMutation(),
    invalidate: [adminGetBlogPostQueryKey({ path: { id: postId } })],
    successMessage: t(($) => $.admin.blogPosts.suggestQueued),
    errorMessage: t(($) => $.admin.crud.saveError),
  });

  return (
    <Card withBorder p="xl" maw={720}>
      <Title order={4} mb="md">
        {t(($) => $.admin.blogPosts.relatedCourses)}
      </Title>
      <MultiSelect
        data={options}
        value={selected}
        onChange={setSelected}
        searchable
        placeholder={t(($) => $.admin.crud.search)}
      />
      <Group justify="flex-end" mt="md">
        <Button
          variant="light"
          color="grape"
          loading={suggestMutation.isPending}
          onClick={() => suggestMutation.mutate({ path: { id: postId } })}
        >
          {t(($) => $.admin.blogPosts.suggestRelatedCourses)}
        </Button>
        <Button
          loading={mutation.isPending}
          onClick={() =>
            mutation.mutate({
              path: { id: postId },
              body: { courseIds: selected.map(Number) },
            })
          }
        >
          {t(($) => $.admin.crud.save)}
        </Button>
      </Group>
    </Card>
  );
}
