import type { PropsWithChildren } from "react";
import { Grid, Text, Paper, Stack, Container, SimpleGrid, Title, Box, Center, Card } from '@mantine/core';

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type {
  LanguageCategory,
  LanguageLandingPage,
  LanguageLandingPageForLists, LanguageLandingPageQnaItem, LeadCrud
} from "@/types";

import CourseBlock from "@/components/CourseBlock";
import * as Routes from "@/routes.js";
import i18next from "i18next";
import { usePage } from "@inertiajs/react";
import { SharedProps } from "@/types";
import LeadFormBlock from "@/components/LeadFormBlock";
import MarkdownViewer from "@/components/MarkdownViewer";

type Props = PropsWithChildren & {
  categoryLandingPages: LanguageLandingPageForLists[];
  qnaItems: LanguageLandingPageQnaItem[];
  landingPages: LanguageLandingPage[];
  courseCategory: LanguageCategory;
  lead: LeadCrud;
};

export default function Show({
  courseCategory,
  qnaItems,
  categoryLandingPages,
  lead,
}: Props) {
  const { t } = useTranslation();

  const {
    auth: { user },
  } = usePage<SharedProps>().props;

  const items = [
    {
      name: t("language_categories.index.header"),
      url: Routes.language_categories_path(),
    },
    {
      name: courseCategory.header!,
      url: Routes.language_category_path(courseCategory.slug!),
    },
  ];

  return (
    <ApplicationLayout items={items} header={courseCategory.header!}>
      <Container>
        {courseCategory.description && (
          <Grid mb="xl">
            <Grid.Col span={{ base: 12, sm: 8 }}>
              <Text size="lg">{courseCategory.description}</Text>
            </Grid.Col>
          </Grid>
        )}

        <SimpleGrid cols={{ base: 2, xs: 3, sm: 4 }} mb="xl">
          {categoryLandingPages.map((lp) => (
            <CourseBlock key={lp.id} landingPage={lp} />
          ))}
        </SimpleGrid>

        {!user.guest && i18next.language === 'ru' && (
          <Grid align="center" justify="space-between" gutter={0}>
            <Grid.Col span={{ base: 12, xs: 7 }}>
              <Center>
                <Text fz={40} mb="xs" fw="bold">{t("home.index.consultation")}</Text>
              </Center>
            </Grid.Col>
            <Grid.Col span={{ base: 12, sm: 5 }}>
              <Card withBorder shadow="sm" p="xl">
                <LeadFormBlock lead={lead} />
              </Card>
            </Grid.Col>
          </Grid>
        )}

        {qnaItems.length > 0 && (
          <Stack py="xl">
            <Title order={2}>
              {t("languages.show.sort_questions")}
            </Title>
            <SimpleGrid cols={{ base: 1, xs: 2 }}>
              {qnaItems.map((item) => (
                <Box key={item.id}>
                  <Text size="lg" fw="bold">
                    {item.question}
                  </Text>
                  <MarkdownViewer>{item.answer}</MarkdownViewer>
                </Box>
              ))}
            </SimpleGrid>
          </Stack>
        )}
      </Container>
    </ApplicationLayout>
  );
}
