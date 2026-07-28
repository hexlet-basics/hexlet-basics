import {
  Anchor,
  Box,
  Container,
  Divider,
  Group,
  SimpleGrid,
  Stack,
  Text,
  Tooltip,
} from "@mantine/core";
import {
  IconBrandGithub,
  IconBrandYoutube,
  IconSend,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { listCoursesOptions } from "@/client/@tanstack/react-query.gen";

// Split a list into `count` roughly-equal chunks (replaces legacy es-toolkit
// `chunk`, which isn't a dependency of the Go stack).
function chunk<T>(items: T[], size: number): T[][] {
  if (size <= 0) return [items];
  const result: T[][] = [];
  for (let i = 0; i < items.length; i += size) {
    result.push(items.slice(i, i + size));
  }
  return result;
}

// Footer, ported from legacy FooterBlock. Course columns are driven by the
// hey-api generated `listCourses` hook and link to the real course route.
// Legacy also linked out to about/blog/reviews/legal/category pages; those
// aren't ported to the Go stack yet, so their columns are added back as each
// page's route lands (no hardcoded hrefs, no placeholder pages).
export default function Footer() {
  const { t, i18n } = useTranslation();
  const { data: courses } = useQuery(listCoursesOptions());

  const landingPages = courses ?? [];
  const landingGroups = chunk(
    landingPages,
    Math.ceil(landingPages.length / 2) || 1,
  );

  return (
    <Box mt={100} py="lg" fz="sm">
      <Divider mb="xl" />
      <Container size="lg" pt="lg">
        <footer>
          <SimpleGrid cols={{ base: 2, xs: 3 }}>
            <Stack gap="sm">
              <Anchor component={Link} to="/languages" fw="bold" fz="sm">
                {t(($) => $.layouts.shared.all_courses)}
              </Anchor>
            </Stack>

            {landingGroups.map((group, index) => (
              // biome-ignore lint/suspicious/noArrayIndexKey: two stable groups
              <Stack key={index} gap="sm">
                <Text fw="bold" fz="sm">
                  {t(($) => $.layouts.shared.footer.courses, {
                    group: String(index + 1),
                  })}
                </Text>
                {group.map((lp) => (
                  <Anchor
                    key={lp.id}
                    fz="sm"
                    renderRoot={(props) => (
                      <Link
                        to="/languages/$slug"
                        params={{ slug: lp.slug }}
                        {...props}
                      />
                    )}
                  >
                    {lp.name}
                  </Anchor>
                ))}
              </Stack>
            ))}
          </SimpleGrid>

          <Divider my="xl" />

          <SimpleGrid cols={{ base: 2, sm: 3 }}>
            <Stack>
              <Group align="top">
                <Anchor
                  href="https://github.com/hexlet-basics"
                  underline="never"
                  aria-label="Project Repository On GitHub"
                  target="_blank"
                  rel="noreferrer"
                >
                  <IconBrandGithub />
                </Anchor>
                <Tooltip label={t(($) => $.common.telegram_channel)}>
                  <Anchor
                    href="https://t.me/hexlet_ru"
                    underline="never"
                    aria-label={t(($) => $.common.telegram_channel)}
                    target="_blank"
                    rel="noreferrer"
                  >
                    <IconSend />
                  </Anchor>
                </Tooltip>
                <Anchor
                  href="https://www.youtube.com/@HexletOrg"
                  underline="never"
                  aria-label="Youtube Channel"
                  target="_blank"
                  rel="noreferrer"
                >
                  <IconBrandYoutube />
                </Anchor>
              </Group>
              <Anchor href="mailto:support@hexlet.io">support@hexlet.io</Anchor>
            </Stack>

            <Stack gap="xs">
              <Anchor href="tel:+78001002247">8 800 100 22 47</Anchor>
              <Anchor href="tel:+74950852162">+7 495 085 21 62</Anchor>
            </Stack>

            {i18n.language === "ru" && (
              <Stack gap={0}>
                <Anchor
                  href="https://ru.hexlet.io"
                  fz="sm"
                  mb="xs"
                  target="_blank"
                  rel="noreferrer"
                >
                  ООО «Хекслет Рус»
                </Anchor>
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
  );
}
