import { usePage } from '@inertiajs/react';
import { Alert, Card, Container, Group, SimpleGrid, Text } from '@mantine/core';
import dayjs from 'dayjs';
import i18next from 'i18next';
import { CircleUser } from 'lucide-react';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import XPaging from '@/components/XPaging';
import XssContent from '@/components/XssContent';
import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type { BreadcrumbItem, SharedProps } from '@/types';
import type { Pagy, Review } from '@/types/serializers';

type Props = PropsWithChildren & {
  reviews: Review[];
  pagy: Pagy;
};

export default function New({ reviews, pagy }: Props) {
  const { suffix } = usePage<SharedProps>().props;
  const { t } = useTranslation();
  const header = t('reviews.index.header');

  const items: BreadcrumbItem[] = [
    {
      name: header,
      url: Routes.blog_posts_url({ suffix }),
    },
  ];

  return (
    <ApplicationLayout items={items} header={header}>
      <Container>
        {i18next.language === 'ru' && (
          <Alert mb="xl">
            <XssContent>{t('reviews.index.add_review')}</XssContent>
          </Alert>
        )}
        <SimpleGrid cols={{ base: 1, xs: 2 }}>
          {reviews.map((review) => (
            <Card
              key={review.id}
              shadow="sm"
              padding="lg"
              radius="md"
              withBorder
            >
              <Group mb="md">
                <CircleUser size={20} />
                <Text fw={500} size="lg">
                  {review.full_name}
                </Text>
              </Group>
              <Text mb="sm">{review.body}</Text>
              <Group justify="space-between" c="dimmed" mt="auto">
                <Text size="sm">{dayjs(review.created_at).format('LL')}</Text>
              </Group>
            </Card>
          ))}
        </SimpleGrid>
        <XPaging pagy={pagy} only={['reviews']} />
      </Container>
    </ApplicationLayout>
  );
}
