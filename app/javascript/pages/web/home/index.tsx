import { Head, usePage } from "@inertiajs/react";
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
import { useId } from "@mantine/hooks";
import {
  IconBook2,
  IconBrandGithub,
  IconRobot,
  IconSend,
} from "@tabler/icons-react";
import i18next from "i18next";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import type { FAQPage, Question, WithContext } from "schema-dts";
import { Banner } from "@/components/Banner";
import BlogPostBlock from "@/components/BlogPostBlock";
import CourseBlock from "@/components/CourseBlock";
import LeadFormBlock from "@/components/LeadFormBlock";
import MarkdownViewer from "@/components/MarkdownViewer";
import SignUpFormBlock from "@/components/SignUpFormBlock";
import XssContent from "@/components/XssContent";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import { hasObjectKey } from "@/lib/utils";
import { getResourceUrl } from "@/resources";
import * as Routes from "@/routes.js";
import type {
  BlogPost,
  LanguageMember,
  LeadCrud,
  Review,
  UserSignUpForm,
} from "@/types";

type Props = PropsWithChildren & {
  blogPosts: BlogPost[];
  courseMembersByCourseId: LanguageMember[];
  reviews: Review[];
  newUser: UserSignUpForm;
  lead: LeadCrud;
};

// TODO: move to locales
type ReviewItem = {
  name: string;
  avatar: string;
  body: string;
};

const reviews: Record<string, ReviewItem[]> = {
  ru: [
    {
      name: "Александр Авдошкин",
      avatar: getResourceUrl("avdoshkin.jpg"),
      body: `Если бы не коронавирус, выполнил бы всё в заход (в смысле каждый день по несколько пунктов в теме).
Изучаю с нуля, ваш портал очень ориентирован на новичков. Спасибо вам большое!`,
    },
    {
      name: "Сергей Тюрин",
      avatar: getResourceUrl("tyrin.jpg"),
      body: `Очень всё доступно даже для полного профана вроде меня. Эта вводная по JS вошла в мой туговатый ум,
складно как недостающий пазл. Всем кидаю линк на эту страничку.`,
    },
    {
      name: "Элиях Клейман",
      avatar: getResourceUrl("user-avatar.png"),
      body: `Для меня это первый курс для новичка. Понравилось тем, что вся информация структурирована
и дана по мере изучения материала в иерархичном порядке, что значительно повышает и желание к обучению`,
    },
  ],
  en: [
    {
      name: "Aleksandr Avdoshkin",
      avatar: getResourceUrl("avdoshkin.jpg"),
      body: "As someone with zero coding skills, I'd say that CodeBasics is focused on newcomers. Thank you very much!",
    },
    {
      name: "Sergei Tyurin",
      avatar: getResourceUrl("tyrin.jpg"),
      body: "This is all very approachable even for a dummy like me. Now I show people this platform when I get the chance.",
    },
    {
      name: "Eliyah Kleyman",
      avatar: getResourceUrl("user-avatar.png"),
      body: `For me, it was my very first programming course. I liked it because all the information is very well
structured and given in a clear hierarchical order. It motivated me a lot to move forward in my studies.`,
    },
  ],
};

