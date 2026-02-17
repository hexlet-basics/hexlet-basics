import { Alert, Card, Container, Grid, List, Stack, Text } from "@mantine/core";
import { Trans, useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
import LeadFormBlock from "@/components/LeadFormBlock";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import type { LanguageLandingPage, LeadCrud } from "@/types";

type Props = {
  courseLandingPage: LanguageLandingPage;
  lead: LeadCrud;
};

export default function Success(props: Props) {
  const { t } = useTranslation();

  const { courseLandingPage, lead } = props;

  const header = t(($) => $.languages.success.header, {
    name: courseLandingPage.header,
  });

  return (
    <ApplicationLayout header={header} center>
      <Container>
        <Alert mb="xl">
          <Trans
            t={t}
            i18nKey={($) => $.languages.success.add_review}
            components={{
              a: (
                <AppAnchor
                  external
                  href="https://taplink.cc/codebasics_reviews"
                />
              ),
            }}
          />
        </Alert>
        <Grid gutter="xl">
          <Grid.Col span={{ base: 12, sm: 7 }} mb="xl">
            <Stack>
              <Text>{t(($) => $.languages.success.description)}</Text>
              <Text fw="bold">
                {t(($) => $.languages.success.choose_your_path)}
              </Text>
              <List>
                <List.Item>
                  <Trans
                    t={t}
                    i18nKey={($) => $.languages.success.changing_career_html}
                    components={{
                      a: (
                        <AppAnchor
                          external
                          href="https://ru.hexlet.io/courses_for_beginners?utm_source=code-basics&utm_medium=referral&utm_campaign=courses_for_beginners&utm_content=finished_course_page"
                        />
                      ),
                    }}
                  />
                </List.Item>
                <List.Item>
                  <Trans
                    t={t}
                    i18nKey={($) => $.languages.success.getting_new_skill_html}
                    components={{
                      a: (
                        <AppAnchor
                          external
                          href="https://ru.hexlet.io/courses_for_programmers?utm_source=code-basics&utm_medium=referral&utm_campaign=courses_for_beginners&utm_content=finished_course_page"
                        />
                      ),
                    }}
                  />
                </List.Item>
              </List>
              <Text fw="bold">
                {t(($) => $.languages.success.struggle_choosing)}
              </Text>
              <Text>{t(($) => $.languages.success.leave_request)}</Text>
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
