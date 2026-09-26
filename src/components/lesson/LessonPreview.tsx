import { Divider, Stack } from "@mantine/core";
import { useTranslation } from "react-i18next";

// The learner's markup, rendered under the editor as they type it.
//
// What is rendered is whatever a learner put in the buffer, so it goes where it
// cannot reach the application around it: `srcdoc` with an empty `sandbox` gives
// the frame an opaque origin and withholds `allow-same-origin` and
// `allow-scripts` alike — nothing in it can read the page's cookies or storage,
// run code, submit a form or navigate the player away. Legacy wrote into a
// same-origin frame directly; the courses shown here are HTML and CSS, which
// need none of what the sandbox takes away.
//
// The workspace imports this lazily, so a course without a preview never
// fetches it.
export default function LessonPreview({ code }: { code: string }) {
  const { t } = useTranslation();

  return (
    <Stack gap={0}>
      <Divider color="gray.4" />
      <iframe
        title={t(($) => $.courses.lessons.show.preview.title)}
        sandbox=""
        srcDoc={code}
        width="100%"
        style={{ border: 0, display: "block" }}
      />
    </Stack>
  );
}