export default function Index({
  blogPosts,
  newUser,
  lead,
  courseMembersByCourseId,
}: Props) {
  const { t } = useTranslation();

  const baseId = useId("courses");

  const {
    locale,
    landingPagesForLists,
    auth: { user },
  } = usePage().props;

  const faq = t(($) => $.faq.main, {
    returnObjects: true,
  });

  const entities: Question[] = Object.values(faq).map((item) => {
    return {
      "@type": "Question",
      name: item.question,
      acceptedAnswer: {
        "@type": "Answer",
        text: item.answer,
      },
    };
  });
  // https://developers.google.com/search/docs/appearance/structured-data/faqpage
  const qaSchema: WithContext<FAQPage> = {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    mainEntity: entities,
  };

  return (
    <ApplicationLayout>
      <Head>
        {Object.keys(faq).length > 0 && (
          <script type="application/ld+json">{JSON.stringify(qaSchema)}</script>
        )}
      </Head>
      {/*<Banner />*/}
      <Container
        ta="center"
        size="sm"
        my={{
          base: "xl",
          xs: 60,
          sm: 70,
          md: 80,
          lg: 120,
        }}
      >
        <Badge
          component="a"
          target="_blank"
          c="gray"
          style={{ cursor: "pointer" }}
          href="https://github.com/orgs/hexlet-basics/repositories"
          leftSection={<IconBrandGithub size={10} />}
          // autoContrast
          size="sm"
          variant="default"
        >
          {t(($) => $.home.index.hero.source_code)}
        </Badge>

        <Title
          mb="xs"
          lh={1}
          order={1}
          fz={{
            base: 25,
            xs: 30,
            sm: 40,
            md: 50,
          }}
        >
          {t(($) => $.home.index.hero.free_programming_courses)}
        </Title>

        <XssContent c="gray" fz="xl" mb="xl">
          {t(($) => $.home.index.hero.fastest_way_to_start_coding)}
        </XssContent>

        <Center>
          <Button component="a" href={`#${baseId}`} size="lg" me="sm">
            {t(($) => $.home.index.hero.try)}
          </Button>
          {user.guest && (
            <Button
              component="a"
              href={Routes.new_user_path()}
              size="lg"
              variant="outline"
            >
              {t(($) => $.home.index.sign_up)}
            </Button>
          )}
        </Center>
      </Container>
      <Container
        size="lg"
        my={{
          base: "xl",
          xs: 60,
          sm: 70,
          md: 80,
          lg: 100,
        }}
      >
        <SimpleGrid cols={{ base: 1, xs: 2, md: 3 }} my="xl">
          <Card bg="blue.0" p="xl" c="blue.9">
            <IconBook2 />
            <Text fw="bold" fz="h2">
              {t(($) => $.home.index.hero.courses_count)}
            </Text>
            {t(($) => $.home.index.hero.courses_count_description)}
          </Card>
          <Card bg="green.0" p="xl" c="green.9">
            <IconSend />
            <Text fw="bold" fz="h2">
              {t(($) => $.home.index.hero.community_count)}
            </Text>
            {t(($) => $.home.index.hero.community_count_description)}
          </Card>
          <Card bg="violet.0" p="xl" c="violet.9">
            <IconRobot />
            <Text fw="bold" fz="h2">
              {t(($) => $.home.index.hero.ai_count)}
            </Text>
            {t(($) => $.home.index.hero.ai_count_description)}
          </Card>
        </SimpleGrid>
      </Container>
      <Container size="lg" my="xl">
        <Title id={baseId} order={2} mb="xl">
          {t(($) => $.home.languages.courses)}
        </Title>

        <SimpleGrid spacing="md" cols={{ base: 2, xs: 3, md: 4 }}>
          {landingPagesForLists.map((lp) => (
            <CourseBlock
              lazy
              landingPage={lp}
              key={lp.id}
              courseMember={courseMembersByCourseId[lp.id]}
            />
          ))}
        </SimpleGrid>
      </Container>
      <Container size="lg" my="xl">
        <Box mb="xl">
          <Title order={2}>{t(($) => $.home.index.reviews)}</Title>
          <Divider />
        </Box>

        <SimpleGrid spacing="md" cols={{ base: 1, xs: 2, md: 3 }}>
          {reviews[hasObjectKey(reviews, locale) ? locale : "ru"].map(
            (review) => (
              <Box key={review.name}>
                <Group mb="lg">
                  <Image
                    src={review.avatar}
                    fit="contain"
                    radius="xl"
                    loading="lazy"
                    w={40}
                    // width={50}
                    // height={50}
                    alt={`Аватар пользователя ${review.name}`}
                  />
                  <Text fw="bold">{review.name}</Text>
                </Group>

                <Text fs="italic">{review.body}</Text>
              </Box>
            ),
          )}
        </SimpleGrid>
      </Container>
      {blogPosts.length > 0 && (
        <Container size="lg">
          <Box mb="xl">
            <Title order={2}>{t(($) => $.home.index.blog_posts)}</Title>
            <Divider />
          </Box>
          <SimpleGrid spacing="md" cols={{ base: 1, xs: 2, md: 3 }}>
            {blogPosts.map((post) => (
              <BlogPostBlock key={post.id} post={post} lazy />
            ))}
          </SimpleGrid>
        </Container>
      )}
      <Container size="lg" my="xl">
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
      {user.guest && (
        <Container size="lg" mt={100}>
          <Grid align="center" justify="space-between" gutter={0}>
            <Grid.Col span={{ base: 12, md: 6 }} mb="lg">
              <Center>
                <Text
                  fz={{
                    base: 25,
                    xs: 30,
                    sm: 40,
                    md: 50,
                  }}
                  fw="bold"
                >
                  {t(($) => $.home.index.join)}
                </Text>
              </Center>
            </Grid.Col>
            <Grid.Col span={{ base: 12, md: 5 }}>
              <Card withBorder shadow="sm" p="xl">
                <SignUpFormBlock userDto={newUser} />
              </Card>
            </Grid.Col>
          </Grid>
        </Container>
      )}
      {!user.guest && i18next.language === "ru" && (
        <Container size="lg" mt={100}>
          <Grid align="center" justify="space-between" gutter={0}>
            <Grid.Col span={{ base: 12, md: 6 }}>
              <Center>
                <Text fz={40} fw="bold">
                  {t(($) => $.home.index.consultation)}
                </Text>
              </Center>
            </Grid.Col>
            <Grid.Col span={{ base: 12, md: 4 }}>
              <Card withBorder shadow="sm" p="xl">
                <LeadFormBlock leadDto={lead} />
              </Card>
            </Grid.Col>
          </Grid>
        </Container>
      )}
    </ApplicationLayout>
  );
}
