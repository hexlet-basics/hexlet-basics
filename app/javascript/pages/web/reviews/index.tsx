import type { PropsWithChildren } from "react";
import { Card, Container } from "react-bootstrap";

import Paging from "@/components/Paging";
import Application from "@/pages/layouts/Application";
import type {
	Language,
	LanguageCategory,
	Pagy,
	Review,
} from "@/types/serializers";

type Props = PropsWithChildren & {
	languageCategories: LanguageCategory[];
	languages: Language[];
	reviews: Review[];
	pagy: Pagy;
};

export default function New({
	languageCategories,
	languages,
	reviews,
	pagy,
}: Props) {
	return (
		<Application languageCategories={languageCategories} languages={languages}>
			<Container>
				{reviews.map((review) => (
					<Card key={review.id} className="mb-3">
						<Card.Title>{review.full_name}</Card.Title>
						<Card.Body>{review.body}</Card.Body>
					</Card>
				))}
				<Paging pagy={pagy} />
			</Container>
		</Application>
	);
}
