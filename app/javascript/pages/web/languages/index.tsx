import {
  Box,
  Card,
  Container,
  List,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useTranslation } from "react-i18next";
import CourseBlock from "@/components/CourseBlock";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import { propsForExternalLink } from "@/lib/utils";
import type { LanguageLandingPageForLists } from "@/types";

type Props = {
  catalogLandingPages: LanguageLandingPageForLists[];
};

export default function Index({ catalogLandingPages }: Props) {
  const { t, i18n } = useTranslation();
  const isRuLocale = i18n.language === "ru";

  return (
    <ApplicationLayout header={t(($) => $.pages.languages.index.header)}>
      <Container size="lg" my="xl">
        <Box mb="xxl">
          <SimpleGrid spacing="md" cols={{ base: 2, xs: 3, md: 4 }}>
            {catalogLandingPages.map((landingPage) => (
              <CourseBlock
                key={landingPage.id}
                lazy
                landingPage={landingPage}
              />
            ))}
          </SimpleGrid>
        </Box>

        {isRuLocale && (
          <Card withBorder p="xl">
            <Stack gap="md">
              <Title order={2}>
                {t(($) => $.pages.languages.index.hexlet.title)}
              </Title>
              <Text>
                {t(($) => $.pages.languages.index.hexlet.description)}
              </Text>
              <Box>
                <Text fw="bold" mb="sm">
                  {t(($) => $.pages.languages.index.hexlet.programs)}
                </Text>
                <List spacing="xs">
                  <List.Item>
                    <a
                      href="https://ru.hexlet.io/programs/python"
                      {...propsForExternalLink()}
                    >
                      {t(($) => $.pages.languages.index.hexlet.python)}
                    </a>
                  </List.Item>
                  <List.Item>
                    <a
                      href="https://ru.hexlet.io/programs/frontend"
                      {...propsForExternalLink()}
                    >
                      {t(($) => $.pages.languages.index.hexlet.frontend)}
                    </a>
                  </List.Item>
                  <List.Item>
                    <a
                      href="https://ru.hexlet.io/programs/java"
                      {...propsForExternalLink()}
                    >
                      {t(($) => $.pages.languages.index.hexlet.java)}
                    </a>
                  </List.Item>
                  <List.Item>
                    <a
                      href="https://ru.hexlet.io/programs/go"
                      {...propsForExternalLink()}
                    >
                      {t(($) => $.pages.languages.index.hexlet.go)}
                    </a>
                  </List.Item>
                  <List.Item>
                    <a
                      href="https://ru.hexlet.io/programs/devops-engineer-from-scratch"
                      {...propsForExternalLink()}
                    >
                      {t(($) => $.pages.languages.index.hexlet.devops)}
                    </a>
                  </List.Item>
                  <List.Item>
                    <a
                      href="https://ru.hexlet.io/programs/data-analytics"
                      {...propsForExternalLink()}
                    >
                      {t(($) => $.pages.languages.index.hexlet.analytics)}
                    </a>
                  </List.Item>
                </List>
              </Box>
            </Stack>
          </Card>
        )}
      </Container>
    </ApplicationLayout>
  );
}
