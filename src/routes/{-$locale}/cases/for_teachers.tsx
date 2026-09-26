import {
  Box,
  Card,
  Container,
  Group,
  Image,
  List,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { createFileRoute, notFound } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import learningVideo from "@/assets/course-landing-page/learning_ru.mp4";
import discussionIcon from "@/assets/for-school-teachers-page/discussion-icon.svg";
import integrationIcon from "@/assets/for-school-teachers-page/integration-icon.svg";
import interactiveFormatIcon from "@/assets/for-school-teachers-page/interactive-format-icon.svg";
import practiceIcon from "@/assets/for-school-teachers-page/practice-icon.svg";
import { ButtonLink } from "@/components/RouterLink";
import { siteUrl } from "@/lib/site-url";

// The case for school and university teachers, ported from legacy
// cases/for_teachers. Legacy served it for ru only (`require_russian_locale`),
// so its one video is the ru one.
export const Route = createFileRoute("/{-$locale}/cases/for_teachers")({
  beforeLoad: ({ params }) => {
    if (params.locale !== "ru") throw notFound();
  },
  // Legacy meta: title, the description cut as Rails `truncate(160)` cut it,
  // Open Graph and a summary Twitter card. Legacy pointed the canonical at the
  // book page by mistake; it is the page's own URL here.
  head: ({ match }) => {
    const { t } = match.context.i18n;
    const title = t(($) => $.cases.for_teachers.title);
    const description = truncate(
      t(($) => $.cases.for_teachers.meta.description),
      160,
    );
    return {
      meta: [
        { title: `CodeBasics | ${title}` },
        { name: "description", content: description },
        { property: "og:title", content: title },
        { property: "og:description", content: description },
        { name: "twitter:card", content: "summary" },
        { name: "twitter:site", content: t(($) => $.links.hexlet_twitter_handle) },
      ],
      links: [{ rel: "canonical", href: siteUrl(match.pathname) }],
    };
  },
  component: ForTeachers,
});

// Rails String#truncate: at most `length` characters, the "..." omission
// included, cut at a character boundary.
function truncate(text: string, length: number): string {
  return text.length > length ? `${text.slice(0, length - 3)}...` : text;
}

// The card images, named by the `img` field of each card string.
const CARD_ICONS: Record<string, string> = {
  "discussion-icon": discussionIcon,
  "integration-icon": integrationIcon,
  "interactive-format-icon": interactiveFormatIcon,
  "practice-icon": practiceIcon,
};

function ForTeachers() {
  const { t } = useTranslation();
  const { user } = Route.useRouteContext();

  const interactiveApproachItems = t(($) => $.cases.for_teachers.interactive_approach_list, {
    returnObjects: true,
  });
  const earlyCareerGuidanceItems = t(($) => $.cases.for_teachers.early_career_guidance_list, {
    returnObjects: true,
  });
  const howToLearnCards = t(($) => $.cases.for_teachers.how_to_learn_programming_cards, {
    returnObjects: true,
  });

  return (
    <Container py="md" size="lg">
      <Title order={1} mb="sm">
        {t(($) => $.cases.for_teachers.header)}
      </Title>
      <Text size="lg" c="dimmed" mb="xl">
        {t(($) => $.cases.for_teachers.description)}
      </Text>
      <ButtonLink to="/{-$locale}" size="lg" variant="filled" px="xl" mb={50}>
        {t(($) => $.cases.for_teachers.try)}
      </ButtonLink>

      <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="xl" my="xl">
        <Title order={2} pe="xl">
          {t(($) => $.cases.for_teachers.integrate_into_education)}
        </Title>
        <Text size="lg" c="dimmed">
          {t(($) => $.cases.for_teachers.lay_programming_foundations)}
        </Text>
      </SimpleGrid>

      <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="xl" my="xl">
        <Stack>
          <Text fw="bold" fz="h4">
            {t(($) => $.cases.for_teachers.interactive_approach)}
          </Text>
          <List>
            {interactiveApproachItems.map((item) => (
              <List.Item mb="sm" key={item}>
                {item}
              </List.Item>
            ))}
          </List>
          <Text fw="bold" fz="h4">
            {t(($) => $.cases.for_teachers.early_career_guidance)}
          </Text>
          <List>
            {earlyCareerGuidanceItems.map((item) => (
              <List.Item mb="sm" key={item}>
                {item}
              </List.Item>
            ))}
          </List>
        </Stack>
        <Box>
          <video
            src={learningVideo}
            autoPlay
            loop
            muted
            playsInline
            style={{
              width: "100%",
              height: "60%",
              objectFit: "cover",
              borderRadius: "var(--mantine-radius-md)",
            }}
          />
        </Box>
      </SimpleGrid>

      <Card my="xl" radius="lg" bg="dark" c="white">
        <Title order={2} mb="xl">
          {t(($) => $.cases.for_teachers.integrate_now)}
        </Title>
        <Group justify="space-between" align="end" mt="md" gap="md" wrap="wrap">
          <Text c="gray.5">{t(($) => $.cases.for_teachers.open_browser_and_sign_up)}</Text>
          {/* A visitor is asked to sign up; a learner is sent to pick a course. */}
          {user ? (
            <ButtonLink to="/{-$locale}" variant="white" c="dark" size="lg" px="xl">
              {t(($) => $.cases.for_teachers.select_course)}
            </ButtonLink>
          ) : (
            <ButtonLink to="/{-$locale}/users/new" variant="white" c="dark" size="lg" px="xl">
              {t(($) => $.cases.for_teachers.sign_up)}
            </ButtonLink>
          )}
        </Group>
      </Card>

      <Box my="xl">
        <Title order={2} mb="xl">
          {t(($) => $.cases.for_teachers.how_to_learn_programming)}
        </Title>
        <SimpleGrid cols={{ base: 1, md: 2, xl: 4 }} spacing="lg">
          {howToLearnCards.map((item) => (
            <Card key={item.title}>
              <Image
                src={CARD_ICONS[item.img]}
                alt={item.img}
                w={90}
                h={65}
                mb="xl"
                loading="lazy"
                fit="contain"
              />
              <Text mb="sm" fw={500}>
                {item.title}
              </Text>
              <Text c="dimmed">{item.subtitle}</Text>
            </Card>
          ))}
        </SimpleGrid>
      </Box>
    </Container>
  );
}
