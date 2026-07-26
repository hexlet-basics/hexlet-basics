import {
  Anchor,
  Box,
  Card,
  Center,
  Container,
  List,
  Loader,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { listCoursesOptions } from "@/client/@tanstack/react-query.gen";
import CourseBlock from "@/components/CourseBlock";

// Course catalog, rendered at both `/` and `/languages` (legacy URL kept for
// compat). Data is prefetched in each route's SSR loader, so `useQuery` resolves
// from the dehydrated cache on first paint — no loading flash.
export default function CourseCatalog() {
  const { t, i18n } = useTranslation();
  const isRuLocale = i18n.language === "ru";
  const { data, isPending, isError } = useQuery(listCoursesOptions());

  if (isPending)
    return (
      <Center h="100vh">
        <Loader />
      </Center>
    );

  if (isError || !data)
    return (
      <Center h="100vh">
        <Text c="red">Failed to load courses</Text>
      </Center>
    );

  return (
    <Container size="lg" my="xl">
      <Title order={1} mb="xl">
        {t(($) => $.pages.languages.index.header)}
      </Title>

      <Box mb={48}>
        <SimpleGrid spacing="md" cols={{ base: 2, xs: 3, md: 4 }}>
          {data.map((item) => (
            <CourseBlock key={item.id} lazy item={item} />
          ))}
        </SimpleGrid>
      </Box>

      {isRuLocale && (
        <Card withBorder p="xl">
          <Stack gap="md">
            <Title order={2}>
              {t(($) => $.pages.languages.index.hexlet.title)}
            </Title>
            <Text>{t(($) => $.pages.languages.index.hexlet.description)}</Text>
            <Box>
              <Text fw="bold" mb="sm">
                {t(($) => $.pages.languages.index.hexlet.programs)}
              </Text>
              <List spacing="xs">
                <List.Item>
                  <Anchor
                    href="https://ru.hexlet.io/programs/python"
                    target="_blank"
                    rel="noreferrer"
                  >
                    {t(($) => $.pages.languages.index.hexlet.python)}
                  </Anchor>
                </List.Item>
                <List.Item>
                  <Anchor
                    href="https://ru.hexlet.io/programs/frontend"
                    target="_blank"
                    rel="noreferrer"
                  >
                    {t(($) => $.pages.languages.index.hexlet.frontend)}
                  </Anchor>
                </List.Item>
                <List.Item>
                  <Anchor
                    href="https://ru.hexlet.io/programs/java"
                    target="_blank"
                    rel="noreferrer"
                  >
                    {t(($) => $.pages.languages.index.hexlet.java)}
                  </Anchor>
                </List.Item>
                <List.Item>
                  <Anchor
                    href="https://ru.hexlet.io/programs/go"
                    target="_blank"
                    rel="noreferrer"
                  >
                    {t(($) => $.pages.languages.index.hexlet.go)}
                  </Anchor>
                </List.Item>
                <List.Item>
                  <Anchor
                    href="https://ru.hexlet.io/programs/devops-engineer-from-scratch"
                    target="_blank"
                    rel="noreferrer"
                  >
                    {t(($) => $.pages.languages.index.hexlet.devops)}
                  </Anchor>
                </List.Item>
                <List.Item>
                  <Anchor
                    href="https://ru.hexlet.io/programs/data-analytics"
                    target="_blank"
                    rel="noreferrer"
                  >
                    {t(($) => $.pages.languages.index.hexlet.analytics)}
                  </Anchor>
                </List.Item>
              </List>
            </Box>
          </Stack>
        </Card>
      )}
    </Container>
  );
}
