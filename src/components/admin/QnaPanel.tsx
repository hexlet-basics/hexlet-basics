import {
  ActionIcon,
  Button,
  Card,
  Group,
  Stack,
  TextInput,
  Textarea,
  Title,
  Tooltip,
} from "@mantine/core";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import type { QnaItem, QnaItemInput } from "@/client/types.gen";
import { zQnaItemInput } from "@/client/zod.gen";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";
import { submitForm, useAppForm } from "@/lib/form";

interface QnaPanelProps {
  items: QnaItem[];
  onCreate: (input: QnaItemInput) => void;
  onUpdate: (id: number, input: QnaItemInput) => void;
  onDelete: (id: number) => void;
  isPending?: boolean;
}

// The nested FAQ editor shown on a category's / landing page's edit screen
// (legacy nested `qna_items` resources). Presentational: the route owns the
// Parent-scoped queries and mutations, so the panel works for both parents.
export function QnaPanel({ items, onCreate, onUpdate, onDelete, isPending }: QnaPanelProps) {
  const { t } = useTranslation();

  return (
    <Card withBorder p="xl" maw={720}>
      <Title order={4} mb="md">
        {t(($) => $.admin.qna.title)}
      </Title>
      <Stack>
        {items.map((item) => (
          <QnaItemForm
            key={item.id}
            item={item}
            onSave={(input) => onUpdate(item.id, input)}
            onDelete={() => onDelete(item.id)}
            isPending={isPending}
          />
        ))}
        <QnaItemForm onSave={onCreate} isPending={isPending} />
      </Stack>
    </Card>
  );
}

// One question and its answer: an existing item (saved, or deleted) when `item`
// is given, otherwise the blank form that adds one. The generated
// `zQnaItemInput` is the rule for both, so an empty field cannot be submitted.
function QnaItemForm({
  item,
  onSave,
  onDelete,
  isPending,
}: {
  item?: QnaItem;
  onSave: (input: QnaItemInput) => void;
  onDelete?: () => void;
  isPending?: boolean;
}) {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();
  const form = useAppForm({
    defaultValues: { question: item?.question ?? "", answer: item?.answer ?? "" },
    validators: { onSubmit: zQnaItemInput },
    onSubmit: ({ value, formApi }) => {
      onSave(value);
      if (!item) formApi.reset();
    },
  });

  return (
    <Stack component="form" gap="xs" onSubmit={submitForm(form)}>
      {!item && <Title order={6}>{t(($) => $.admin.qna.add)}</Title>}
      <Group align="flex-start" wrap="nowrap">
        <Stack gap="xs" flex={1}>
          <form.AppField name="question">
            {(field) => <field.TextField label={t(($) => $.models.attributes.qna_item.question)} />}
          </form.AppField>
          <form.AppField name="answer">
            {(field) => (
              <field.TextareaField
                minRows={2}
                label={t(($) => $.models.attributes.qna_item.answer)}
              />
            )}
          </form.AppField>
        </Stack>
        {item && onDelete && (
          <Tooltip label={t(($) => $.admin.crud.delete)}>
            <ActionIcon
              color="red"
              variant="light"
              mt={28}
              aria-label={t(($) => $.admin.crud.delete)}
              onClick={() => confirmDelete({ description: item.question, onConfirm: onDelete })}
            >
              ✕
            </ActionIcon>
          </Tooltip>
        )}
      </Group>
      <Group justify="flex-end">
        <Button type="submit" size="xs" variant={item ? "light" : "filled"} loading={isPending}>
          {item ? t(($) => $.admin.crud.save) : t(($) => $.admin.crud.create)}
        </Button>
      </Group>
    </Stack>
  );
}
