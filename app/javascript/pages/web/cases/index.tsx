import {
  Anchor,
  Card,
  Container,
  Group,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { IconArrowRight } from "@tabler/icons-react";
import { useTranslation } from "react-i18next";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";

export default function Index() {
  const { t } = useTranslation();

  return (
    <ApplicationLayout>
      <Container py="md" mih="100%">
        <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} py="md">
          <Card bg="gray.0" p="xl" radius="xl" shadow="sm" h="100%">
            <Stack h="100%">
              <Title order={4} mb="xs" fw="bold">
                {t(($) => $.cases.index.for_teachers)}
              </Title>
              <Anchor
                href={Routes.for_teachers_cases_path()}
                td="none"
                mt="auto"
              >
                <Group gap={6}>
                  <Text>{t(($) => $.cases.index.link)}</Text>
                  <IconArrowRight size={16} />
                </Group>
              </Anchor>
            </Stack>
          </Card>
        </SimpleGrid>
      </Container>
    </ApplicationLayout>
  );
}
