import { Head, Link, usePage } from '@inertiajs/react';
import {
  Accordion,
  Box,
  Button,
  Card,
  Center,
  Container,
  Flex,
  Grid,
  Group,
  Image,
  List,
  NumberFormatter,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import dayjs from 'dayjs';
import i18next from 'i18next';
import { Clock, Users } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import type { Product, WithContext } from 'schema-dts';
import AppAnchor from '@/components/Elements/AppAnchor';
import LeadFormBlock from '@/components/LeadFormBlock';
import MarkdownViewer from '@/components/MarkdownViewer';
import SignUpFormBlock from '@/components/SignUpFormBlock';
import XssContent from '@/components/XssContent';
import codeIllustration from '@/images/code.svg';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import { getResourceUrl } from '@/resources';
import * as Routes from '@/routes.js';
import type { BreadcrumbItem, Language, LeadCrud, SharedProps } from '@/types';
import type {
  LanguageCategory,
  LanguageLandingPage,
  LanguageLandingPageQnaItem,
  LanguageLesson,
  LanguageMember,
  LanguageModule,
  Review,
  UserSignUpForm,
} from '@/types/serializers';

type Props = {
  lead: LeadCrud;
  course: Language;
  courseMember?: LanguageMember;
  courseCategory?: LanguageCategory;
  courseLandingPage: LanguageLandingPage;
  qnaItems: LanguageLandingPageQnaItem[];
  firstLesson: LanguageLesson;
  newUser: UserSignUpForm;
  nextLesson?: LanguageLesson;
  courseModules: LanguageModule[];
  lessonsByModuleId: {
    [moduleId: number]: LanguageLesson[];
  };
  reviews: Review[];
};

export default function Show({
  firstLesson,
  nextLesson,
  courseLandingPage,
  courseMember,
  courseCategory,
  courseModules,
  lessonsByModuleId,
  course,
  qnaItems,
  newUser,
  lead,
}: Props) {
  const { t } = useTranslation();
  const {
    auth: { user },
    locale,
  } = usePage<SharedProps>().props;

  const productSchema: WithContext<Product> = {
    '@context': 'https://schema.org',
    '@type': 'Product',
    description: courseLandingPage.description,
    image: course.cover_list_variant,
    offers: {
      '@type': 'Offer',
      price: '0',
      priceCurrency: 'RUB',
      availability: 'https://schema.org/InStock',
    },
    aggregateRating: {
      '@type': 'AggregateRating',
      ratingValue: course.rating_value,
      ratingCount: course.rating_count,
    },
    name: courseLandingPage.header,
  };

  const breadcrumbItems: BreadcrumbItem[] = [
    {
      name: courseCategory?.name ?? '-',
      url: courseCategory
        ? Routes.language_category_url(courseCategory.slug!)
        : '#',
    },
    {
      name: courseLandingPage.header,
      url: Routes.language_url(courseLandingPage.slug),
    },
  ];

  return (
    <ApplicationLayout items={breadcrumbItems}>
      <Head>
        <script type="application/ld+json">
          {JSON.stringify(productSchema)}
        </script>
      </Head>

      <Container size="lg">
        <Grid>
          <Grid.Col span={{ base: 12, sm: 7 }}>
            <Text size="sm" c="dimmed">
              {t('languages.show.free_course', {
                name: courseLandingPage.header,
              })}
            </Text>
            <Title order={1} mb="lg" fz={48}>
              {courseLandingPage.header}
            </Title>
            <Text size="lg" mb="sm">
              {courseLandingPage.description}
            </Text>
            <Group mb="xl">
              <Group me="lg">
                <Users size={16} />
                <NumberFormatter
                  thousandSeparator
                  value={course.members_count}
                />
              </Group>
              <Group>
                <Clock size={16} />
                {t('languages.show.updated_at', {
                  date: dayjs(course.current_version!.created_at).format('LL'),
                })}
              </Group>
            </Group>

            <Group>
              {!courseMember && (
                <Button
                  size="lg"
                  component={Link}
                  href={Routes.language_lesson_path(
                    course.slug!,
                    firstLesson.slug,
                  )}
                >
                  {t('languages.show.try')}
                </Button>
              )}
              {courseMember && nextLesson && (
                <Button
                  size="lg"
                  component={Link}
                  href={Routes.language_lesson_path(
                    course.slug!,
                    nextLesson.slug!,
                  )}
                >
                  {t('languages.show.continue')}
                </Button>
              )}
              {courseMember?.state === 'finished' && !nextLesson && (
                <Button
                  size="lg"
                  variant="outline"
                  component={Link}
                  href={Routes.language_lesson_path(
                    course.slug!,
                    firstLesson.slug!,
                  )}
                >
                  {t('languages.show.restart')}
                </Button>
              )}
              {courseMember?.state === 'finished' &&
                !nextLesson &&
                course.hexlet_program_landing_page && (
                  <Button
                    size="lg"
                    component="a"
                    target="_blank"
                    href={`${course.hexlet_program_landing_page}?utm_source=code-basics&utm_medium=referral`}
                  >
                    {t('languages.show.hexlet_program_link')}
                  </Button>
                )}
            </Group>
          </Grid.Col>
          <Grid.Col visibleFrom="sm" span={{ base: 12, sm: 5 }}>
            <Image
              src={codeIllustration}
              // w="auto"
              fit="cover"
              alt={t('languages.show.cover_image')}
            />
          </Grid.Col>
        </Grid>

        <Grid my="xl">
          <Grid.Col span={{ base: 12, lg: 9 }}>
            <Title order={2} size="h1" mb="md" responsive>
              {courseLandingPage.used_in_header}
            </Title>
            <Text>{courseLandingPage.used_in_description}</Text>
          </Grid.Col>
        </Grid>

        <SimpleGrid cols={{ base: 1, lg: 2 }} py="xl">
          {courseLandingPage.outcomes_image && (
            <Image
              src={courseLandingPage.outcomes_image}
              width="100%"
              height="auto"
              loading="lazy"
              alt={t('languages.show.learning_preview')}
              style={{
                borderRadius: 'var(--mantine-radius-xl)',
                boxShadow: 'var(--mantine-shadow-lg)',
              }}
            />
          )}

          <Box>
            <Title order={2} size="h1" mb="md" responsive>
              {courseLandingPage.outcomes_header}
            </Title>
            <Text>{courseLandingPage.outcomes_description}</Text>
          </Box>
        </SimpleGrid>

        <Card bg="indigo.0" my={{ base: 'lg', sm: 80 }} p="xl">
          <Group>
            <Text responsive="3rem">
              {t('languages.show.course_graduates')}
            </Text>
            <Button
              bg="dark"
              size="xl"
              component={Link}
              href={Routes.language_lesson_path(course.slug!, firstLesson.slug)}
            >
              {t('languages.show.start', { name: courseLandingPage.name })}
            </Button>
          </Group>
        </Card>

        <Box>
          <Title order={2} size="h1" my="xl" responsive>
            {t('languages.show.learning_program', {
              name: courseLandingPage.name,
            })}
          </Title>
          <Accordion defaultValue="0">
            {courseModules.map((m, index) => (
              <Accordion.Item key={m.id} value={index.toString()} py="lg">
                <Accordion.Control>
                  <Title order={3} responsive>
                    {m.name!}
                  </Title>
                </Accordion.Control>
                <Accordion.Panel>
                  <Grid>
                    <Grid.Col span={{ base: 12, xs: 4 }}>
                      <List>
                        {(lessonsByModuleId[m.id] ?? []).map((l) => (
                          <List.Item key={l.id}>
                            <AppAnchor
                              href={Routes.language_lesson_path(
                                course.slug!,
                                l.slug!,
                              )}
                            >
                              {l.name}
                            </AppAnchor>
                          </List.Item>
                        ))}
                      </List>
                    </Grid.Col>
                    <Grid.Col span={{ base: 12, xs: 8 }}>
                      {m.description}
                    </Grid.Col>
                  </Grid>
                </Accordion.Panel>
              </Accordion.Item>
            ))}
          </Accordion>
        </Box>

        <Box my={{ base: 'lg', sm: 80 }}>
          <Title order={2} fz="h1" mb="xl" responsive>
            {t('languages.show.about_learning')}
          </Title>
          <SimpleGrid cols={{ base: 1, lg: 2 }}>
            <Box>
              <Text fw="bold">{t('languages.show.convenient format')}</Text>
              <XssContent mb="md">
                {t('languages.show.learning_conveniently')}
              </XssContent>
              <Text fw="bold">{t('languages.show.browser_practice')}</Text>
              <XssContent mb="md">
                {t('languages.show.real_life_challenges')}
              </XssContent>
              <Text fw="bold">{t('languages.show.ai_without_limits')}</Text>
              <XssContent mb="md">
                {t('languages.show.ai_explanation')}
              </XssContent>
            </Box>

            <Box h="100%" w="100%">
              <video
                className="w-full h-full object-cover"
                src={getResourceUrl(
                  `course-landing-page/learning_${locale}.mp4`,
                )}
                autoPlay
                loop
                muted
                playsInline
              />
            </Box>
          </SimpleGrid>
        </Box>

        <Card bg="indigo.0" p="xl">
          <Grid>
            <Grid.Col span={{ base: 12, sm: 8 }} display="flex">
              <Text responsive="3rem">
                {t('languages.show.demo_description')}
              </Text>
            </Grid.Col>
            <Grid.Col span={{ base: 12, sm: 4 }}>
              <Flex h="100%" justify="end">
                <Button
                  my="auto"
                  size="xl"
                  component={Link}
                  href={Routes.language_lesson_path(
                    course.slug!,
                    firstLesson.slug,
                  )}
                >
                  {t('languages.show.demo_start')}
                </Button>
              </Flex>
            </Grid.Col>
          </Grid>
        </Card>

        {/* {i18next.language === "ru" && ( */}
        {/*   <Box> */}
        {/*     <Grid> */}
        {/*       <Grid.Col className="position-relative p-0"> */}
        {/*         <div className="hexlet-basics-community-image" /> */}
        {/*       </Grid.Col> */}
        {/*       <Grid.Col className="col-lg-7 col-xl-6 p-4 p-md-5"> */}
        {/*         <div className="d-flex flex-column justify-content-center py-3 py-lg-4"> */}
        {/*           <div className="display-5 fw-semibold lh-1 mb-4"> */}
        {/*             {t("languages.show.more_than_support")} */}
        {/*           </div> */}
        {/*           <div className="pe-lg-5"> */}
        {/*             <div className="mb-5 pe-xl-4"> */}
        {/*               {t("languages.show.about_developer_community")} */}
        {/*             </div> */}
        {/*             <a */}
        {/*               className="btn btn-secondary" */}
        {/*               href="https://t.me/HexletLearningBot" */}
        {/*               target="_blank" */}
        {/*               rel="noopener noreferrer" */}
        {/*             > */}
        {/*               <span>{t("languages.show.join")}</span> */}
        {/*             </a> */}
        {/*           </div> */}
        {/*         </div> */}
        {/*       </Grid.Col> */}
        {/*     </Grid> */}
        {/*   </Box> */}
        {/* )} */}

        {qnaItems.length > 0 && (
          <Box my="xl" py="xl">
            <Text fz="h1" fw="bold" mb="xl">
              {t('languages.show.sort_questions')}
            </Text>
            <SimpleGrid cols={{ base: 1, xs: 2 }}>
              {qnaItems.map((item) => (
                <Stack key={item.id} gap={0}>
                  <Text fw="bold">{item.question}</Text>
                  <MarkdownViewer>{item.answer}</MarkdownViewer>
                </Stack>
              ))}
            </SimpleGrid>
          </Box>
        )}

        {!user.guest && i18next.language === 'ru' && (
          <Container size="lg">
            <Grid align="center" justify="space-between" gutter={0}>
              <Grid.Col span={{ base: 12, md: 6 }}>
                <Center mb="xl">
                  <Text fz={40} fw="bold">
                    {t('home.index.consultation')}
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

        {user.guest && (
          <Container size="lg" mt={100}>
            <Grid align="center" justify="space-between" gutter={0}>
              <Grid.Col span={{ base: 12, md: 6 }}>
                <Center mb="lg">
                  <Text fw="bold" responsive="3rem">
                    {t('home.index.join')}
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
      </Container>
    </ApplicationLayout>
  );
}
