import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";

import { Head, Link, usePage } from "@inertiajs/react";
import { Container, Grid, Stack, Text, Group, Anchor, Title, Divider, Box } from '@mantine/core';
import { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import type { Organization, WithContext } from "schema-dts";

function FooterLink(props: { href: string } & PropsWithChildren) {
  const { href, children } = props

  return <Anchor
    c="dimmed"
    component={Link}
    href={href}
  >
    {children}
  </Anchor>
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
        <Container size="lg" p="xl" fz="xs">
          <footer>
            <Grid justify="space-around">
              <Grid.Col span={{ base: 12, sm: 6, md: 3 }}>
                <Text c="dimmed" size="sm">
                  © 2020 mantine.dev. All rights reserved.
                </Text>
              </Grid.Col>
              <Grid.Col span={{ base: 12, sm: 6, md: 3 }}>
                <Stack gap={3}>
                  {landingPagesForFooter.map((lp) => (
                    <FooterLink key={lp.id} href={Routes.language_path(lp.slug)}>
                      {lp.header}
                    </FooterLink>
                  ))}
                </Stack>
              </Grid.Col>

              <Grid.Col span={{ base: 12, sm: 6, md: 3 }}>
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

              <Grid.Col span={{ base: 12, sm: 6, md: 3 }} order={{ base: 1, lg: 0 }}>
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
                </Stack>
              </Grid.Col>
            </Grid>

            <Divider my="xl" />

            <Group justify="space-between" align="center">
              <Text>{`© ${new Date().getFullYear()}`}</Text>
              <Group gap="md">
                <Anchor
                  href="https://github.com/hexlet-basics"
                  target="_blank"
                  rel="noreferrer nofollow"
                  underline="never"
                >
                  <i className="bi bi-github" />
                </Anchor>
              </Group>
            </Group>

          </footer>
        </Container>
      </Box>
    </>
  );
}
