import { getImageUrl } from "@/images";
import * as Routes from "@/routes.js";
import type { Language } from "@/types/serializers";
import type { SharedProps } from "@/types/types";
import { usePage } from "@inertiajs/react";
import type { PropsWithChildren } from "react";
import { Card } from "react-bootstrap";
import { useTranslation } from "react-i18next";

type Props = PropsWithChildren & {
  course: Language;
};

export default function CourseBlock({ course }: Props) {
  const { t } = useTranslation();
  const { t: tJS } = useTranslation("js");
  const { suffix } = usePage<SharedProps>().props;

  return (
    <Card className="h-100 shadow-sm bg-body-tertiary">
      <Card.Img
        variant="top"
        src={getImageUrl(course.cover!)}
        alt={course.name!}
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
            href={Routes.language_path(course.slug!, { suffix })}
            className="stretched-link text-decoration-none link-body-emphasis"
          >
            {course.name}
          </a>
        </Card.Title>
        <div className="text-muted">
          <span className="text-nowrap d-inline-block me-4">
            <span className="bi bi-clock me-2" />
            {tJS("hours", { count: course.duration })}
          </span>
          <span className="text-nowrap d-inline-block">
            <span className="bi bi-people me-2" />
            {course.members_count}
          </span>
        </div>
      </Card.Body>
    </Card>
  );
}
