import {
  Card,
  Center,
  Container,
  Grid,
  List,
  Stack,
  Text,
} from '@mantine/core';
import { Undo2 } from 'lucide-react';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import AppAnchor from '@/components/Elements/AppAnchor';
import LeadFormBlock from '@/components/LeadFormBlock';
import XssContent from '@/components/XssContent';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import type { LeadCrud } from '@/types';

type Props = PropsWithChildren & {
  lead: LeadCrud;
  from?: string;
};

export default function New({ lead, from }: Props) {
  // const { t: tAr } = useTranslation("activerecord");
  // const { t: tHelpers } = useTranslation("helpers");
  const { t: tViews } = useTranslation('web');
  const helpItemsList = tViews(($) => $.leads.new.help_items, {
    returnObjects: true,
  });
  const header = tViews(($) => $.leads.new.header);

  console.log(from);
  return (
    <ApplicationLayout header={header} center>
      <Container py="xl">
        <Grid gutter="xl">
          <Grid.Col span={{ base: 12, lg: 7 }} mb="xl">
            <Stack gap="md">
              <Text>{tViews(($) => $.leads.new.description)}</Text>
              <Text fw={500}>{tViews(($) => $.leads.new.how_can_we_help)}</Text>
              <List>
                {helpItemsList.map((item, index) => (
                  <List.Item key={item}>
                    <XssContent>{item}</XssContent>
                  </List.Item>
                ))}
              </List>
              <XssContent>{tViews(($) => $.leads.new.do_it)}</XssContent>
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
                    {tViews(($) => $.leads.new.return)}
                  </Text>
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
