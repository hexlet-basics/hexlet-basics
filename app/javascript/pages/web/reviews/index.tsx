import { usePage } from "@inertiajs/react";
import { Alert, Card, Container, Group, SimpleGrid, Text } from "@mantine/core";
import { IconUserCircle } from "@tabler/icons-react";
import dayjs from "dayjs";
import i18next from "i18next";
import type { PropsWithChildren } from "react";
import { Trans, useTranslation } from "react-i18next";
import AppAnchor from "@/components/Elements/AppAnchor";
import XPaging from "@/components/XPaging";
import ApplicationLayout from "@/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { BreadcrumbItem } from "@/types";
import type { Pagy, Review } from "@/types/serializers";

type Props = PropsWithChildren & {
  reviews: Review[];
  pagy: Pagy;
};

export default function New({ reviews, pagy }: Props) {
  const { suffix } = usePage().props;
  const { t } = useTranslation();
  const header = t(($) => $.reviews.index.header);

  const items: BreadcrumbItem[] = [
    {
      name: header,
      url: Routes.blog_posts_url({ suffix }),
    },
  ];

  return (
    <ApplicationLayout items={items} header={header}>
      <Container>
        {i18next.language === "ru" && (
          <Alert mb="xl">
            <Trans
              t={t}
              i18nKey={($) => $.reviews.index.add_review}
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
                <IconUserCircle size={20} />
                <Text fw={500} size="lg">
                  {review.full_name}
                </Text>
              </Group>
              <Text mb="sm">{review.body}</Text>
              <Group justify="space-between" c="dimmed" mt="auto">
                <Text size="sm">{dayjs(review.created_at).format("LL")}</Text>
              </Group>
            </Card>
          ))}
        </SimpleGrid>
        <XPaging pagy={pagy} />
      </Container>
    </ApplicationLayout>
  );
}
