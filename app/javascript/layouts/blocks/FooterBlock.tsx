import { Head, usePage } from "@inertiajs/react";
import {
  Anchor,
  Box,
  Container,
  Divider,
  Group,
  SimpleGrid,
  Stack,
  Text,
} from "@mantine/core";
import {
  IconBrandGithub,
  IconBrandYoutube,
  IconSend,
} from "@tabler/icons-react";
import { chunk } from "es-toolkit";
import i18next from "i18next";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import type { Organization, WithContext } from "schema-dts";
import AppAnchor from "@/components/Elements/AppAnchor";
import * as Routes from "@/routes.js";

function FooterLink(
  props: {
    href: string;
    external?: boolean;
    pseudo?: boolean;
  } & PropsWithChildren,
) {
  const { href, children, external = false } = props;

  return (
    <AppAnchor
      fz="sm"
      // c="dimmed"
      href={href}
      external={external}
    >
      {children}
    </AppAnchor>
  );
}

export default function FooterBlock() {
  const { t } = useTranslation();

  const { landingPagesForFooter, courseCategories } = usePage().props;

  const organization: WithContext<Organization> = {
    "@context": "https://schema.org",
    "@type": "Organization",
    email: t(($) => $.common.organization.email),
    name: t(($) => $.common.organization.legal_name),
    telephone: t(($) => $.common.organization.phone),
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
                  {t(($) => $.layouts.shared.footer.codebasics)}
                </Text>
                {i18next.language === "ru" && (
                  <FooterLink href={Routes.map_path()}>
                    {t(($) => $.layouts.shared.footer.sitemap)}
                  </FooterLink>
                )}
                <FooterLink href={Routes.page_path("about")} pseudo>
                  {t(($) => $.layouts.shared.footer.about)}
                </FooterLink>
                <FooterLink href={Routes.blog_posts_path()} pseudo>
                  {t(($) => $.layouts.shared.footer.blog)}
                </FooterLink>
                <FooterLink href={Routes.reviews_path()}>
                  {t(($) => $.layouts.shared.footer.reviews)}
                </FooterLink>
                <FooterLink pseudo href={Routes.page_path("authors")}>
                  {t(($) => $.layouts.shared.footer.authors)}
                </FooterLink>
              </Stack>

              {landingGroups.map((group, index) => (
                // biome-ignore lint/suspicious/noArrayIndexKey: -
                <Stack key={index} gap="sm">
                  <Text fw="bold" fz="sm">
                    {t(($) => $.layouts.shared.footer.courses, {
                      number: index + 1,
                    })}
                  </Text>
                  {group.map((lp) => (
                    <FooterLink
                      key={lp.id}
                      href={Routes.language_path(lp.slug)}
                    >
                      {lp.name}
                    </FooterLink>
                  ))}
                </Stack>
              ))}

              <Stack gap="sm">
                <Text fw="bold" fz="sm">
                  {t(($) => $.layouts.shared.footer.categories)}
                </Text>
                {courseCategories.map((category) => (
                  <FooterLink
                    key={category.id}
                    href={Routes.language_category_path(category.slug!)}
                  >
                    {category.name}
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
                    aria-label="Project Repository On IconBrandGithub"
                  >
                    <IconBrandGithub />
                  </AppAnchor>
                  <AppAnchor
                    href="https://ttttt.me/hexlet_ru"
                    external
                    underline="never"
                    aria-label="Hexlet Telegram Channel"
                  >
                    <IconSend />
                  </AppAnchor>
                  <AppAnchor
                    href="https://www.youtube.com/@HexletOrg"
                    external
                    underline="never"
                    aria-label="Youtbe Channel"
                  >
                    <IconBrandYoutube />
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
                <FooterLink pseudo href={Routes.page_path("tos")}>
                  {t(($) => $.layouts.shared.footer.tos)}
                </FooterLink>
                <FooterLink pseudo href={Routes.page_path("privacy")}>
                  {t(($) => $.layouts.shared.footer.privacy)}
                </FooterLink>
                <FooterLink pseudo href={Routes.page_path("cookie_policy")}>
                  {t(($) => $.layouts.shared.footer.cookie_policy)}
                </FooterLink>
              </Stack>

              {i18next.language === "ru" && (
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
