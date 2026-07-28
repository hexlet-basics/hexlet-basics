import { Center, Loader, Stack, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminGetCourseCategoryOptions,
  adminGetCourseCategoryQueryKey,
  adminListCourseCategoriesQueryKey,
  adminUpdateCourseCategoryMutation,
} from "@/client/@tanstack/react-query.gen";
import { zCourseCategoryInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import {
  courseCategoryToForm,
  courseCategoryToInput,
  useCourseCategoryFields,
} from "@/components/admin/resources/courseCategory";

export const Route = createFileRoute(
  "/{-$locale}/admin/language_categories/$id",
)({
  component: EditCourseCategory,
});

function EditCourseCategory() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const fields = useCourseCategoryFields();
  const { id } = Route.useParams();
  const categoryId = Number(id);
  const [serverError, setServerError] = useState<string | null>(null);

  const { data, isLoading } = useQuery(
    adminGetCourseCategoryOptions({ path: { id: categoryId } }),
  );

  const backToList = () =>
    navigate({ to: "/{-$locale}/admin/language_categories" });

  const mutation = useMutation({
    ...adminUpdateCourseCategoryMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: adminListCourseCategoriesQueryKey(),
      });
      queryClient.invalidateQueries({
        queryKey: adminGetCourseCategoryQueryKey({ path: { id: categoryId } }),
      });
      notifications.show({
        color: "green",
        message: t(($) => $.admin.crud.updated),
      });
      backToList();
    },
    onError: () => setServerError(t(($) => $.admin.crud.saveError)),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.edit)}</Title>
      {isLoading || !data ? (
        <Center py="xl">
          <Loader />
        </Center>
      ) : (
        <CrudForm
          // Remount per row so the form re-seeds its default values from the
          // loaded category instead of keeping a stale initial snapshot.
          key={data.id}
          fields={fields}
          schema={zCourseCategoryInput}
          defaultValues={courseCategoryToForm(data)}
          onSubmit={async (values) => {
            setServerError(null);
            await mutation.mutateAsync({
              path: { id: categoryId },
              body: courseCategoryToInput(values),
            });
          }}
          submitLabel={t(($) => $.admin.crud.save)}
          onCancel={backToList}
          isPending={mutation.isPending}
          serverError={serverError}
        />
      )}
    </Stack>
  );
}
