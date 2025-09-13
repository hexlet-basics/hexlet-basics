import { Head, usePage } from '@inertiajs/react';
import {
  Anchor,
  Box,
  Container,
  Divider,
  Grid,
  Group,
  SimpleGrid,
  Stack,
  Text,
} from '@mantine/core';
import { chunk } from 'es-toolkit';
import i18next from 'i18next';
import { Github, Send, Youtube } from 'lucide-react';
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
    <AppAnchor
      fz="sm"
      // c="dimmed"
      href={href}
      external={external}
      pseudo={pseudo}
    >
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

  const landingGroups = chunk(
    landingPagesForFooter,
    Math.ceil(landingPagesForFooter.length / 2),
  );

  return (
    <>
      <Head>
        <script type="application/ld+json">
          {JSON.stringify(organization)}
        </script>
      </Head>

      <Box bg="gray.0" mt={100} py="lg" fz="sm">
        <Container size="lg" pt="lg">
          <footer>
            <SimpleGrid cols={{ base: 2, xs: 4 }}>
              <Stack gap="sm">
                <Text fz="sm" fw="bold">
                  {tLayouts('shared.footer.codebasics')}
                </Text>
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
              </Stack>

              {landingGroups.map((group, index) => (
                // biome-ignore lint/suspicious/noArrayIndexKey: -
                <Stack key={index} gap="sm">
                  <Text fw="bold" fz="sm">
                    {tLayouts('shared.footer.courses', { number: index + 1 })}
                  </Text>
                  {group.map((lp) => (
                    <FooterLink
                      key={lp.id}
                      href={Routes.language_path(lp.slug)}
                    >
                      {lp.header}
                    </FooterLink>
                  ))}
                </Stack>
              ))}

              <Stack gap="sm">
                <Text fw="bold" fz="sm">
                  {tLayouts('shared.footer.categories')}
                </Text>
                {courseCategories.map((category) => (
                  <FooterLink
                    key={category.id}
                    href={Routes.language_category_path(category.slug!)}
                  >
                    {category.header}
                  </FooterLink>
                ))}
              </Stack>
            </SimpleGrid>

            <Divider my="xl" />

            <SimpleGrid cols={{ base: 2, sm: 4 }}>
              <Stack>
                <Group align="top">
                  <AppAnchor
                    href="https://github.com/hexlet-basics"
                    external
                    underline="never"
                    aria-label="Project Repository On Github"
                  >
                    <Github />
                  </AppAnchor>
                  <AppAnchor
                    href="https://ttttt.me/hexlet_ru"
                    external
                    underline="never"
                    aria-label="Hexlet Telegram Channel"
                  >
                    <Send />
                  </AppAnchor>
                  <AppAnchor
                    href="https://www.youtube.com/@HexletOrg"
                    external
                    underline="never"
                    aria-label="Youtbe Channel"
                  >
                    <Youtube />
                  </AppAnchor>
                </Group>
                <Anchor href="mailto:support@hexlet.io">
                  support@hexlet.io
                </Anchor>
              </Stack>

              <Stack gap="xs">
                <Anchor href="tel:+78001002247">8 800 100 22 47</Anchor>
                <Anchor href="tel:+74950852162">+7 495 085 21 62</Anchor>
              </Stack>

              <Stack gap="xs">
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

              {i18next.language === 'ru' && (
                <Stack gap={0}>
                  <AppAnchor
                    href="https://ru.hexlet.io"
                    fz="sm"
                    external
                    mb="xs"
                  >
                    ООО «Хекслет Рус»
                  </AppAnchor>
                  <Text fz="sm">
                    108813 г. Москва, вн.тер.г. поселение Московский, г.
                    Московский, ул. Солнечная, д. 3А, стр. 1, помещ. 10/3
                  </Text>
                  <Text fz="sm">ОГРН 1217300010476</Text>
                </Stack>
              )}
            </SimpleGrid>
          </footer>
        </Container>
      </Box>
    </>
  );
}
