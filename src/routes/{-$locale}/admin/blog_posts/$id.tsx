import { Button, Card, Center, Group, Loader, MultiSelect, Stack, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminGetBlogPostOptions,
  adminGetBlogPostQueryKey,
  adminListBlogPostsQueryKey,
  adminListCoursesOptions,
  adminSetBlogPostRelatedCoursesMutation,
  adminUpdateBlogPostMutation,
} from "@/client/@tanstack/react-query.gen";
import { zBlogPostInput } from "@/client/zod.gen";
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

  const backToList = () => navigate({ to: "/{-$locale}/admin/blog_posts" });

  const mutation = useResourceMutation({
    mutation: adminUpdateBlogPostMutation(),
    invalidate: [adminListBlogPostsQueryKey(), adminGetBlogPostQueryKey({ path: { id: postId } })],
    successMessage: t(($) => $.admin.crud.updated),
    errorMessage: t(($) => $.admin.crud.saveError),
    onDone: backToList,
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <Center py="xl">
          <Loader />
        </Center>
      ) : (
        <>
          <CrudForm
            key={data.id}
            fields={fields}
            schema={zBlogPostInput}
            defaultValues={blogPostToForm(data)}
            onSubmit={(values) => mutation.mutate({ path: { id: postId }, body: values })}
            submitLabel={t(($) => $.admin.crud.save)}
            onCancel={backToList}
            isPending={mutation.isPending}
          />
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

// The promoted-courses editor (legacy `related_courses` member action, now an
// explicit set). Selection order is the display order the backend persists.
function RelatedCoursesPanel({ postId, initial }: { postId: number; initial: number[] }) {
  const { t } = useTranslation();
  const [selected, setSelected] = useState(initial.map(String));

  // The course catalog is small (a few dozen), so one large page covers it.
  const courses = useQuery(adminListCoursesOptions({ query: { page: 1, perPage: 100 } }));
  const options = (courses.data?.items ?? []).map((course) => ({
    value: String(course.id),
    label: course.slug,
  }));

  const mutation = useResourceMutation({
    mutation: adminSetBlogPostRelatedCoursesMutation(),
    invalidate: [adminGetBlogPostQueryKey({ path: { id: postId } })],
    successMessage: t(($) => $.admin.crud.updated),
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
