import { Head, usePage } from '@inertiajs/react';
import {
  Anchor,
  Box,
  Container,
  Divider,
  Grid,
  Group,
  Stack,
  Text,
} from '@mantine/core';
import i18next from 'i18next';
import { Github, Send } from 'lucide-react';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';
import type { Organization, WithContext } from 'schema-dts';
import AppAnchor from '@/components/Elements/AppAnchor';
import * as Routes from '@/routes.js';
import type { SharedProps } from '@/types';

function FooterLink(
  props: {
    href: string;
    external?: boolean;
    pseudo?: boolean;
  } & PropsWithChildren,
) {
  const { href, children, external = false, pseudo = false } = props;

  return (
    <AppAnchor c="dimmed" href={href} external={external} pseudo={pseudo}>
      {children}
    </AppAnchor>
  );
}

export default function FooterBlock() {
  const { t: tLayouts } = useTranslation('layouts');
  const { t: tCommon } = useTranslation('common');
  const { landingPagesForFooter, courseCategories } =
    usePage<SharedProps>().props;

  const organization: WithContext<Organization> = {
    '@context': 'https://schema.org',
    '@type': 'Organization',
    email: tCommon('organization.email'),
    name: tCommon('organization.legal_name'),
    telephone: tCommon('organization.phone'),
  };

  return (
    <>
      <Head>
        <script type="application/ld+json">
          {JSON.stringify(organization)}
        </script>
      </Head>

      <Box bg="gray.0" mt={100} pt="lg">
        <Container size="lg" pt="lg">
          <footer>
            <Grid justify="space-b">
              {i18next.language === 'ru' && (
                <Grid.Col span={{ base: 12, xs: 6, md: 3 }}>
                  <Stack gap="sm">
                    <Anchor href="tel:+78001002247" fz="h3">
                        8 800 100 22 47
                    </Anchor>
                    <Anchor href="tel:+74950852162" fz="h3">
                        +7 495 085 21 62
                    </Anchor>
                    <Anchor href="mailto:support@hexlet.io">
                        support@hexlet.io
                    </Anchor>
                  </Stack>
                  <Stack mt="sm" gap={0}>
                    <AppAnchor href={Routes.root_path()} mb="xs">
                      ООО «Хекслет Рус»
                    </AppAnchor>
                    <Text>
                      108813 г. Москва, вн.тер.г. поселение Московский, г.
                      Московский
                    </Text>
                    <Text mb="xs">
                      ул. Солнечная, д. 3А, стр. 1, помещ. 10/3
                    </Text>
                    <Text>ОГРН 1217300010476</Text>
                  </Stack>
                </Grid.Col>
              )}

              <Grid.Col span={{ base: 12, xs: 6, md: 3 }}>
                <Text fw="bold">{tLayouts('shared.footer.courses')}</Text>
                <Stack gap={3}>
                  {landingPagesForFooter.map((lp) => (
                    <FooterLink
                      key={lp.id}
                      href={Routes.language_path(lp.slug)}
                    >
                      {lp.header}
                    </FooterLink>
                  ))}
                </Stack>
              </Grid.Col>

              <Grid.Col span={{ base: 12, xs: 6, md: 3 }}>
                <Stack gap={3}>
                  <Text fw="bold">{tLayouts('shared.footer.categories')}</Text>
                  {courseCategories.map((category) => (
                    <FooterLink
                      key={category.id}
                      href={Routes.language_category_path(category.slug!)}
                    >
                      {category.header}
                    </FooterLink>
                  ))}
                </Stack>
              </Grid.Col>

              <Grid.Col
                span={{ base: 12, xs: 6, md: 3 }}
                order={{ base: 1, lg: 0 }}
              >
                <Stack gap={3}>
                  {i18next.language === 'ru' && (
                    <FooterLink href={Routes.map_path()}>
                      {tLayouts('shared.footer.sitemap')}
                    </FooterLink>
                  )}
                  <FooterLink href={Routes.page_path('about')} pseudo>
                    {tLayouts('shared.footer.about')}
                  </FooterLink>
                  <FooterLink href={Routes.blog_posts_path()} pseudo>
                    {tLayouts('shared.footer.blog')}
                  </FooterLink>
                  <FooterLink href={Routes.reviews_path()}>
                    {tLayouts('shared.footer.reviews')}
                  </FooterLink>
                  <FooterLink pseudo href={Routes.page_path('authors')}>
                    {tLayouts('shared.footer.authors')}
                  </FooterLink>
                  <FooterLink pseudo href={Routes.page_path('tos')}>
                    {tLayouts('shared.footer.tos')}
                  </FooterLink>
                  <FooterLink pseudo href={Routes.page_path('privacy')}>
                    {tLayouts('shared.footer.privacy')}
                  </FooterLink>
                  <FooterLink pseudo href={Routes.page_path('cookie_policy')}>
                    {tLayouts('shared.footer.cookie_policy')}
                  </FooterLink>
                </Stack>
              </Grid.Col>
            </Grid>

            <Divider my="xl" />

            <Group justify="space-between" align="center">
              <Text>{`© ${new Date().getFullYear()}`}</Text>
              <Group gap="md">
                <AppAnchor
                  href="https://ttttt.me/hexlet_ru"
                  external
                  underline="never"
                  aria-label="Hexlet Telegram Channel"
                >
                  <Send />
                </AppAnchor>
                <AppAnchor
                  href="https://github.com/hexlet-basics"
                  external
                  underline="never"
                  aria-label="Project Repository On Github"
                >
                  <Github />
                </AppAnchor>
              </Group>
            </Group>
          </footer>
        </Container>
      </Box>
    </>
  );
}
