import { Link, usePage } from '@inertiajs/react';
import {
  Box,
  Button,
  Container,
  Divider,
  Grid,
  Image,
  Paper,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import i18next from 'i18next';
import { Compass, FileText, List, Users } from 'lucide-react';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import LeadFormBlock from '@/components/LeadFormBlock';
import bookCoverImg from '@/images/profession-developer-book-cover.webp';
import bookToc from '@/lib/book.ts';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout.tsx';
import * as Routes from '@/routes.js';
import type { LeadCrud, SharedProps } from '@/types';

interface Props extends PropsWithChildren {
  bookRequested: boolean;
  lead: LeadCrud;
}

export default function Show({ bookRequested, lead }: Props) {
  const { t } = useTranslation();
  const { auth } = usePage<SharedProps>().props;

  const features = [
    {
      key: 'direction',
      title: 'books.show.features.direction',
      explanation: 'books.show.features.direction_explanation',
      icon: Compass,
    },
    {
      key: 'plan',
      title: 'books.show.features.plan',
      explanation: 'books.show.features.plan_explanation',
      icon: List,
    },
    {
      key: 'resume',
      title: 'books.show.features.resume',
      explanation: 'books.show.features.resume_explanation',
      icon: FileText,
    },
    {
      key: 'interview',
      title: 'books.show.features.interview',
      explanation: 'books.show.features.interview_explanation',
      icon: Users,
    },
  ] as const;

  return (
    <ApplicationLayout>
      <Container size="lg" my="xl">
        <Grid>
          <Grid.Col span={{ base: 12, md: 7 }}>
            <Text c="dimmed" fz="lg">
              {t('books.show.freebook')}
            </Text>
            <Title order={1} mb="lg">
              {t('books.show.header')}
            </Title>
            <Text fz="lg">{t('books.show.description')}</Text>

            {!bookRequested ? (
              <Button
                component={Link}
                href={Routes.create_request_book_url()}
                mt="md"
                method="post"
                variant="outline"
                size="lg"
              >
                {t('books.show.request')}
              </Button>
            ) : (
              <Button
                component="a"
                href={Routes.download_book_url()}
                target="_blank"
                mt="md"
                variant="outline"
                size="lg"
              >
                {t('books.show.download')}
              </Button>
            )}

            <Grid mt={40} gutter="md">
              {features.map(({ key, title, explanation, icon: Icon }) => (
                <Grid.Col span={{ base: 12, md: 6 }} key={key}>
                  <Stack gap={4}>
                    <Box style={{ display: 'flex', alignItems: 'center' }}>
                      <Icon size={20} style={{ marginRight: 8 }} />
                      <Text fw={700} fz="lg">
                        {t(title)}
                      </Text>
                    </Box>
                    <Text>{t(explanation)}</Text>
                  </Stack>
                </Grid.Col>
              ))}
            </Grid>
          </Grid.Col>

          <Grid.Col span={{ base: 12, md: 5 }}>
            <Image
              src={bookCoverImg}
              alt="Book cover"
              radius="md"
              fit="contain"
            />
          </Grid.Col>
        </Grid>

        <Title order={2} mt={40} mb={20}>
          {t('books.show.toc')}
        </Title>

        {bookToc.map((item, index) => (
          <Box key={item.title} py="md">
            <Grid align="start">
              <Grid.Col span={{ base: 12, md: 2 }}>
                <Text fw={700}>
                  {t('books.show.chapter', { number: index + 1 })}
                </Text>
              </Grid.Col>
              <Grid.Col span={{ base: 12, md: 4 }}>{item.title}</Grid.Col>
              <Grid.Col span={{ base: 12, md: 6 }}>
                <ul style={{ margin: 0, paddingLeft: 20 }}>
                  {item.subsections.map((subsection) => (
                    <li key={subsection}>{subsection}</li>
                  ))}
                </ul>
              </Grid.Col>
            </Grid>
            {index !== bookToc.length - 1 && <Divider my="md" />}
          </Box>
        ))}

        {!auth.user.guest && i18next.language === 'ru' && (
          <Grid align="center" mt={60}>
            <Grid.Col span={{ base: 12, lg: 7 }}>
              <Title order={2}>{t('home.index.consultation')}</Title>
            </Grid.Col>
            <Grid.Col span={{ base: 12, lg: 5 }}>
              <Paper p="lg" radius="md" withBorder>
                <LeadFormBlock leadDto={lead} />
              </Paper>
            </Grid.Col>
          </Grid>
        )}
      </Container>
    </ApplicationLayout>
  );
}
