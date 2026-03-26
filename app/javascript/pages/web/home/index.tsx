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
import BlogPostBlock from "@/components/BlogPostBlock";
import CourseBlock from "@/components/CourseBlock";
import LeadFormBlock from "@/components/LeadFormBlock";
import MarkdownViewer from "@/components/MarkdownViewer";
import SignUpFormBlock from "@/components/SignUpFormBlock";
import { Surface } from "@/components/Surface";
import XssContent from "@/components/XssContent";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import { propsForExternalLink } from "@/lib/utils";
import {
  reviewShowcaseAvatars,
  reviewShowcaseOrder,
} from "@/pages/web/shared/reviews";
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

export default function Index({
  blogPosts,
  newUser,
  lead,
  courseMembersByCourseId,
}: Props) {
  const { t } = useTranslation();

  const baseId = useId("courses");

  const {
    landingPagesForLists,
    auth: { user },
  } = usePage().props;

  const faq = t(($) => $.faq.main, {
    returnObjects: true,
  });
  const reviewTexts = t(($) => $.shared.reviews_showcase, {
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
          sm: "xxl",
        }}
      >
        <Badge
          component="a"
          c="gray"
          style={{ cursor: "pointer" }}
          href="https://github.com/orgs/hexlet-basics/repositories"
          leftSection={<IconBrandGithub size={10} />}
          // autoContrast
          size="sm"
          variant="default"
          {...propsForExternalLink()}
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
          sm: "xxl",
        }}
      >
        <SimpleGrid cols={{ base: 1, xs: 2, md: 3 }} my="xl">
          <Surface color="blue" p="xl">
            <IconBook2 />
            <Text fw="bold" fz="h2">
              {t(($) => $.home.index.hero.courses_count)}
            </Text>
            {t(($) => $.home.index.hero.courses_count_description)}
          </Surface>
          <Surface color="green" p="xl">
            <IconSend />
            <Text fw="bold" fz="h2">
              {t(($) => $.home.index.hero.community_count)}
            </Text>
            {t(($) => $.home.index.hero.community_count_description)}
          </Surface>
          <Surface color="violet" p="xl">
            <IconRobot />
            <Text fw="bold" fz="h2">
              {t(($) => $.home.index.hero.ai_count)}
            </Text>
            {t(($) => $.home.index.hero.ai_count_description)}
          </Surface>
        </SimpleGrid>
      </Container>
      <Container
        size="lg"
        my={{
          base: "xl",
          sm: "xxl",
        }}
      >
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
      <Container
        size="lg"
        my={{
          base: "xl",
          sm: "xxl",
        }}
      >
        <Box mb="xl">
          <Title order={2}>{t(($) => $.home.index.reviews)}</Title>
          <Divider />
        </Box>

        <SimpleGrid spacing="md" cols={{ base: 1, xs: 2, md: 3 }}>
          {reviewShowcaseOrder.map((reviewId) => {
            const review = reviewTexts[reviewId];

            return (
              <Box key={review.name}>
                <Group mb="lg">
                  <Image
                    src={reviewShowcaseAvatars[reviewId]}
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
            );
          })}
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
      <Container
        size="lg"
        my={{
          base: "xl",
          sm: "xxl",
        }}
      >
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
        <Container
          size="lg"
          my={{
            base: "xl",
            sm: "xxl",
          }}
        >
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
