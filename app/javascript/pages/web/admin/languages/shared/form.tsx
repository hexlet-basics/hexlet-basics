import { Button, FileInput, Select, TextInput } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { useAppForm } from '@/hooks/useAppForm';
import { enumToOptions } from '@/lib/utils';
import type { HttpRouterMethod } from '@/types';
import type LanguageCrud from '@/types/serializers/LanguageCrud';

type Props = {
  data: LanguageCrud;
  url: string;
  method?: HttpRouterMethod;
};

// const locales = [
//   { name: 'Russian', code: 'ru' },
//   { name: 'English', code: 'en' },
// ];

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

  const {
    getInputProps,
    getSelectProps,
    getFileInputProps,
    submit,
    formState: { isSubmitting },
  } = useAppForm<LanguageCrud>({
    url,
    method: method ?? 'post',
    container: data,
  });

  return (
    <form onSubmit={submit}>
      <Select
        {...getSelectProps(
          'progress',
          languageProgressEnumOptions,
          'id',
          'name',
        )}
      />
      <Select
        {...getSelectProps(
          'learn_as',
          languageLearnAsEnumOptions,
          'id',
          'name',
        )}
      />
      <TextInput {...getInputProps('slug')} />
      <TextInput {...getInputProps('hexlet_program_landing_page')} />
      <TextInput {...getInputProps('openai_assistant_id')} />
      <FileInput {...getFileInputProps('cover')} name="cover" />

      <Button type="submit" loading={isSubmitting}>
        {tHelpers('submit.save')}
      </Button>
    </form>
  );
}
