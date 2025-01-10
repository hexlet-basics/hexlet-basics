import cn from "classnames";
import _ from "lodash";
import type { InputHTMLAttributes } from "react";
import { Form } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import {
  type FormProps,
  Form as InertiaForm,
  type NestedObject,
  useInertiaInput,
} from "use-inertia-form";

type Props = InputHTMLAttributes<HTMLInputElement> & {
  model?: string;
  name: string;
};

export const XForm = <TForm extends NestedObject>({
  children,
  railsAttributes = true,
  className,
  ...props
}: FormProps<TForm>) => {
  return (
    <InertiaForm
      className={`form ${className}`}
      railsAttributes={railsAttributes}
      {...props}
    >
      {children}
    </InertiaForm>
  );
};

export function XInput({ name, model, ...props }: Props) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tAm } = useTranslation("activemodel");

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput({
    name,
    model,
  });
  // console.log(form, value)

  const errors = error ? _.castArray(error) : [];

  const path = `attributes.${form.model}.${name}`;
  const label = tAr(path, tAm(path));

  const controlClasses = cn({
    "is-invalid": error,
  });

  return (
    <Form.Group className="mb-4">
      <Form.FloatingLabel label={label}>
        <Form.Control
          name={inputName}
          value={value}
          id={inputId}
          placeholder={label}
          className={controlClasses}
          onChange={(e) => setValue(e.target.value)}
          {...props}
        />
        {error && (
          <Form.Control.Feedback type="invalid">
            {errors.map((e) => (
              <div key={e}>{e}</div>
            ))}
          </Form.Control.Feedback>
        )}
      </Form.FloatingLabel>
    </Form.Group>
  );
}

export function XCheck({ name, model, type, ...props }: Props) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tAm } = useTranslation("activemodel");

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput({
    name,
    model,
  });
  // console.log(form, value)

  const errors = error ? _.castArray(error) : [];

  const path = `attributes.${form.model}.${name}`;
  const label = tAr(path, tAm(path));

  const controlClasses = cn({
    "is-invalid": error,
  });

  return (
    <Form.Group className="mb-4">
      <Form.Check
        label={label}
        name={inputName}
        type="checkbox"
        defaultChecked={value as boolean}
        id={inputId}
        placeholder={label}
        className={controlClasses}
        onChange={(e) => setValue(e.target.checked)}
        {...props}
      />
      {error && (
        <Form.Control.Feedback type="invalid">
          {errors.map((e) => (
            <div key={e}>{e}</div>
          ))}
        </Form.Control.Feedback>
      )}
    </Form.Group>
  );
}
