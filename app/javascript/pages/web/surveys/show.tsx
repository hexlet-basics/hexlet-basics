import { type PropsWithChildren } from "react";
import { Button, Card, Container, Row } from "react-bootstrap";

import * as Routes from "@/routes.js";
import { useTranslation } from "react-i18next";

import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import { Survey, SurveyItem } from "@/types";
import { router, usePage } from "@inertiajs/react";

type Props = PropsWithChildren & {
  survey: Survey
  surveyItems: SurveyItem[]
};

export default function Show({ survey, surveyItems }: Props) {
  const { t } = useTranslation();

  const page = usePage()
  // is there a better way?
  const queryParams = new URLSearchParams(page.url.split('?')[1]);

  const handleAnswer = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    const id = Number((e.currentTarget as HTMLButtonElement).dataset.id);
    router.post(
      Routes.survey_answers_url(survey.slug!),
      { from: queryParams.get('from'), survey_answer: { survey_item_id: id } }
    )
  }

  return (
    <ApplicationLayout>
      <Container className="my-5">
        <Row className="justify-content-center">
          <div className="col-sm-9 col-md-8 col-lg-7">
            <Card className="p-3">
              <Card.Body>
                <div className="mb-4">
                  <Card.Title as="h1">
                    {survey.question}
                  </Card.Title>
                  {survey.description && <p>{survey.description}</p>}
                </div>
                {surveyItems.map((item) => (
                  <div
                    className="mb-2"
                    key={item.id}
                  >
                    <Button
                      data-id={item.id}
                      onClick={handleAnswer}
                      variant="outline-secondary"
                    >{item.value}</Button>
                  </div>
                ))}
              </Card.Body>
            </Card>
          </div>
        </Row>
      </Container>
    </ApplicationLayout>
  );
}
