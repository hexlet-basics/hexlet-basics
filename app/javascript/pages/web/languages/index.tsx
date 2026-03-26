import { Link } from "@inertiajs/react";
import {
  Box,
  Card,
  Container,
  Divider,
  Group,
  Image,
  List,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useTranslation } from "react-i18next";
import CourseBlock from "@/components/CourseBlock";
import { HoverLift } from "@/components/HoverLift";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import { propsForExternalLink } from "@/lib/utils";
import {
  reviewShowcaseAvatars,
  reviewShowcaseOrder,
} from "@/pages/web/shared/reviews";
import * as Routes from "@/routes.js";
import type { LanguageCategory, LanguageLandingPageForLists } from "@/types";

type Props = {
  catalogLandingPages: LanguageLandingPageForLists[];
  categories: LanguageCategory[];
};

export default function Index({ catalogLandingPages, categories }: Props) {
  const { t, i18n } = useTranslation();
  const isRuLocale = i18n.language === "ru";
  const reviewTexts = t(($) => $.shared.reviews_showcase, {
    returnObjects: true,
  });

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

        <Box mb="xxl">
          <Box mb="xl">
            <Title order={2}>
              {t(($) => $.pages.languages.index.categories)}
            </Title>
            <Divider />
          </Box>
          <SimpleGrid cols={{ base: 1, md: 3 }} spacing="md">
            {categories.map((category) => {
              const href = Routes.language_category_path(category.slug!);

              return (
                <Card
                  component={Link}
                  key={category.id}
                  href={href}
                  withBorder
                  radius="md"
                  p="lg"
                  td="none"
                >
                  <HoverLift h="100%">
                    <Stack gap="xs">
                      <Text fw="bold">{category.header ?? category.name}</Text>
                      {category.description && (
                        <Text c="dimmed">{category.description}</Text>
                      )}
                    </Stack>
                  </HoverLift>
                </Card>
              );
            })}
          </SimpleGrid>
        </Box>

        <Box mb="xxl">
          <Box mb="xl">
            <Title order={2}>{t(($) => $.pages.languages.index.reviews)}</Title>
            <Divider />
          </Box>
          <SimpleGrid spacing="md" cols={{ base: 1, xs: 2, md: 3 }}>
            {reviewShowcaseOrder.map((reviewId) => {
              const review = reviewTexts[reviewId];

              return (
                <Box key={review.name}>
                  <Group mb="lg">
                    <Image
                      src={reviewShowcaseAvatars[reviewId]}
                      fit="contain"
                      radius="xl"
                      loading="lazy"
                      w={40}
                      alt={review.name}
                    />
                    <Text fw="bold">{review.name}</Text>
                  </Group>
                  <Text fs="italic">{review.body}</Text>
                </Box>
              );
            })}
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
