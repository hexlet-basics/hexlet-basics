import {
  Box,
  Button,
  Checkbox,
  Fieldset,
  FileInput,
  Select,
  Textarea,
  TextInput,
} from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { useAppForm } from '@/hooks/useAppForm';
import type {
  HttpRouterMethod,
  Language,
  LanguageLandingPage,
  LanguageLandingPageCrud,
  LanguageLandingPageQnaItemCrud,
} from '@/types';

type Props = {
  data: LanguageLandingPageCrud;
  url: string;
  method?: HttpRouterMethod;
  languages: Language[];
  landingPages: LanguageLandingPage[];
};

// const locales = [
//   { name: "Russian", code: "ru" },
//   { name: "English", code: "en" },
// ];

export default function Form({
  data,
  landingPages,
  url,
  method,
  languages,
}: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    getFileInputProps,
    getSelectProps,
    submit,
    useArrayField,
    formState: { isSubmitting },
  } = useAppForm<LanguageLandingPageCrud>({
    url,
    method: method ?? 'post',
    container: data,
  });

  const qnaField = useArrayField('qna_items_attributes');
  const defaultQna: LanguageLandingPageQnaItemCrud = {
    id: null,
    question: '',
    answer: '',
    _destroy: false,
  };

  return (
    <form onSubmit={submit}>
      <Checkbox {...getInputProps('main')} />
      <Checkbox {...getInputProps('listed')} />
      <Checkbox {...getInputProps('footer')} />

      <Select
        {...getSelectProps('state', data.meta.state_events, 'key', 'value')}
      />
      <Select {...getSelectProps('language_id', languages, 'id', 'slug')} />
      <Select
        {...getSelectProps(
          'landing_page_to_redirect_id',
          landingPages,
          'id',
          'header',
        )}
      />
      <TextInput {...getInputProps('slug')} />
      <TextInput {...getInputProps('order')} />
      <TextInput {...getInputProps('meta_title')} />
      <Textarea {...getInputProps('meta_description')} rows={3} />
      <TextInput {...getInputProps('name')} />
      <TextInput {...getInputProps('header')} />
      <Textarea {...getInputProps('description')} rows={5} />

      <TextInput {...getInputProps('used_in_header')} />
      <Textarea {...getInputProps('used_in_description')} rows={5} />

      <FileInput {...getFileInputProps('outcomes_image')} />
      <TextInput {...getInputProps('outcomes_header')} />
      <Textarea {...getInputProps('outcomes_description')} rows={5} />

      <Fieldset>
        {qnaField.fields.map((field, index) => (
          <Box key={field._internalId}>
            <input
              type="hidden"
              {...getInputProps(`qna_items_attributes.${index}.id`)}
            />
            <TextInput
              {...getInputProps(`qna_items_attributes.${index}.question`)}
            />
            <Textarea
              {...getInputProps(`qna_items_attributes.${index}.answer`)}
              rows={5}
            />
            <Checkbox
              {...getInputProps(`qna_items_attributes.${index}._destroy`)}
            />
            {!field.id && (
              <Button
                variant="outline"
                color="red"
                mt="xs"
                onClick={() => qnaField.remove(index)}
              >
                {tHelpers('crud.remove')}
              </Button>
            )}
          </Box>
        ))}
        <Button
          variant="light"
          mt="sm"
          onClick={() => qnaField.append(defaultQna)}
        >
          {tHelpers('crud.add')}
        </Button>
      </Fieldset>

      <Button type="submit" loading={isSubmitting}>
        {tHelpers('submit.save')}
      </Button>
    </form>
  );
}
