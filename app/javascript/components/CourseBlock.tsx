import * as Routes from "@/routes.js";
import type { SharedProps } from "@/types";
import type {
  LanguageLandingPageForLists,
  LanguageMember,
} from "@/types/serializers";
import { usePage } from "@inertiajs/react";
import type { PropsWithChildren } from "react";
import { Card } from "react-bootstrap";

type Props = PropsWithChildren & {
  // course: Language;
  landingPage: LanguageLandingPageForLists;
  courseMember?: LanguageMember;
};

export default function CourseBlock({ landingPage, courseMember }: Props) {
  const { suffix } = usePage<SharedProps>().props;

  return (
    <Card className="h-100 shadow-sm bg-body-tertiary">
      <Card.Img
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
      <Card.Body>
        <Card.Title>
          <a
            href={Routes.language_path(landingPage.slug, { suffix })}
            className="stretched-link text-decoration-none link-body-emphasis h2"
          >
            {courseMember && <i className="me-3 bi bi-trophy" />}
            {landingPage.header}
          </a>
        </Card.Title>
        <div className="text-muted">
          <span className="text-nowrap d-inline-block me-4">
            <span className="bi bi-clock me-2" />
            {landingPage.duration}
          </span>
          <span className="text-nowrap d-inline-block">
            <span className="bi bi-people me-2" />
            {landingPage.members_count}
          </span>
        </div>
      </Card.Body>
    </Card>
  );
}
