import type { Language } from "@/types/serializers";
import * as Routes from "@/routes.js";
import type { PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { assetPath } from "@/lib/utils";
import { Card } from "react-bootstrap";

type Props = PropsWithChildren & {
	course: Language;
};

export default function CourseBlock({ course }: Props) {
	const { t } = useTranslation();
	const { t: tJS } = useTranslation("js");
	return (
		<Card className="h-100 shadow-sm">
			<Card.Img
        variant="top"
				src={assetPath(course.cover!)}
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
						href={Routes.language_path(course.slug!)}
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
