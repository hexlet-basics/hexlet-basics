import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import type {
  LanguageLandingPageForLists,
  LanguageMember,
} from "@/types/serializers";
import { usePage } from "@inertiajs/react";
import type { PropsWithChildren } from "react";
import { Button, Card } from "react-bootstrap";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren & {
  // course: Language;
  landingPage: LanguageLandingPageForLists;
  courseMember?: LanguageMember;
  continueButton?: boolean;
};

export default function CourseBlock({
  landingPage,
  courseMember,
  continueButton,
}: Props) {
  const { suffix } = usePage<SharedProps>().props;
  const { t: tHelpers } = useTranslation("helpers");

  return (
    <>
      <Card className="border-0 border-bottom shadow-sm bg-body-tertiary h-100">
        <Card.Img
          loading="lazy"
          variant="top"
          src={landingPage.language.cover_list_variant}
          alt={landingPage.header}
        />
        {/* {languageMembersByLanguage && */}
        {/*   languageMembersByLanguage[language.id]?.finished && ( */}
        {/*     <div className="card-img-overlay m-auto bg-light text-center shadow-sm p-1"> */}
        {/*       {t(".course_finished")} */}
        {/*     </div> */}
        {/*   )} */}
        <Card.Body className="p-2 p-md-3">
          <Card.Title>
            <a
              href={Routes.language_url(landingPage.slug)}
              className="stretched-link text-decoration-none link-body-emphasis h2"
            >
              {courseMember && <i className="me-3 bi bi-trophy" />}
              {landingPage.header}
            </a>
          </Card.Title>
          <div className="text-muted">
            <span className="text-nowrap d-block d-sm-inline-block me-4">
              <span className="bi bi-clock me-2" />
              {landingPage.duration}
            </span>
            <span className="text-nowrap d-block d-sm-inline-block">
              <span className="bi bi-people me-2" />
              {landingPage.members_count}
            </span>
          </div>
        </Card.Body>
      </Card>
      {continueButton && (
        <Button variant="outline-primary" className="mt-2 w-100">
          {tHelpers("continue")}
        </Button>
      )}
    </>
  );
}
