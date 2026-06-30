import type { Method } from "@inertiajs/core";
import { Button, Textarea, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import type { StaffRoleCrud } from "@/types";

type Props = {
  data: StaffRoleCrud;
  url: string;
  method?: Method;
};

export default function Form({ data, url, method }: Props) {
  const { t } = useTranslation();

  const { onSubmit, processing, form } = useAppForm(data, {
    url,
    method: method ?? "post",
  });

  return (
    <form onSubmit={onSubmit}>
      <TextInput
        label={t(($) => $.admin.management.roles.form.name)}
        {...form.getInputProps("name")}
      />
      <Textarea
        label={t(($) => $.admin.management.roles.form.description)}
        {...form.getInputProps("description")}
        rows={3}
      />
      <Button type="submit" loading={processing} mt="md">
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
