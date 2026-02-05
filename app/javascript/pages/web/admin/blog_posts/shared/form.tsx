import { Button, FileInput, Select, Textarea, TextInput } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { useAppForm } from '@/hooks/useAppForm';
import type { BlogPostCrud, HttpRouterMethod } from '@/types';

type Props = {
  data: BlogPostCrudWithMeta;
  url: string;
  method?: HttpRouterMethod;
};

type StateOption = {
  key: string;
  value: string;
};

type BlogPostCrudWithMeta = BlogPostCrud & {
  meta?: {
    states?: StateOption[];
  };
};

const locales = [
  { name: 'Russian', code: 'ru' },
  { name: 'English', code: 'en' },
];

export default function Form({ data, url, method }: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    getSelectProps,
    getFileInputProps,
    submit,
    formState: { isSubmitting },
  } = useAppForm<BlogPostCrud>({
    url,
    method: method ?? 'post',
    container: data,
  });
  const states = data.meta?.states ?? [];

  return (
    <form onSubmit={submit}>
      <Select {...getSelectProps('state', states, 'key', 'value')} />
      <TextInput {...getInputProps('name')} />
      <FileInput {...getFileInputProps('cover')} name="cover" />
      <TextInput {...getInputProps('slug')} />
      <Textarea {...getInputProps('description')} rows={5} />
      <Textarea {...getInputProps('body')} rows={12} />
      <Button type="submit" mt="xl" loading={isSubmitting}>
        {tHelpers(($) => $.submit.save)}
      </Button>
    </form>
  );
}
