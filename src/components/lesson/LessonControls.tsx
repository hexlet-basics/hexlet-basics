import { ActionIcon, Box, Button, Divider, Group, Stack } from "@mantine/core";
import { modals } from "@mantine/modals";
import { IconPlayerPlay, IconRepeat } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";

// The bar under the workspace, where the buttons that act on the exercise live.
// It sits below the tabs rather than inside the editor pane so that running a
// solution is one press away from whichever tab the learner is on.
//
// Previous and next join it with #775.
export default function LessonControls({
  onReset,
  onRun,
  running,
}: {
  onReset: () => void;
  onRun: () => void;
  running: boolean;
}) {
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

          {/* While a check is in flight the button says so and refuses a second
              press: one solution is running, and submitting it twice would tell
              the learner nothing new. */}
          <Button
            leftSection={<IconPlayerPlay size={18} />}
            onClick={onRun}
            loading={running}
            disabled={running}
          >
            {t(($) => $.courses.lessons.show.controls.run)}
          </Button>
        </Group>
      </Box>
    </Stack>
  );
}
