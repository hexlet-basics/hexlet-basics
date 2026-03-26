import { Link } from "@inertiajs/react";
import { Card, Container, Group, SimpleGrid, Stack, Text } from "@mantine/core";
import { IconArrowRight } from "@tabler/icons-react";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { LanguageCategory } from "@/types/serializers";

type Props = PropsWithChildren & {
  categories: LanguageCategory[];
};

export default function Index({ categories }: Props) {
  const { t } = useTranslation();
  const header = t(($) => $.language_categories.index.header);

  const items = [
    {
      name: header,
      url: Routes.language_categories_path(),
    },
  ];

  return (
    <ApplicationLayout items={items} header={header}>
      <Container size="lg">
        <SimpleGrid py="md" cols={{ base: 1, md: 2, lg: 3 }}>
          {categories.map((category) => (
            <Card
              component={Link}
              key={category.id}
              href={Routes.language_category_path(category.slug!)}
              padding="xl"
              radius="md"
              withBorder
              h="100%"
              td="none"
            >
              <Stack h="100%">
                <Text size="xl" fw={700}>
                  {category.header}
                </Text>
                <Group mt="auto" c="blue">
                  <Text>{t(($) => $.language_categories.index.link)}</Text>
                  <IconArrowRight size={16} />
                </Group>
              </Stack>
            </Card>
          ))}
        </SimpleGrid>
      </Container>
    </ApplicationLayout>
  );
}
