import {
  Accordion,
  Badge,
  Box,
  Button,
  Card,
  Center,
  Container,
  Divider,
  Grid,
  Group,
  Image,
  SimpleGrid,
  Text,
  Title,
} from "@mantine/core";
import { IconBook2, IconBrandGithub, IconRobot, IconSend } from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import type { ReactElement, ReactNode } from "react";
import { useTranslation } from "react-i18next";
import {
  getCurrentUserOptions,
  listBlogPostsOptions,
  listCoursesOptions,
} from "@/client/@tanstack/react-query.gen";
import avdoshkinAvatar from "@/assets/avdoshkin.jpg";
import tyurinAvatar from "@/assets/tyrin.jpg";
import kleymanAvatar from "@/assets/user-avatar.png";
import { Banner } from "@/components/Banner";
import BlogPostBlock from "@/components/BlogPostBlock";
import CourseBlock from "@/components/CourseBlock";
import MarkdownViewer from "@/components/MarkdownViewer";
import { ButtonLink } from "@/components/RouterLink";
import SignUpFormBlock from "@/components/SignUpFormBlock";

const COURSES_ANCHOR = "courses";

// Reviews of the showcase, in legacy's order, with the avatar each one carries.
const SHOWCASE = [
  { avatar: avdoshkinAvatar, id: "avdoshkin" },
  { avatar: tyurinAvatar, id: "tyurin" },
  { avatar: kleymanAvatar, id: "kleyman" },
] as const;

// A figure tile of the hero row. The tint is Mantine's `light` variant colour
// variable, which resolves per colour scheme, so it needs no scheme check in
// JS and cannot differ between the server and the first client render.
function Tile({
  color,
  icon,
  title,
  children,
}: {
  color: "blue" | "green" | "violet";
  icon: ReactNode;
  title: string;
  children: ReactNode;
}): ReactElement {
  return (
    <Card p="xl" radius="md" bg={`var(--mantine-color-${color}-light)`}>
      {icon}
      <Text fw="bold" fz="h2">
        {title}
      </Text>
      {children}
    </Card>
  );
}

