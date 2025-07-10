import { router, usePage } from '@inertiajs/react';
import {
  Button,
  Card,
  Center,
  Container,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import type { PropsWithChildren } from 'react';
import { useTranslation } from 'react-i18next';

import ApplicationLayout from '@/pages/layouts/ApplicationLayout';
import * as Routes from '@/routes.js';
import type { Survey, SurveyItem, SurveyScenario } from '@/types';

type Props = PropsWithChildren & {
  survey: Survey;
  scenario: SurveyScenario;
  surveyItems: SurveyItem[];
};

export default function Show({ scenario, survey, surveyItems }: Props) {
  const { t } = useTranslation();

  const page = usePage();
  // is there a better way?
  const queryParams = new URLSearchParams(page.url.split('?')[1]);

  const handleAnswer = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    const id = Number((e.currentTarget as HTMLButtonElement).dataset.id);
    router.post(Routes.scenario_survey_answers_url(scenario.id, survey.id), {
      from: queryParams.get('from'),
      survey_answer: { survey_item_id: id },
    });
  };

  return (
    <ApplicationLayout>
      <Container my="xl">
        <Center>
          <Card
            withBorder
            radius="lg"
            p="xl"
            w={{ base: '100%', sm: '90%', md: '80%', lg: '70%' }}
          >
            <Stack gap="md">
              <Title order={1} mb="lg">
                {survey.question}
              </Title>
              {survey.description && <Text mb="lg">{survey.description}</Text>}

              <Stack gap="xs">
                {surveyItems.map((item) => (
                  <Button
                    key={item.id}
                    data-id={item.id}
                    onClick={handleAnswer}
                    variant="outline"
                  >
                    {item.value}
                  </Button>
                ))}
              </Stack>
            </Stack>
          </Card>
        </Center>
      </Container>
    </ApplicationLayout>
  );
}
