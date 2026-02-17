import { CodeHighlight } from "@mantine/code-highlight";
import { Alert } from "@mantine/core";
import { IconAlertTriangle, IconCheck } from "@tabler/icons-react";
import { Trans, useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
import { useLessonStore } from "../store.tsx";

export default function OutputTab() {
  // const { lessonMember } = usePage().props;
  const result = useLessonStore((state) => state.result);
  const processState = useLessonStore((state) => state.processState);
  const output = useLessonStore((state) => state.output);
  const passed = useLessonStore((state) => state.passed);
  const { t } = useTranslation();

  if (processState !== "checked") {
    return null;
  }

  return (
    <>
      <Alert
        color={passed ? "green" : "yellow"}
        icon={
          passed ? <IconCheck size={16} /> : <IconAlertTriangle size={16} />
        }
        withCloseButton={false}
        fz="sm"
        variant="light"
        radius="sm"
        my="xs"
      >
        <Trans
          t={t}
          i18nKey={($) => $.common.check[result!].message}
          components={{
            a: (
              <AppAnchor external href="https://help.hexlet.io/article/20623" />
            ),
          }}
        />
      </Alert>
      <CodeHighlight code={output} language="ansi" withCopyButton={false} />
    </>
  );
}
