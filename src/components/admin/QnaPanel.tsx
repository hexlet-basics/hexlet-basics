import {
  ActionIcon,
  Button,
  Card,
  Group,
  Stack,
  Textarea,
  TextInput,
  Title,
  Tooltip,
} from "@mantine/core";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import type { QnaItem, QnaItemInput } from "@/client/types.gen";
import { useDeleteConfirmation } from "@/hooks/useDeleteConfirmation";

interface QnaPanelProps {
  items: QnaItem[];
  onCreate: (input: QnaItemInput) => void;
  onUpdate: (id: number, input: QnaItemInput) => void;
  onDelete: (id: number) => void;
  isPending?: boolean;
}

// The nested FAQ editor shown on a category's / landing page's edit screen
// (legacy nested `qna_items` resources). Presentational: the route owns the
// parent-scoped queries and mutations, so the panel works for both parents.
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
        <NewQnaItemForm onCreate={onCreate} isPending={isPending} />
      </Stack>
    </Card>
  );
}

function QnaItemForm({
  item,
  onSave,
  onDelete,
  isPending,
}: {
  item: QnaItem;
  onSave: (input: QnaItemInput) => void;
  onDelete: () => void;
  isPending?: boolean;
}) {
  const { t } = useTranslation();
  const confirmDelete = useDeleteConfirmation();
  const [question, setQuestion] = useState(item.question);
  const [answer, setAnswer] = useState(item.answer);

  return (
    <Stack gap="xs">
      <Group align="flex-start" wrap="nowrap">
        <Stack gap="xs" style={{ flex: 1 }}>
          <TextInput
            label={t(($) => $.models.attributes.qna_item.question)}
            value={question}
            onChange={(event) => setQuestion(event.currentTarget.value)}
          />
          <Textarea
            label={t(($) => $.models.attributes.qna_item.answer)}
            autosize
            minRows={2}
            value={answer}
            onChange={(event) => setAnswer(event.currentTarget.value)}
          />
        </Stack>
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
      </Group>
      <Group justify="flex-end">
        <Button
          size="xs"
          variant="light"
          loading={isPending}
          disabled={question.trim() === "" || answer.trim() === ""}
          onClick={() => onSave({ question, answer })}
        >
          {t(($) => $.admin.crud.save)}
        </Button>
      </Group>
    </Stack>
  );
}

function NewQnaItemForm({
  onCreate,
  isPending,
}: {
  onCreate: (input: QnaItemInput) => void;
  isPending?: boolean;
}) {
  const { t } = useTranslation();
  const [question, setQuestion] = useState("");
  const [answer, setAnswer] = useState("");

  return (
    <Stack gap="xs">
      <Title order={6}>{t(($) => $.admin.qna.add)}</Title>
      <TextInput
        label={t(($) => $.models.attributes.qna_item.question)}
        value={question}
        onChange={(event) => setQuestion(event.currentTarget.value)}
      />
      <Textarea
        label={t(($) => $.models.attributes.qna_item.answer)}
        autosize
        minRows={2}
        value={answer}
        onChange={(event) => setAnswer(event.currentTarget.value)}
      />
      <Group justify="flex-end">
        <Button
          size="xs"
          loading={isPending}
          disabled={question.trim() === "" || answer.trim() === ""}
          onClick={() => {
            onCreate({ question, answer });
            setQuestion("");
            setAnswer("");
          }}
        >
          {t(($) => $.admin.crud.create)}
        </Button>
      </Group>
    </Stack>
  );
}
