import dayjs from "dayjs";
import type { PropsWithChildren } from "react";
import { Card, Col, Container, Row } from "react-bootstrap";

import XPaging from "@/components/XPaging";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import type { Pagy, Review } from "@/types/serializers";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren & {
  reviews: Review[];
  pagy: Pagy;
};

export default function New({ reviews, pagy }: Props) {
  const { t } = useTranslation();

  return (
    <ApplicationLayout header={t("reviews.index.title")}>
      <Container>
        <Row className="mb-5 row-cols-1 row-cols-lg-2">
          {reviews.map((review) => (
            <Col key={review.id} className="mb-5">
              <Card key={review.id} className="mb-3 h-100">
                <Card.Body className="d-flex flex-column">
                  <div className="mb-auto">
                    <Card.Title>{review.full_name}</Card.Title>
                    <Card.Text>{review.body}</Card.Text>
                  </div>
                  <div className="d-flex">
                    {/* {review.language && ( */}
                    {/*   <Card.Link */}
                    {/*     href={Routes.language_path(review.language.slug!)} */}
                    {/*     className="me-auto link-body-emphasis" */}
                    {/*   > */}
                    {/*     {review.language.name} */}
                    {/*   </Card.Link> */}
                    {/* )} */}
                    <div>{dayjs(review.created_at).format("YYYY-MM-DD")}</div>
                  </div>
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
        <XPaging pagy={pagy} />
      </Container>
    </ApplicationLayout>
  );
}
