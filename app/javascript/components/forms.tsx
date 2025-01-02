import type BaseModel from "@/types/types";
import type { PropsWithChildren } from "react";
import { Form } from "react-bootstrap";
import { useTranslation } from "react-i18next";

type Props<T> = PropsWithChildren & {
	model: T;
	attribute: keyof T;
};

export function XInput<T extends BaseModel>({ model, attribute }: Props<T>) {
	const { t: tAr } = useTranslation("activerecord");
	const modelName = model.type;
	const label = tAr(`attributes.${modelName}.${attribute}`) as string;

	return (
		<Form.Group className="mb-4">
			<Form.FloatingLabel label={label}>
				<Form.Control placeholder={attribute as string} />
			</Form.FloatingLabel>
		</Form.Group>
	);
}
