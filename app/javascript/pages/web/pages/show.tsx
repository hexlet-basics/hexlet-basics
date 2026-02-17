import { Container, Grid, Title, Typography } from "@mantine/core";
import i18next from "i18next";
import type { PropsWithChildren } from "react";
import ApplicationLayout from "@/layouts/ApplicationLayout.tsx";
import { typographyStyles } from "@/lib/mantine";
import { hasObjectKey } from "@/lib/utils";
import type { User } from "@/types/serializers";
import AboutEn from "./parts/about.en";
import AboutRu from "./parts/about.ru";
import AuthorsEn from "./parts/authors.en";
import AuthorsRu from "./parts/authors.ru";
import CookiePolicyEn from "./parts/cookie_policy.en";
import CookiePolicyRu from "./parts/cookie_policy.ru";
import PrivacyEn from "./parts/privacy.en";
import PrivaryRu from "./parts/privacy.ru";
import TosEn from "./parts/tos.en";
import TosRu from "./parts/tos.ru";

type Props = PropsWithChildren & {
  page: "about" | "tos" | "privacy" | "cookie_policy" | "authors";
  title: string;
  user: User;
};

const mapping = {
  ru: {
    about: AboutRu,
    tos: TosRu,
    privacy: PrivaryRu,
    cookie_policy: CookiePolicyRu,
    authors: AuthorsRu,
  },
  en: {
    about: AboutEn,
    tos: TosEn,
    privacy: PrivacyEn,
    cookie_policy: CookiePolicyEn,
    authors: AuthorsEn,
  },
};

export default function New({ page, title }: Props) {
  const locale = hasObjectKey(mapping, i18next.language)
    ? i18next.language
    : "ru";
  const Component = mapping[locale][page];
  return (
    <ApplicationLayout header={title} center>
      <Container size="xl">
        <Grid justify="center" mb="xl">
          <Grid.Col span={{ base: 12, md: 10, lg: 8 }}>
            <Typography styles={typographyStyles}>
              <Component />
            </Typography>
          </Grid.Col>
        </Grid>
      </Container>
    </ApplicationLayout>
  );
}
