import { Button } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import type { HTTPVerb } from 'use-inertia-form';
import {
  XCheck,
  XDynamicInputs,
  XForm,
  XInput,
  XSelect,
} from '@/components/forms';
import { eventNames } from '@/event_names';
import type { Survey, SurveyItemCrud, SurveyScenarioCrud } from '@/types';

type Props = {
  data: SurveyScenarioCrud;
  surveys: Survey[];
  surveysItems: SurveyItemCrud[];
  url: string;
  method?: HTTPVerb;
};

export default function Form({
  data,
  url,
  method,
  surveys,
  surveysItems,
}: Props) {
  const { t } = useTranslation();
  const { t: tHelpers } = useTranslation('helpers');

  const eventNameOptions: Record<string, string>[] = eventNames.map(
    (event) => ({
      label: event,
      value: event,
    }),
  );

  return (
    <XForm method={method} model="survey_scenario" data={data} to={url}>
      <XInput field="name" />
      <XSelect
        field="survey_item_id"
        items={surveysItems}
        valueField="id"
        labelField="value_for_select"
      />

      <XDynamicInputs
        model="triggers"
        label="Triggers"
        emptyData={{ event_name: null, event_threshold_count: null }}
      >
        <XSelect
          field="event_name"
          items={eventNameOptions}
          valueField="value"
          labelField="label"
        />
        <XInput field="event_threshold_count" type="number" />
        <XCheck field="_destroy" />
      </XDynamicInputs>

      <XDynamicInputs
        model="items"
        label="Items"
        emptyData={{ survey_id: null }}
      >
        <XSelect
          field="survey_id"
          valueField="id"
          labelField="question"
          items={surveys}
        />
        <XInput field="order" type="number" />
        <XCheck field="_destroy" />
      </XDynamicInputs>

      <Button type="submit">{tHelpers('submit.save')}</Button>
    </XForm>
  );
}
