import {
  Box,
  Button,
  Fieldset,
  Select,
  Textarea,
  TextInput,
} from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { enums } from '@/generated/enums';
import { useAppForm } from '@/hooks/useAppForm';
import { enumToSelectData } from '@/lib/utils';
import type { HttpRouterMethod, SurveyCrud, SurveyItemCrud } from '@/types';

type Props = {
  data: SurveyCrud;
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
  } = useAppForm<SurveyCrud>({
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
            <TextInput {...getInputProps(`items_attributes.${index}.value`)} />
            <TextInput
              {...getInputProps(`items_attributes.${index}.tag_list`)}
            />
            <Select
              {...getSelectProps(
                `items_attributes.${index}.state`,
                enumToSelectData(enums.surveyItemState),
                'value',
                'label',
              )}
            />
            <TextInput {...getInputProps(`items_attributes.${index}.order`)} />
            <Button
              variant="outline"
              color="red"
              mt="xs"
              onClick={() => itemsField.remove(index)}
            >
              {tHelpers(($) => $.crud.remove)}
            </Button>
          </Box>
        ))}
        <Button
          variant="light"
          mt="sm"
          onClick={() => itemsField.append(defaultItem)}
        >
          {tHelpers(($) => $.crud.add)}
        </Button>
      </Fieldset>
      <Button type="submit" loading={isSubmitting}>
        {tHelpers(($) => $.submit.save)}
      </Button>
    </form>
  );
}
