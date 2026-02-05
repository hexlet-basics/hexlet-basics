import { Alert, Card, Container, Grid, List, Stack, Text } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import LeadFormBlock from '@/components/LeadFormBlock';
import XssContent from '@/components/XssContent';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import type { LanguageLandingPage, LeadCrud } from '@/types';

type Props = {
  courseLandingPage: LanguageLandingPage;
  lead: LeadCrud;
};

export default function Success(props: Props) {
  const { t: tViews } = useTranslation('web');
  const { courseLandingPage, lead } = props;

  const header = tViews(($) => $.languages.success.header, {
    name: courseLandingPage.header,
  });

  return (
    <ApplicationLayout header={header} center>
      <Container>
        <Alert mb="xl">
          <XssContent>
            {tViews(($) => $.languages.success.add_review)}
          </XssContent>
        </Alert>
        <Grid gutter="xl">
          <Grid.Col span={{ base: 12, sm: 7 }} mb="xl">
            <Stack>
              <Text>{tViews(($) => $.languages.success.description)}</Text>
              <Text fw="bold">
                {tViews(($) => $.languages.success.choose_your_path)}
              </Text>
              <List>
                <List.Item>
                  <XssContent>
                    {tViews(($) => $.languages.success.changing_career_html)}
                  </XssContent>
                </List.Item>
                <List.Item>
                  <XssContent>
                    {tViews(($) => $.languages.success.getting_new_skill_html)}
                  </XssContent>
                </List.Item>
              </List>
              <Text fw="bold">
                {tViews(($) => $.languages.success.struggle_choosing)}
              </Text>
              <Text>{tViews(($) => $.languages.success.leave_request)}</Text>
            </Stack>
          </Grid.Col>
          <Grid.Col span={{ base: 12, sm: 5 }}>
            <Card withBorder p="xl">
              <LeadFormBlock autoFocus leadDto={lead} />
            </Card>
          </Grid.Col>
        </Grid>
      </Container>
    </ApplicationLayout>
  );
}
