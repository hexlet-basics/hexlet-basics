import { Alert, Anchor, Card, Container, Grid, List, Stack, Text, Title } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { isAxiosError } from "axios";
import { Trans, useTranslation } from "react-i18next";
import { getCourseOptions } from "@/client/@tanstack/react-query.gen";
import LeadFormBlock from "@/components/LeadFormBlock";
import { requireAuth } from "@/lib/auth";

// Outbound links legacy hardcoded in the page rather than in the locale files.
const reviewsHref = "https://taplink.cc/codebasics_reviews";
const careerHref =
  "https://ru.hexlet.io/courses_for_beginners?utm_source=code-basics&utm_medium=referral&utm_campaign=courses_for_beginners&utm_content=finished_course_page";
const skillsHref =
  "https://ru.hexlet.io/courses_for_programmers?utm_source=code-basics&utm_medium=referral&utm_campaign=courses_for_beginners&utm_content=finished_course_page";

// The course completion page, ported from legacy languages#success: where the
// lesson player sends a learner who finished the Course.
//
// Legacy required a session (a guest goes to sign in and comes back), then
// bounced anyone whose Enrollment is not finished to the course page with a
// warning. The finished signal is the Enrollment's state, not `progress`: the
// progress payload is computed for guests too, from the signed cookie, so it
// could say "finished" for someone who has no Enrollment at all. A learner with
// no Enrollment is bounced the same way; legacy crashed on that case.
export const Route = createFileRoute("/{-$locale}/languages/$slug/success")({
  beforeLoad: ({ context, location }) => {
    requireAuth(context.user, location.href);
  },
  loader: async ({ context, params }) => {
    try {
      const view = await context.queryClient.ensureQueryData(
        getCourseOptions({ path: { slug: params.slug } }),
      );
      if (view.enrollment?.state !== "finished") {
        throw redirect({
          to: "/{-$locale}/languages/$slug",
          params,
          search: { unfinished: true },
        });
      }
    } catch (error) {
      // Legacy sent an unknown course to the catalog.
      if (isAxiosError(error) && error.response?.status === 404) {
        throw redirect({ to: "/{-$locale}/languages", params: { locale: params.locale } });
      }
      throw error;
    }
  },
  component: Success,
});

function Success() {
  const { t } = useTranslation();
  const { slug } = Route.useParams();
  const { data } = useQuery(getCourseOptions({ path: { slug } }));

  const external = { target: "_blank", rel: "noopener noreferrer" } as const;

  return (
    <Container my="xl">
      <Title order={1} ta="center" mb="xl">
        {t(($) => $.courses.success.header, {
          name: data?.landingPage?.header ?? data?.course.name ?? "",
        })}
      </Title>
      <Alert mb="xl">
        <Trans
          t={t}
          i18nKey={($) => $.courses.success.add_review}
          components={{ a: <Anchor href={reviewsHref} {...external} /> }}
        />
      </Alert>
      <Grid gap="xl">
        <Grid.Col span={{ base: 12, sm: 7 }} mb="xl">
          <Stack>
            <Text>{t(($) => $.courses.success.description)}</Text>
            <Text fw="bold">{t(($) => $.courses.success.choose_your_path)}</Text>
            <List>
              <List.Item>
                <Trans
                  t={t}
                  i18nKey={($) => $.courses.success.changing_career_html}
                  components={{ a: <Anchor href={careerHref} {...external} /> }}
                />
              </List.Item>
              <List.Item>
                <Trans
                  t={t}
                  i18nKey={($) => $.courses.success.getting_new_skill_html}
                  components={{ a: <Anchor href={skillsHref} {...external} /> }}
                />
              </List.Item>
            </List>
            <Text fw="bold">{t(($) => $.courses.success.struggle_choosing)}</Text>
            <Text>{t(($) => $.courses.success.leave_request)}</Text>
          </Stack>
        </Grid.Col>
        <Grid.Col span={{ base: 12, sm: 5 }}>
          <Card withBorder p="xl">
            {/* Legacy put the learner straight into the form. */}
            {/* oxlint-disable-next-line jsx-a11y/no-autofocus */}
            <LeadFormBlock autoFocus />
          </Card>
        </Grid.Col>
      </Grid>
    </Container>
  );
}
