import type { Method } from "@inertiajs/core";
import {
  Box,
  Button,
  Checkbox,
  Fieldset,
  Select,
  TextInput,
} from "@mantine/core";
import { useTranslation } from "react-i18next";
import { eventNames } from "@/generated/event_names";
import { useAppForm } from "@/hooks/useAppForm";
import { arrayToSelectData } from "@/lib/utils";
import type {
  Survey,
  SurveyItemCrud,
  SurveyScenarioCrud,
  SurveyScenarioItemCrud,
  SurveyScenarioTriggerCrud,
} from "@/types";

type Props = {
  data: SurveyScenarioCrud;
  surveys: Survey[];
  surveysItems: SurveyItemCrud[];
  url: string;
  method?: Method;
};

export default function Form({
  data,
  url,
  method,
  surveys,
  surveysItems,
}: Props) {
  const { t } = useTranslation();

  const payload = data;
  const surveysItemsSelectData = arrayToSelectData(
    surveysItems,
    "id",
    "value_for_select",
  );
  const surveysSelectData = arrayToSelectData(surveys, "id", "question");

  const { onSubmit, processing, form } = useAppForm(payload, {
    url,
    method: method ?? "post",
  });

  const triggersCollection = form.useCollection<SurveyScenarioTriggerCrud>(
    "triggers_attributes",
  );
  const itemsCollection =
    form.useCollection<SurveyScenarioItemCrud>("items_attributes");
  const defaultTrigger: SurveyScenarioTriggerCrud = {
    id: null,
    event_name: null,
    event_threshold_count: null,
    _destroy: false,
    meta: { model: "", relations: {} },
  };
  const defaultItem: SurveyScenarioItemCrud = {
    id: null,
    survey_id: null,
    // survey: null,
    order: null,
    _destroy: false,
    meta: { model: "", relations: {} },
  };

  const eventNameOptions = eventNames.map((event) => ({
    label: event,
    value: event,
  }));

  return (
    <form onSubmit={onSubmit}>
      <TextInput {...form.getInputProps("name")} />
      <Select
        {...form.getSelectProps("survey_item_id", surveysItemsSelectData)}
      />
      <Fieldset>
        {triggersCollection.forms.map((triggerForm) => (
          <Box key={`${triggerForm.index}-${triggerForm.data.id ?? "new"}`}>
            <Select
              {...triggerForm.getSelectProps("event_name", eventNameOptions)}
            />
            <TextInput
              {...triggerForm.getInputProps("event_threshold_count")}
              type="number"
            />
            <Checkbox {...triggerForm.getCheckboxProps("_destroy")} />
            {!triggerForm.data.id && (
              <Button
                type="button"
                variant="outline"
                color="red"
                mt="xs"
                onClick={() => triggersCollection.remove(triggerForm.index)}
              >
                {t(($) => $.helpers.crud.remove)}
              </Button>
            )}
          </Box>
        ))}
        <Button
          type="button"
          variant="light"
          mt="sm"
          onClick={() => triggersCollection.add(defaultTrigger)}
        >
          {t(($) => $.helpers.crud.add)}
        </Button>
      </Fieldset>
      <Fieldset>
        {itemsCollection.forms.map((itemForm) => (
          <Box key={`${itemForm.index}-${itemForm.data.id ?? "new"}`}>
            <Select
              {...itemForm.getSelectProps("survey_id", surveysSelectData)}
            />
            <TextInput {...itemForm.getInputProps("order")} type="number" />
            <Checkbox {...itemForm.getCheckboxProps("_destroy")} />
            {!itemForm.data.id && (
              <Button
                type="button"
                variant="outline"
                color="red"
                mt="xs"
                onClick={() => itemsCollection.remove(itemForm.index)}
              >
                {t(($) => $.helpers.crud.remove)}
              </Button>
            )}
          </Box>
        ))}
        <Button
          type="button"
          variant="light"
          mt="sm"
          onClick={() => itemsCollection.add(defaultItem)}
        >
          {t(($) => $.helpers.crud.add)}
        </Button>
      </Fieldset>
      <Button type="submit" loading={processing}>
        {t(($) => $.helpers.submit.save)}
      </Button>
    </form>
  );
}
