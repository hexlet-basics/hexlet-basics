import type { PropsWithChildren } from "react";
import { Container, Grid, Card, Text, List, Stack, Title, Center } from '@mantine/core';

import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { LeadCrud } from "@/types";
import LeadFormBlock from "@/components/LeadFormBlock";
import XssContent from "@/components/XssContent";
import AppAnchor from "@/components/AppAnchor";
import { Undo2 } from "lucide-react";

type Props = PropsWithChildren & {
  lead: LeadCrud
  from?: string
};

export default function New({ lead, from }: Props) {
  const { t } = useTranslation();
  // const { t: tAr } = useTranslation("activerecord");
  // const { t: tHelpers } = useTranslation("helpers");
  const { t: tViews } = useTranslation("web");
  const helpItems = tViews('leads.new.help_items', { returnObjects: true })
  const header = tViews("leads.new.header")

  console.log(from)
  return (
    <ApplicationLayout header={header} center>
      <Container py="xl">
        <Grid gutter="xl">
          <Grid.Col span={{ base: 12, lg: 7 }} mb="xl">
            <Stack gap="md">
              <Text>{tViews("leads.new.description")}</Text>
              <Text fw={500}>{tViews('leads.new.how_can_we_help')}</Text>
              <List>
                {helpItems.map((item, index) =>
                  <List.Item key={index}>
                    <XssContent>
                      {item}
                    </XssContent>
                  </List.Item>
                )}
              </List>
              <XssContent>{tViews('leads.new.do_it')}</XssContent>
            </Stack>
          </Grid.Col>
          <Grid.Col span={{ base: 12, lg: 5 }}>
            <Card withBorder p="xl">
              <LeadFormBlock lead={lead} />
            </Card>
            {from && (
              <Center mt="xl">
                <AppAnchor href={from} c="grey">
                  <Text component="span" me="xs">{tViews('leads.new.return')}</Text>
                  <Undo2 size={14} />
                </AppAnchor>
              </Center>
            )}
          </Grid.Col>
        </Grid>
      </Container>
    </ApplicationLayout>
  );
}
