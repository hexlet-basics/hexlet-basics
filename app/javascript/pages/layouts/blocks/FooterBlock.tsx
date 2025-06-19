import AppAnchor from "@/components/AppAnchor";
import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";

import { Head, usePage } from "@inertiajs/react";
import { Container, Grid, Stack, Text, Group, Anchor, Title, Divider, Box } from '@mantine/core';
import i18next from "i18next";
import { Github, Send } from "lucide-react";
import { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import type { Organization, WithContext } from "schema-dts";

function FooterLink(props: { href: string } & PropsWithChildren) {
  const { href, children } = props

  return <AppAnchor
    c="dimmed"
    href={href}
  >
    {children}
  </AppAnchor>
}

export default function FooterBlock() {
  const { t: tLayouts } = useTranslation("layouts");
  const { t: tCommon } = useTranslation("common");
  const { landingPagesForFooter, courseCategories } = usePage<SharedProps>().props;

  const organization: WithContext<Organization> = {
    "@context": "https://schema.org",
    "@type": "Organization",
    email: tCommon("organization.email"),
    name: tCommon("organization.legal_name"),
    telephone: tCommon("organization.phone"),
  };

  return (
    <>
      <Head>
        <script type="application/ld+json">
          {JSON.stringify(organization)}
        </script>
      </Head>

      <Box bg="gray.0" mt="xl" pt="lg">
        <Container size="lg" pt="lg">
          <footer>
            <Grid justify="space-b">
              {i18next.language == 'ru' && (
                <Grid.Col span={{ base: 12, xs: 6, md: 3 }}>
                  <Text fz="h3" mb="sm">
                    8 800 100 22 47
                  </Text>
                  <Text fz="h3" mb="sm">
                    +7 495 085 21 62
                  </Text>
                  <Anchor href="mailto:support@hexlet.io">
                    support@hexlet.io
                  </Anchor>
                  <Stack mt="sm" gap={0}>
                    <Text mb="xs">
                      ООО «Хекслет Рус»
                    </Text>
                    <Text>
                      108813 г. Москва, вн.тер.г. поселение Московский, г. Московский
                    </Text>
                    <Text mb="xs">
                      ул. Солнечная, д. 3А, стр. 1, помещ. 10/3
                    </Text>
                    <Text>
                      ОГРН 1217300010476
                    </Text>
                  </Stack>
                </Grid.Col>
              )}

              <Grid.Col span={{ base: 12, xs: 6, md: 3 }}>
                <Title order={6}>{tLayouts('shared.footer.courses')}</Title>
                <Stack gap={3}>
                  {landingPagesForFooter.map((lp) => (
                    <FooterLink key={lp.id} href={Routes.language_path(lp.slug)}>
                      {lp.header}
                    </FooterLink>
                  ))}
                </Stack>
              </Grid.Col>

              <Grid.Col span={{ base: 12, xs: 6, md: 3 }}>
                <Stack gap={3}>
                  <Title order={6}>{tLayouts('shared.footer.categories')}</Title>
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

              <Grid.Col span={{ base: 12, xs: 6, md: 3 }} order={{ base: 1, lg: 0 }}>
                <Stack gap={3}>
                  <FooterLink
                    href={Routes.page_path("about")}
                  >
                    {tLayouts("shared.footer.about")}
                  </FooterLink>
                  <FooterLink
                    href={Routes.blog_posts_path()}
                  >
                    {tLayouts("shared.footer.blog")}
                  </FooterLink>
                  <FooterLink
                    href={Routes.reviews_path()}
                  >
                    {tLayouts("shared.footer.reviews")}
                  </FooterLink>
                  <FooterLink
                    href={Routes.page_path("authors")}
                  >
                    {tLayouts("shared.footer.authors")}
                  </FooterLink>
                  <FooterLink
                    href={Routes.page_path("tos")}
                  >
                    {tLayouts("shared.footer.tos")}
                  </FooterLink>
                  <FooterLink
                    href={Routes.page_path("privacy")}
                  >
                    {tLayouts("shared.footer.privacy")}
                  </FooterLink>
                  <FooterLink
                    href={Routes.page_path("cookie_policy")}
                  >
                    {tLayouts("shared.footer.cookie_policy")}
                  </FooterLink>
                </Stack>
              </Grid.Col>
            </Grid>

            <Divider my="xl" />

            <Group justify="space-between" align="center">
              <Text>{`© ${new Date().getFullYear()}`}</Text>
              <Group gap="md">
                <Anchor
                  href="https://ttttt.me/hexlet_ru"
                  target="_blank"
                  rel="noreferrer nofollow"
                  underline="never"
                >
                  <Send />
                </Anchor>
                <Anchor
                  href="https://github.com/hexlet-basics"
                  target="_blank"
                  rel="noreferrer nofollow"
                  underline="never"
                >
                  <Github />
                </Anchor>
              </Group>
            </Group>

          </footer>
        </Container>
      </Box>
    </>
  );
}
