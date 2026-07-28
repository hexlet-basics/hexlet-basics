import { Stack, Title } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  adminCreateCourseCategoryMutation,
  adminListCourseCategoriesQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zCourseCategoryInput } from "@/client/zod.gen";
import { CrudForm } from "@/components/admin/CrudForm";
import {
  courseCategoryToInput,
  emptyCourseCategory,
  useCourseCategoryFields,
} from "@/components/admin/resources/courseCategory";

export const Route = createFileRoute(
  "/{-$locale}/admin/language_categories/new",
)({
  component: NewCourseCategory,
});

function NewCourseCategory() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const fields = useCourseCategoryFields();
  const [serverError, setServerError] = useState<string | null>(null);

  const backToList = () =>
    navigate({ to: "/{-$locale}/admin/language_categories" });

  const mutation = useMutation({
    ...adminCreateCourseCategoryMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: adminListCourseCategoriesQueryKey(),
      });
      notifications.show({
        color: "green",
        message: t(($) => $.admin.crud.created),
      });
      backToList();
    },
    onError: () => setServerError(t(($) => $.admin.crud.saveError)),
  });

  return (
    <Stack>
      <Title order={2}>{t(($) => $.admin.crud.new)}</Title>
      <CrudForm
        fields={fields}
        schema={zCourseCategoryInput}
        defaultValues={emptyCourseCategory}
        onSubmit={async (values) => {
          setServerError(null);
          await mutation.mutateAsync({ body: courseCategoryToInput(values) });
        }}
        submitLabel={t(($) => $.admin.crud.create)}
        onCancel={backToList}
        isPending={mutation.isPending}
        serverError={serverError}
      />
    </Stack>
  );
}
