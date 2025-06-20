import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { TypeAnimation } from "react-type-animation";

import type {
  BlogPost,
  LanguageMember,
  Review,
  User,
} from "@/types/serializers";

import BlogPostBlock from "@/components/BlogPostBlock";
import CourseBlock from "@/components/CourseBlock";
import SignUpFormBlock from "@/components/SignUpFormBlock";
import XssContent from "@/components/XssContent";
import { getImageUrl } from "@/images";
import codeImagePathEn from "@/images/code-basics-coding-en.png";
import codeImagePathRu from "@/images/code-basics-coding-ru.png";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { LeadCrud, SharedProps } from "@/types";
import { Head, usePage } from "@inertiajs/react";
import MarkdownViewer from "@/components/MarkdownViewer";
import type { Question, FAQPage, WithContext } from "schema-dts";
import LeadFormBlock from "@/components/LeadFormBlock";
import i18next from "i18next";
import { Container, Image, Grid, Accordion, Paper, Button, Title, Text, Box, Anchor, Group, SimpleGrid, Center, Card, Divider } from '@mantine/core';

type Props = PropsWithChildren & {
  blogPosts: BlogPost[];
  courseMembersByCourseId: LanguageMember[];
  reviews: Review[];
  newUser: User;
  lead: LeadCrud;
};

const sequence = [
  "TypeScript",
  1000,
  "Java",
  1000,
  "Python",
  1000,
  "PHP",
  1000,
  "Ruby",
  1000,
  "HTML",
  1000,
  "CSS",
  1000,
  "Go",
  1000,
];

const codeImagePaths = {
  ru: codeImagePathRu,
  en: codeImagePathEn,
};

// TODO: move to locales
const reviews = {
  ru: [
    {
      name: "Александр Авдошкин",
      avatar: getImageUrl("avdoshkin.jpg"),
      body: `Если бы не коронавирус, выполнил бы всё в заход (в смысле каждый день по несколько пунктов в теме).
Изучаю с нуля, ваш портал очень ориентирован на новичков. Спасибо вам большое!`,
    },
    {
      name: "Сергей Тюрин",
      avatar: getImageUrl("tyrin.jpg"),
      body: `Очень всё доступно даже для полного профана вроде меня. Эта вводная по JS вошла в мой туговатый ум,
складно как недостающий пазл. Всем кидаю линк на эту страничку.`,
    },
    {
      name: "Элиях Клейман",
      avatar: getImageUrl("user-avatar.png"),
      body: `Для меня это первый курс для новичка. Понравилось тем, что вся информация структурирована
и дана по мере изучения материала в иерархичном порядке, что значительно повышает и желание к обучению`,
    },
  ],
  en: [
    {
      name: "Aleksandr Avdoshkin",
      avatar: getImageUrl("avdoshkin.jpg"),
      body: "As someone with zero coding skills, I'd say that CodeBasics is focused on newcomers. Thank you very much!",
    },
    {
      name: "Sergei Tyurin",
      avatar: getImageUrl("tyrin.jpg"),
      body: "This is all very approachable even for a dummy like me. Now I show people this platform when I get the chance.",
    },
    {
      name: "Eliyah Kleyman",
      avatar: getImageUrl("user-avatar.png"),
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
  const { t: tFaq } = useTranslation("faq");

  const {
    locale,
    landingPagesForLists,
    auth: { user },
  } = usePage<SharedProps>().props;

  const faq = tFaq("main", { returnObjects: true });

  const entities: Question[] = Object.values(faq).map((item) => {
    return {
      "@type": "Question",
      "name": item.question,
      "acceptedAnswer": {
        "@type": "Answer",
        "text": item.answer,
      }
    }
  })
  // https://developers.google.com/search/docs/appearance/structured-data/faqpage
  const qaSchema: WithContext<FAQPage> = {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    "mainEntity": entities,
  };

  return (
    <ApplicationLayout>
      <Head>
        {Object.keys(faq).length > 0 && <script type="application/ld+json">{JSON.stringify(qaSchema)}</script>}
      </Head>

      <Container size="lg">

        <Paper
          shadow="sm"
          my={{ base: "md", xs: "lg", sm: 60 }}
          pt={{ base: "md", xs: "lg", sm: 60 }}
          withBorder
          ps={{ base: "md", xs: "lg", sm: 60 }}
        >
          <Grid overflow="hidden" gutter={0}>
            <Grid.Col span={{ base: 12, md: 7 }}>
              <Title order={1} size="h5" fw="normal" c="dimmed">
                {t("home.hero.free_programming_courses")}
              </Title>
              <Box my="xs" visibleFrom="xs">
                <Text size="55px" fw="bold">
                  {t("home.hero.learn")}{" "}
                  <Text c="blue" component="span" fw="bold">
                    <TypeAnimation
                      // className="text-primary"
                      sequence={sequence}
                      wrapper="span"
                      speed={5}
                      repeat={Number.POSITIVE_INFINITY}
                    />
                  </Text>
                </Text>
              </Box>
              <Box my="xs" hiddenFrom="xs">
                <Text fz="h1" fw="bold">
                  {t("home.hero.learn_xs")}
                </Text>
              </Box>
              <XssContent mb="xl">
                {t("home.hero.fastest_way_to_start_coding")}
              </XssContent>
              <Button mb="lg" size="xl" component="a" href="#courses">
                {t("home.hero.try")}
              </Button>
            </Grid.Col>

            <Grid.Col
              span={{ base: 12, md: 5 }}
            >

              <Image
                alt="Code Basics learning preview"
                fetchPriority="high"
                w="auto"
                width="720"
                height="356"
                src={codeImagePaths[locale]}
              />

            </Grid.Col>
          </Grid>
        </Paper>
      </Container>


      <Container size="lg" my="xl">
        <Anchor
          id="courses"
          href="#courses"
        >
          <Title order={2} mb="xl">
            {t("home.languages.courses")}
          </Title>
        </Anchor>

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
          <Title order={2}>
            {t("home.index.reviews")}
          </Title>
          <Divider />
        </Box>

        <SimpleGrid spacing="md" cols={{ base: 1, xs: 2, md: 3 }}>
          {reviews[locale].map((review) => (
            <Box key={review.name}>
              <Group mb="lg">
                <Image
                  src={review.avatar}
                  fit="contain"
                  radius="100%"
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
          ))}
        </SimpleGrid>
      </Container>


      {blogPosts.length > 0 && (
        <Container size="lg">
          <Box mb="xl">
            <Title order={2}>
              {t("home.index.blog_posts")}
            </Title>
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
          <Title order={2}>
            {tFaq("header")}
          </Title>
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
        <Container size="lg">
          <h2 className="mb-4">{t("home.index.join")}</h2>
          <Grid gutter="md">
            <Grid.Col span={{ base: 12, md: 6 }}>
              <SignUpFormBlock user={newUser} />
            </Grid.Col>
          </Grid>
        </Container>
      )}

      {!user.guest && i18next.language === 'ru' && (
        <Container size="lg" mt={100}>
          <Grid align="center" justify="space-between" gutter={0}>
            <Grid.Col span={{ base: 12, md: 6 }}>
              <Center>
                <Text fz={40} fw="bold">{t("home.index.consultation")}</Text>
              </Center>
            </Grid.Col>
            <Grid.Col span={{ base: 12, md: 4 }}>
              <Card withBorder shadow="sm" p="xl">
                <LeadFormBlock lead={lead} />
              </Card>
            </Grid.Col>
          </Grid>
        </Container>
      )}

    </ApplicationLayout>
  );
}
