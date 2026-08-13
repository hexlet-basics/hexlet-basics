import { ActionIcon, Box, Divider, Group, Stack } from "@mantine/core";
import { modals } from "@mantine/modals";
import { IconRepeat } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";

// The bar under the workspace, where the buttons that act on the exercise live.
// It carries reset today; run and the two navigation buttons join it with their
// own tickets, which is why it sits below the tabs rather than inside the editor
// pane — from here it is reachable from whichever tab the learner is on.
export default function LessonControls({ onReset }: { onReset: () => void }) {
  const { t } = useTranslation();

  // Resetting throws away whatever the learner has written, and there is no undo
  // once the buffer is overwritten — so it asks first.
  const confirmReset = () =>
    modals.openConfirmModal({
      title: t(($) => $.common.confirm),
      labels: { confirm: t(($) => $.common.boolean.yes), cancel: t(($) => $.common.boolean.no) },
      onConfirm: onReset,
    });

  return (
    <Stack gap={0}>
      <Divider />
      <Box py="sm">
        <Group justify="center">
          <ActionIcon
            variant="light"
            size="lg"
            onClick={confirmReset}
            aria-label={t(($) => $.helpers.reset)}
          >
            <IconRepeat size={18} />
          </ActionIcon>
        </Group>
      </Box>
    </Stack>
  );
}
