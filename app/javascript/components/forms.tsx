import axios from "axios";
import cn from "classnames";
import _ from "lodash";
import {
  AutoComplete,
  type AutoCompleteChangeEvent,
  type AutoCompleteCompleteEvent,
} from "primereact/autocomplete";
import { useEffect, useState, type InputHTMLAttributes } from "react";
import { Form } from "react-bootstrap";
import type { AsProp } from "react-bootstrap/esm/helpers";
import { useTranslation } from "react-i18next";
import {
  type FormProps,
  Form as InertiaForm,
  type NestedObject,
  useInertiaInput,
} from "use-inertia-form";

type Props = InputHTMLAttributes<HTMLInputElement> &
  AsProp & {
    model?: string;
    name: string;
  };

export function XForm<TForm extends NestedObject>({
  children,
  railsAttributes = true,
  className,
  ...props
}: FormProps<TForm>) {
  return (
    <InertiaForm
      className={`form ${className}`}
      railsAttributes={railsAttributes}
      {...props}
    >
      {children}
    </InertiaForm>
  );
}

export function XInput({ name, model, as, ...props }: Props) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tAm } = useTranslation("activemodel");

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput({
    name,
    model,
  });

  // console.log(form)

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
          as={as}
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

type XSelectProps<
  T extends Record<string, unknown>,
  K extends keyof T,
> = Props & {
  items?: T[];
  has?: K;
  source?: string;
  labelField: K;
  valueField: K;
};

export function XSelect<T extends Record<string, unknown>, K extends keyof T>({
  name,
  model,
  items = [],
  has,
  source,
  valueField,
  labelField,
  type,
  ...props
}: XSelectProps<T, K>) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tAm } = useTranslation("activemodel");

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput({
    name,
    model,
  });

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    if (has) {
      const defaultItem = _.get(form.data, [form.model!, has]);
      setSelected(defaultItem);
    } else {
      const defaultItem = _.first(
        filteredItems.filter((i) => i[valueField] === value),
      );
      setSelected(defaultItem);
    }
  }, []);

  const errors = error ? _.castArray(error) : [];

  const path = `attributes.${form.model}.${name}`;
  const label = tAr(path, tAm(path));

  const controlClasses = cn({
    "is-invalid": error,
  });

  const [filteredItems, setFilteredItems] = useState(items);
  const [selected, setSelected] = useState<T | null | undefined>(null);

  const handleChange = (e: AutoCompleteChangeEvent) => {
    setValue(e.value[valueField]);
    setSelected(e.value);
  };

  const search = async (event: AutoCompleteCompleteEvent) => {
    const query = event.query.trim().toLowerCase();
    // console.log("query:", query);
    if (query === "") {
      setFilteredItems([...items]);
    } else {
      if (source) {
        const res = await axios.get<T[]>(source, { params: { query } });
        // console.log(res);
        setFilteredItems(res.data);
      } else {
        const newFilteredItems = items.filter((item) =>
          String(item[labelField]).toLowerCase().startsWith(query),
        );
        setFilteredItems(newFilteredItems);
      }
    }
  };

  return (
    <Form.Group className="mb-4">
      <AutoComplete
        id={inputId}
        placeholder={label}
        className={controlClasses}
        name={inputName}
        field={String(labelField)}
        value={selected}
        suggestions={filteredItems}
        completeMethod={search}
        onChange={handleChange}
        dropdown
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
