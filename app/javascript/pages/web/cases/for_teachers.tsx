import { getResourceUrl } from "@/resources";
import learningEnVideo from "@/images/course-landing-page/learning_en.mp4";
import learningRuVideo from "@/images/course-landing-page/learning_ru.mp4";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import { usePage, Link } from "@inertiajs/react";
import i18next from "i18next";
import { useTranslation } from "react-i18next";
import {
  Box,
  Button,
  Card,
  Container,
  Group,
  Image,
  List,
  SimpleGrid,
  Stack,
  Text,
  Title,
} from "@mantine/core";

export default function ForTeachersBlock() {
  const {
    auth: { user },
  } = usePage<SharedProps>().props;
  const { t } = useTranslation();

  const interactiveApproachList = t(
    "cases.for_teachers.interactive_approach_list",
    { returnObjects: true }
  );
  const earlyCareerGuidanceList = t(
    "cases.for_teachers.early_career_guidance_list",
    { returnObjects: true }
  );
  const howToLearnCards = t(
    "cases.for_teachers.how_to_learn_programming_cards",
    { returnObjects: true }
  );

  return (
    <ApplicationLayout>
      <Container py="md" size="lg">

        <Title order={1} mb="sm">{t("cases.for_teachers.header")}</Title>
        <Text size="lg" c="dimmed" mb="xl">
          {t("cases.for_teachers.description")}
        </Text>
        <Button
          component={Link}
          href={Routes.root_path()}
          size="lg"
          variant="filled"
          px="xl"
          mb={50}
        >
          {t("cases.for_teachers.try")}
        </Button>



        <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="xl" my="xl">
          <Title order={2} pe="xl">{t("cases.for_teachers.integrate_into_education")}</Title>
          <Text size="lg" c="dimmed">
            {t("cases.for_teachers.lay_programming_foundations")}
          </Text>
        </SimpleGrid>

        <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="xl" my="xl">
          <Stack>
            <Text fw="bold" fz="h4">{t("cases.for_teachers.interactive_approach")}</Text>
            <List>
              {interactiveApproachList.map((item) => (
                <List.Item mb="sm" key={item}>{item}</List.Item>
              ))}
            </List>
            <Text fw="bold" fz="h4">{t("cases.for_teachers.early_career_guidance")}</Text>
            <List>
              {earlyCareerGuidanceList.map((item, index) => (
                <List.Item mb="sm" key={item}>{item}</List.Item>
              ))}
            </List>
          </Stack>
          <Box>
            <video
              src={i18next.language === "en" ? learningEnVideo : learningRuVideo}
              autoPlay
              loop
              muted
              playsInline
              style={{
                width: '100%',
                height: '60%',
                objectFit: 'cover', // or 'contain' depending on the desired behavior
                borderRadius: 'var(--mantine-radius-md)', // optional
              }}
            />
          </Box>
        </SimpleGrid>

        <Card my="xl" radius="lg" bg="dark" c="white">

          <Title order={2} mb="xl">{t("cases.for_teachers.integrate_now")}</Title>
          <Group justify="space-between" align="end" mt="md" gap="md" wrap="wrap">
            <Text c="gray.5">{t("cases.for_teachers.open_browser_and_sign_up")}</Text>
            {user.guest ? (
              <Button
                component={Link}
                href={Routes.new_user_path()}
                variant="white"
                  c="dark"
                size="lg"
                  px="xl"
              >
                {t("cases.for_teachers.sign_up")}
              </Button>
            ) : (
                <Button
                  component={Link}
                  href={Routes.root_path()}
                  variant="white"
                  c="dark"
                  size="lg"
                  px="xl"
                >
                  {t("cases.for_teachers.select_course")}
                </Button>
              )}
          </Group>
        </Card>

        <Box my="xl">
          <Title order={2} mb="xl">
            {t("cases.for_teachers.how_to_learn_programming")}
          </Title>
          <SimpleGrid cols={{ base: 1, md: 2, xl: 4 }} spacing="lg">
            {howToLearnCards.map((item, index) => (
              <Card key={index} >
                <Image
                  src={getResourceUrl(`for-school-teachers-page/${item.img}.svg`)}
                  alt={item.img}
                  w={90}
                  h={65}
                  mb="xl"
                  loading="lazy"
                  fit="contain"
                />
                <Text mb="sm" fw={500}>{item.title}</Text>
                <Text c="dimmed">{item.subtitle}</Text>
              </Card>
            ))}
          </SimpleGrid>
        </Box>
      </Container>
    </ApplicationLayout>
  );
}

