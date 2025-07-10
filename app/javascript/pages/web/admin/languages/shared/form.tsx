import { Button } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import type { HTTPVerb } from 'use-inertia-form';
import { XFile, XForm, XInput, XSelect } from '@/components/forms';
import { enumToOptions } from '@/lib/utils';
import type LanguageCrud from '@/types/serializers/LanguageCrud';

type Props = {
  data: LanguageCrud;
  url: string;
  method?: HTTPVerb;
};

const locales = [
  { name: 'Russian', code: 'ru' },
  { name: 'English', code: 'en' },
];

export default function Form({ data, url, method }: Props) {
  const { t: tHelpers } = useTranslation('helpers');
  const { t: tEnums } = useTranslation('enumerize');
  const languageProgressEnum = tEnums('language.progress', {
    returnObjects: true,
  });
  const languageProgressEnumOptions = enumToOptions(languageProgressEnum);

  const languageLearnAsEnum = tEnums('language.learn_as', {
    returnObjects: true,
  });
  const languageLearnAsEnumOptions = enumToOptions(languageLearnAsEnum);

  return (
    <XForm method={method} model="language" data={data} to={url}>
      <XSelect
        field="progress"
        labelField="name"
        valueField="id"
        items={languageProgressEnumOptions}
      />
      <XSelect
        field="learn_as"
        labelField="name"
        valueField="id"
        items={languageLearnAsEnumOptions}
      />
      <XInput field="slug" />
      <XInput field="hexlet_program_landing_page" />
      <XInput field="openai_assistant_id" />
      <XFile metaName="cover_thumb_url" field="cover" />

      <Button type="submit">{tHelpers('submit.save')}</Button>
    </XForm>
  );
}
