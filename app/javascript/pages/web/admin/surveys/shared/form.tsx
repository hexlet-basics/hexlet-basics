import {
  Box,
  Button,
  Fieldset,
  Select,
  Textarea,
  TextInput,
} from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { useAppForm } from '@/hooks/useAppForm';
import type {
  HttpRouterMethod,
  SurveyCrudWithAttrs,
  SurveyItemCrud,
} from '@/types';

type Props = {
  data: SurveyCrudWithAttrs;
  url: string;
  method?: HttpRouterMethod;
};

export default function Form({ data, url, method }: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    getSelectProps,
    submit,
    useArrayField,
    formState: { isSubmitting },
  } = useAppForm<SurveyCrudWithAttrs>({
    url,
    method: method ?? 'post',
    container: data, // передаем контейнер целиком
  });

  const itemsField = useArrayField('items_attributes');
  const defaultItem: SurveyItemCrud = {
    id: null,
    survey_id: null,
    value: '',
    tag_list: '',
    state: null,
    order: null,
    value_for_select: null,
    _destroy: false,
  };

  return (
    <form onSubmit={submit}>
      <TextInput {...getInputProps('question')} />
      <TextInput {...getInputProps('slug')} />
      <Textarea {...getInputProps('description')} rows={8} />

      <Fieldset>
        {itemsField.fields.map((field, index) => (
          <Box key={field._internalId}>
            <TextInput {...getInputProps(`items.${index}.value`)} />
            <TextInput {...getInputProps(`items.${index}.tag_list`)} />
            <Select
              {...getSelectProps(
                `items.${index}.state`,
                data.meta.item_states,
                'value',
                'key',
              )}
            />
            <TextInput {...getInputProps(`items.${index}.order`)} />
            <Button
              variant="outline"
              color="red"
              mt="xs"
              onClick={() => itemsField.remove(index)}
            >
              {tHelpers('crud.remove')}
            </Button>
          </Box>
        ))}
        <Button
          variant="light"
          mt="sm"
          onClick={() => itemsField.append(defaultItem)}
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
