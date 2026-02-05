import {
  Box,
  Button,
  Checkbox,
  Fieldset,
  Select,
  TextInput,
} from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { eventNames } from '@/generated/event_names';
import { useAppForm } from '@/hooks/useAppForm';
import type {
  HttpRouterMethod,
  Survey,
  SurveyItemCrud,
  SurveyScenarioCrud,
  SurveyScenarioItemCrud,
  SurveyScenarioTriggerCrud,
} from '@/types';

type Props = {
  data: SurveyScenarioCrud;
  surveys: Survey[];
  surveysItems: SurveyItemCrud[];
  url: string;
  method?: HttpRouterMethod;
};

export default function Form({
  data,
  url,
  method,
  surveys,
  surveysItems,
}: Props) {
  const { t: tHelpers } = useTranslation('helpers');

  const {
    getInputProps,
    getSelectProps,
    submit,
    useArrayField,
    formState: { isSubmitting },
  } = useAppForm<SurveyScenarioCrud>({
    url,
    method: method ?? 'post',
    container: data,
  });

  const triggersField = useArrayField('triggers_attributes');
  const itemsField = useArrayField('items_attributes');
  const defaultTrigger: SurveyScenarioTriggerCrud = {
    id: null,
    event_name: null,
    event_threshold_count: null,
    _destroy: false,
  };
  const defaultItem: SurveyScenarioItemCrud = {
    id: null,
    survey_id: null,
    // survey: null,
    order: null,
    _destroy: false,
  };

  const eventNameOptions = eventNames.map((event) => ({
    label: event,
    value: event,
  }));

  return (
    <form onSubmit={submit}>
      <TextInput {...getInputProps('name')} />
      <Select
        {...getSelectProps(
          'survey_item_id',
          surveysItems,
          'id',
          'value_for_select',
        )}
      />
      <Fieldset>
        {triggersField.fields.map((field, index) => (
          <Box key={field._internalId}>
            <Select
              {...getSelectProps(
                `triggers_attributes.${index}.event_name`,
                eventNameOptions,
                'value',
                'label',
              )}
            />
            <TextInput
              {...getInputProps(
                `triggers_attributes.${index}.event_threshold_count`,
              )}
              type="number"
            />
            <Checkbox
              {...getInputProps(`triggers_attributes.${index}._destroy`)}
            />
            {!field.id && (
              <Button
                variant="outline"
                color="red"
                mt="xs"
                onClick={() => triggersField.remove(index)}
              >
                {tHelpers(($) => $.crud.remove)}
              </Button>
            )}
          </Box>
        ))}
        <Button
          variant="light"
          mt="sm"
          onClick={() => triggersField.append(defaultTrigger)}
        >
          {tHelpers(($) => $.crud.add)}
        </Button>
      </Fieldset>
      <Fieldset>
        {itemsField.fields.map((field, index) => (
          <Box key={field._internalId}>
            <Select
              {...getSelectProps(
                `items_attributes.${index}.survey_id`,
                surveys,
                'id',
                'question',
              )}
            />
            <TextInput
              {...getInputProps(`items_attributes.${index}.order`)}
              type="number"
            />
            <Checkbox
              {...getInputProps(`items_attributes.${index}._destroy`)}
            />
            {!field.id && (
              <Button
                variant="outline"
                color="red"
                mt="xs"
                onClick={() => itemsField.remove(index)}
              >
                {tHelpers(($) => $.crud.remove)}
              </Button>
            )}
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