// The home page, ported from legacy home#index: hero, the three figures, the
// course grid, the review showcase, the latest blog posts, the FAQ, and then
// the registration form for a guest. Legacy's consultation block for a signed-in
// ru learner is gone, here and on the other pages that had it.
export default function HomePage(): ReactElement {
  const { t } = useTranslation();
  const { data: courses } = useQuery(listCoursesOptions());
  const { data: blog } = useQuery(listBlogPostsOptions({ query: { perPage: 3 } }));
  const { data: me } = useQuery(getCurrentUserOptions());
  const user = me?.user ?? null;

  const faq = t(($) => $.faq.main, { returnObjects: true });
  const reviews = t(($) => $.shared.reviews_showcase, { returnObjects: true });
  const posts = blog?.items ?? [];

  return (
    <>
      <Banner />

      <Container ta="center" size="sm" my={{ base: "xl", sm: "xxl" }}>
        <Badge
          component="a"
          c="gray"
          href="https://github.com/orgs/hexlet-basics/repositories"
          leftSection={<IconBrandGithub size={10} />}
          size="sm"
          variant="default"
          target="_blank"
          rel="noopener noreferrer"
        >
          {t(($) => $.home.index.hero.source_code)}
        </Badge>

        <Title mb="xs" lh={1} order={1} fz={{ base: 25, md: 50, sm: 40, xs: 30 }}>
          {t(($) => $.home.index.hero.free_programming_courses)}
        </Title>

        <Text c="gray" fz="xl" mb="xl">
          {t(($) => $.home.index.hero.fastest_way_to_start_coding)}
        </Text>

        <Center>
          <Button component="a" href={`#${COURSES_ANCHOR}`} size="lg" me="sm">
            {t(($) => $.home.index.hero.try)}
          </Button>
          {!user && (
            <ButtonLink to="/{-$locale}/users/new" params size="lg" variant="outline">
              {t(($) => $.home.index.sign_up)}
            </ButtonLink>
          )}
        </Center>
      </Container>

      <Container size="lg" my={{ base: "xl", sm: "xxl" }}>
        <SimpleGrid cols={{ base: 1, md: 3, xs: 2 }} my="xl">
          <Tile color="blue" icon={<IconBook2 />} title={t(($) => $.home.index.hero.courses_count)}>
            {t(($) => $.home.index.hero.courses_count_description)}
          </Tile>
          <Tile
            color="green"
            icon={<IconSend />}
            title={t(($) => $.home.index.hero.community_count)}
          >
            {t(($) => $.home.index.hero.community_count_description)}
          </Tile>
          <Tile color="violet" icon={<IconRobot />} title={t(($) => $.home.index.hero.ai_count)}>
            {t(($) => $.home.index.hero.ai_count_description)}
          </Tile>
        </SimpleGrid>
      </Container>

      <Container size="lg" my={{ base: "xl", sm: "xxl" }}>
        <Title id={COURSES_ANCHOR} order={2} mb="xl">
          {t(($) => $.home.courses.courses)}
        </Title>

        <SimpleGrid cols={{ base: 2, md: 4, xs: 3 }}>
          {(courses ?? []).map((item) => (
            <CourseBlock key={item.id} lazy item={item} />
          ))}
        </SimpleGrid>
      </Container>

      <Container size="lg" my={{ base: "xl", sm: "xxl" }}>
        <Box mb="xl">
          <Title order={2}>{t(($) => $.home.index.reviews)}</Title>
          <Divider />
        </Box>

        <SimpleGrid cols={{ base: 1, md: 3, xs: 2 }}>
          {SHOWCASE.map(({ id, avatar }) => {
            const review = reviews[id];
            return (
              <Box key={id}>
                <Group mb="lg">
                  <Image
                    src={avatar}
                    fit="contain"
                    radius="xl"
                    loading="lazy"
                    w={40}
                    alt={review.name}
                  />
                  <Text fw="bold">{review.name}</Text>
                </Group>
                <Text fs="italic">{review.body}</Text>
              </Box>
            );
          })}
        </SimpleGrid>
      </Container>

      {posts.length > 0 && (
        <Container size="lg">
          <Box mb="xl">
            <Title order={2}>{t(($) => $.home.index.blog_posts)}</Title>
            <Divider />
          </Box>
          <SimpleGrid cols={{ base: 1, md: 3, xs: 2 }}>
            {posts.map((post) => (
              <BlogPostBlock key={post.id} post={post} lazy />
            ))}
          </SimpleGrid>
        </Container>
      )}

      <Container size="lg" my={{ base: "xl", sm: "xxl" }}>
        <Box mb="xl">
          <Title order={2}>{t(($) => $.faq.header)}</Title>
          <Divider />
        </Box>
        <Accordion defaultValue={Object.keys(faq)[0]}>
          {Object.entries(faq).map(([key, item]) => (
            <Accordion.Item key={key} value={key}>
              <Accordion.Control>
                <Box fw="bold">{item.question}</Box>
              </Accordion.Control>
              <Accordion.Panel>
                <MarkdownViewer>{item.answer}</MarkdownViewer>
              </Accordion.Panel>
            </Accordion.Item>
          ))}
        </Accordion>
      </Container>

      {!user && (
        <Container size="lg" mt={100}>
          <Grid align="center" justify="space-between" gap={0}>
            <Grid.Col span={{ base: 12, md: 6 }} mb="lg">
              <Center>
                <Text fz={{ base: 25, md: 50, sm: 40, xs: 30 }} fw="bold">
                  {t(($) => $.home.index.join)}
                </Text>
              </Center>
            </Grid.Col>
            <Grid.Col span={{ base: 12, md: 5 }}>
              <Card withBorder shadow="sm" p="xl">
                <SignUpFormBlock />
              </Card>
            </Grid.Col>
          </Grid>
        </Container>
      )}
    </>
  );
}
