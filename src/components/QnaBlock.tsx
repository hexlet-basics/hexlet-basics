import { Box, SimpleGrid, Stack, Text, Title } from "@mantine/core";
import { useTranslation } from "react-i18next";
import type { QnaItem } from "@/client/types.gen";
import MarkdownViewer from "@/components/MarkdownViewer";

// The "Sorting out the questions" block legacy closed its category and course
// pages with: each question in bold over its markdown answer, two columns from
// the xs breakpoint. Nothing is rendered for a page with no questions.
export default function QnaBlock({ items }: { items: QnaItem[] }) {
  const { t } = useTranslation();
  if (items.length === 0) return null;

  return (
    <Stack py="xl">
      <Title order={2}>{t(($) => $.courses.show.sort_questions)}</Title>
      <SimpleGrid cols={{ base: 1, xs: 2 }}>
        {items.map((item) => (
          <Box key={item.id}>
            <Text size="lg" fw="bold">
              {item.question}
            </Text>
            <MarkdownViewer>{item.answer}</MarkdownViewer>
          </Box>
        ))}
      </SimpleGrid>
    </Stack>
  );
}
