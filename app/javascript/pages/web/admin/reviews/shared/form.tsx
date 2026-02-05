import { Button, Checkbox, Select, Textarea, TextInput } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { useAppForm } from '@/hooks/useAppForm';
import * as Routes from '@/routes.js';
import type { HttpRouterMethod, Language, ReviewCrud } from '@/types';

type StateOption = {
  key: string;
  value: string;
};

type ReviewCrudWithMeta = ReviewCrud & {
  meta?: {
    states?: StateOption[];
  };
};

type Props = {
  data: ReviewCrudWithMeta;
  url: string;
  courses: Language[];
  method?: HttpRouterMethod;
};

export default function Form({ courses, data, url, method }: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    getSelectProps,
    submit,
    formState: { isSubmitting },
  } = useAppForm<ReviewCrud>({
    url,
    method: method ?? 'post',
    container: data, // передаем контейнер целиком
  });
  const states = data.meta?.states ?? [];

  return (
    <form onSubmit={submit}>
      <Select {...getSelectProps('state', states, 'key', 'value')} />
      <Checkbox {...getInputProps('pinned')} />
      <Select {...getSelectProps('language_id', courses, 'id', 'slug')} />
      {/* Для XAutocomplete пока ставим TextInput (можно позже сделать кастомный autocomplete-хук) */}
      <TextInput {...getInputProps('user_id')} />
      <TextInput {...getInputProps('first_name')} autoComplete="name" />
      <TextInput {...getInputProps('last_name')} autoComplete="name" />
      <Textarea {...getInputProps('body')} rows={8} />
      <Button type="submit" loading={isSubmitting}>
        {tHelpers(($) => $.submit.save)}
      </Button>
    </form>
  );
}
