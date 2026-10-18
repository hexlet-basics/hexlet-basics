import { Anchor, Card, Center, Container, Grid, List, Stack, Text, Title } from "@mantine/core";
import { IconArrowBackUp } from "@tabler/icons-react";
import { createFileRoute } from "@tanstack/react-router";
import { Trans, useTranslation } from "react-i18next";
import { z } from "zod";
import LeadFormBlock from "@/components/LeadFormBlock";
import { requireAuth } from "@/lib/auth";

// Only a same-site path is followed back: `from` comes from the URL, and an
// absolute address would turn the form into an open redirect.
const sameSitePath = z
  .string()
  .refine((path) => path.startsWith("/") && !path.startsWith("//"))
  .optional()
  .catch(undefined);

// The consultation request page, ported from legacy leads/new. Legacy required
// a session for the page as well as for the submit, so the guard sits here too.
// `from` is where the visitor came from: the page links back to it and the form
// returns there once sent.
export const Route = createFileRoute("/{-$locale}/leads/new")({
  validateSearch: z.object({ from: sameSitePath }),
  beforeLoad: ({ context, location }) => {
    requireAuth(context.user, location.href);
  },
  component: New,
});

function New() {
  const { t } = useTranslation();
  const { from } = Route.useSearch();

  const helpItems = t(($) => $.leads.new.help_items, { returnObjects: true });

  return (
    <Container py="xl">
      <Title order={1} ta="center" mb="xl">
        {t(($) => $.leads.new.header)}
      </Title>
      <Grid gap="xl">
        <Grid.Col span={{ base: 12, lg: 7 }} mb="xl">
          <Stack gap="md">
            <Text>{t(($) => $.leads.new.description)}</Text>
            <Text fw={500}>{t(($) => $.leads.new.how_can_we_help)}</Text>
            <List>
              {helpItems.map((item) => (
                <List.Item key={item}>{item}</List.Item>
              ))}
            </List>
            <Text>
              <Trans t={t} i18nKey={($) => $.leads.new.do_it} components={{ b: <strong /> }} />
            </Text>
          </Stack>
        </Grid.Col>
        <Grid.Col span={{ base: 12, lg: 5 }}>
          <Card withBorder p="xl">
            <LeadFormBlock redirectHref={from} />
          </Card>
          {from && (
            <Center mt="xl">
              <Anchor href={from} c="gray" size="sm">
                <Text component="span" me="xs">
                  {t(($) => $.leads.new.return)}
                </Text>
                <IconArrowBackUp size={14} />
              </Anchor>
            </Center>
          )}
        </Grid.Col>
      </Grid>
    </Container>
  );
}
