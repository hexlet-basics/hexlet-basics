import type { Method } from "@inertiajs/core";
import {
  Button,
  Checkbox,
  Fieldset,
  FileInput,
  Select,
  Textarea,
  TextInput,
} from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData, enumToOptions } from "@/lib/utils";
import type {
  Language,
  LanguageLandingPage,
  LanguageLandingPageCreate,
  LanguageLandingPageUpdate,
} from "@/types";

type Props = {
  data: LanguageLandingPageCreate | LanguageLandingPageUpdate;
  url: string;
  method?: Method;
  languages: Language[];
  landingPages: LanguageLandingPage[];
};

export default function Form({
  data,
  landingPages,
  url,
  method,
  languages,
}: Props) {
  const { t } = useTranslation();

  const stateEventsEnum = t(
    ($) => $.models.attributes.language_landing_page["state/values"],
    { returnObjects: true },
  );
  const stateEventsSelectData = arrayToSelectData(
    enumToOptions(stateEventsEnum),
    "id",
    "name",
  );
  const languagesSelectData = arrayToSelectData(languages, "id", "slug");
  const landingPagesSelectData = arrayToSelectData(
    landingPages,
    "id",
    "header",
  );

  const { onSubmit, processing, form } = useAppForm(data, {
    url,
    method: method ?? "post",
  });

  return (
    <form onSubmit={onSubmit}>
      <Checkbox
        label={t(($) => $.models.attributes.language_landing_page.main)}
        {...form.getCheckboxProps("main")}
      />
      <Checkbox
        label={t(($) => $.models.attributes.language_landing_page.listed)}
        {...form.getCheckboxProps("listed")}
      />
      <Checkbox
        label={t(($) => $.models.attributes.language_landing_page.footer)}
        {...form.getCheckboxProps("footer")}
      />
      <Select
        label={t(($) => $.models.attributes.language_landing_page.state)}
        {...form.getSelectProps("state", stateEventsSelectData)}
      />
      <Select
        label={t(($) => $.models.attributes.language_landing_page.language_id)}
        {...form.getSelectProps("language_id", languagesSelectData)}
      />
      <Select
        label={t(
          ($) =>
            $.models.attributes.language_landing_page
              .landing_page_to_redirect_id,
        )}
        {...form.getSelectProps(
          "landing_page_to_redirect_id",
          landingPagesSelectData,
        )}
      />
      <TextInput
        label={t(($) => $.models.attributes.base.slug)}
        {...form.getInputProps("slug")}
      />
      <TextInput
        label={t(($) => $.models.attributes.language_landing_page.order)}
        {...form.getInputProps("order")}
      />
      <TextInput
        label={t(($) => $.models.attributes.language_landing_page.meta_title)}
        {...form.getInputProps("meta_title")}
      />
      <Textarea
        label={t(
          ($) => $.models.attributes.language_landing_page.meta_description,
        )}
        {...form.getInputProps("meta_description")}
        rows={3}
      />
      <TextInput
        label={t(($) => $.models.attributes.base.name)}
        {...form.getInputProps("name")}
      />
      <TextInput
        label={t(($) => $.models.attributes.base.header)}
        {...form.getInputProps("header")}
      />
      <Textarea
        label={t(($) => $.models.attributes.base.description)}
        {...form.getInputProps("description")}
        rows={5}
      />
      <TextInput
        label={t(
          ($) => $.models.attributes.language_landing_page.used_in_header,
        )}
        {...form.getInputProps("used_in_header")}
      />
      <Textarea
        label={t(
          ($) => $.models.attributes.language_landing_page.used_in_description,
        )}
        {...form.getInputProps("used_in_description")}
        rows={5}
      />
      <FileInput
        label={t(
          ($) => $.models.attributes.language_landing_page.outcomes_image,
        )}
        {...form.getFileInputProps("outcomes_image")}
      />
      <TextInput
        label={t(
          ($) => $.models.attributes.language_landing_page.outcomes_header,
        )}
        {...form.getInputProps("outcomes_header")}
      />
      <Textarea
        label={t(
          ($) => $.models.attributes.language_landing_page.outcomes_description,
        )}
        {...form.getInputProps("outcomes_description")}
        rows={5}
      />
      <Button type="submit" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
