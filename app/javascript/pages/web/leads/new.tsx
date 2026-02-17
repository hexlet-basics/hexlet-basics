import {
  Card,
  Center,
  Container,
  Grid,
  List,
  Stack,
  Text,
} from "@mantine/core";
import { IconArrowBackUp } from "@tabler/icons-react";
import type { PropsWithChildren } from "react";
import { Trans, useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
import LeadFormBlock from "@/components/LeadFormBlock";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import type { LeadCrud } from "@/types";

type Props = PropsWithChildren & {
  lead: LeadCrud;
  from?: string;
};

export default function New({ lead, from }: Props) {
  const { t } = useTranslation();

  const helpItemsList = t(($) => $.leads.new.help_items, {
    returnObjects: true,
  });
  const header = t(($) => $.leads.new.header);

  return (
    <ApplicationLayout header={header} center>
      <Container py="xl">
        <Grid gutter="xl">
          <Grid.Col span={{ base: 12, lg: 7 }} mb="xl">
            <Stack gap="md">
              <Text>{t(($) => $.leads.new.description)}</Text>
              <Text fw={500}>{t(($) => $.leads.new.how_can_we_help)}</Text>
              <List>
                {helpItemsList.map((item) => (
                  <List.Item key={item}>{item}</List.Item>
                ))}
              </List>
              <Trans
                t={t}
                i18nKey={($) => $.leads.new.do_it}
                components={{
                  b: <strong />,
                }}
              />
            </Stack>
          </Grid.Col>
          <Grid.Col span={{ base: 12, lg: 5 }}>
            <Card withBorder p="xl">
              <LeadFormBlock leadDto={lead} />
            </Card>
            {from && (
              <Center mt="xl">
                <AppAnchor href={from} c="grey">
                  <Text component="span" me="xs">
                    {t(($) => $.leads.new.return)}
                  </Text>
                  <IconArrowBackUp size={14} />
                </AppAnchor>
              </Center>
            )}
          </Grid.Col>
        </Grid>
      </Container>
    </ApplicationLayout>
  );
}
