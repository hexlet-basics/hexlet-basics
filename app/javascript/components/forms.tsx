import axios from "axios";
import cn from "classnames";
import { get } from "es-toolkit/compat";
import { first } from "es-toolkit/compat";
import {
  AutoComplete,
  type AutoCompleteChangeEvent,
  type AutoCompleteCompleteEvent,
} from "primereact/autocomplete";
// import { Editor } from "primereact/editor";
import { useEffect, useState, type InputHTMLAttributes } from "react";
import {
  Badge,
  Form,
  type FormCheckProps,
  type FormControlProps,
} from "react-bootstrap";
import type { AsProp } from "react-bootstrap/esm/helpers";
import { useTranslation } from "react-i18next";
import {
  type FormProps,
  Form as InertiaForm,
  type NestedObject,
  useInertiaInput,
} from "use-inertia-form";

type XFormControlProps = FormControlProps & {
  model?: string;
  name: string;
};

type XFormCheckProps = FormCheckProps & {
  model?: string;
  name: string;
};

export function XForm<TForm extends NestedObject>({
  children,
  railsAttributes = true,
  className,
  ...props
}: FormProps<TForm>) {
  // console.log(railsAttributes, props.data.review)
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

export function XInput({ name, model, as, ...props }: XFormControlProps) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tAm } = useTranslation("activemodel");

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput<
    string | undefined
  >({
    name,
    model,
  });

  const errors = error ? [error] : [];

  const path = `attributes.${form.model}.${name}`;
  const label = tAr(path, tAm(path));

  const controlClasses = cn({
    "is-invalid": errors.length > 0,
  });

  return (
    <Form.Group className="mb-4">
      <Form.FloatingLabel label={label}>
        <Form.Control
          {...props}
          as={as}
          name={inputName}
          value={value}
          id={inputId}
          placeholder={label}
          className={controlClasses}
          onChange={(e) => setValue(e.target.value)}
        />
        {errors && (
          <Form.Control.Feedback type="invalid">
            {errors.map((e) => (
              <div key={e.toString()}>{e}</div>
            ))}
          </Form.Control.Feedback>
        )}
      </Form.FloatingLabel>
    </Form.Group>
  );
}

type XFileProps = InputHTMLAttributes<HTMLInputElement> &
  AsProp & {
    fieldName: string;
    model?: string;
    name: string;
  };

export function XFile({ name, model, fieldName, ...props }: XFileProps) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tAm } = useTranslation("activemodel");

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput<
    undefined | File
  >({
    name,
    model,
  });

  const imageUrl = get(form.data, [form.model!, fieldName]);

  // console.log(form)

  const errors = error ? [error] : [];

  // const path = `attributes.${form.model}.${name}`;
  // const label = tAr(path, tAm(path));

  const controlClasses = cn({
    "is-invalid": errors.length > 0,
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setValue(e.target.files?.[0]);

  return (
    <Form.Group className="mb-4">
      <Form.Control
        type="file"
        name={inputName}
        // value={value}
        className={controlClasses}
        onChange={handleChange}
      />
      {/* {progress && ( */}
      {/*   <progress value={progress.percentage} max="100"> */}
      {/*     {progress.percentage}% */}
      {/*   </progress> */}
      {/* )} */}
      {imageUrl && (
        <div className="m-3">
          <img src={imageUrl} alt={name} />
        </div>
      )}
      {/* <Form.FloatingLabel label={label}> */}
      {/*   <Form.Control */}
      {/*     as={as} */}
      {/*     name={inputName} */}
      {/*     value={value} */}
      {/*     id={inputId} */}
      {/*     placeholder={label} */}
      {/*     className={controlClasses} */}
      {/*     onChange={(e) => setValue(e.target.value)} */}
      {/*     {...props} */}
      {/*   /> */}
      {/*   {errors && ( */}
      {/*     <Form.Control.Feedback type="invalid"> */}
      {/*       {errors.map((e) => ( */}
      {/*         <div key={e}>{e}</div> */}
      {/*       ))} */}
      {/*     </Form.Control.Feedback> */}
      {/*   )} */}
      {/* </Form.FloatingLabel> */}
    </Form.Group>
  );
}

// export function XEditor({ name, model, as, ...props }: Props) {
//   const { t: tAr } = useTranslation("activerecord");
//   const { t: tAm } = useTranslation("activemodel");
//
//   const { inputName, inputId, value, setValue, error, form } = useInertiaInput({
//     name,
//     model,
//   });
//
//   const errors = error ? [error] : [];
//
//   const path = `attributes.${form.model}.${name}`;
//   const label = tAr(path, tAm(path));
//
//   const controlClasses = cn({
//     "is-invalid": errors.length > 0,
//   });
//
//   return (
//     <Form.Group className="mb-4">
//       <Editor
//         {...props}
//         id={inputId}
//         name={inputName}
//         value={String(value)}
//         onTextChange={(e) => setValue(e.htmlValue)}
//       // style={{ height: "500px" }}
//       />
//       {/* <Form.FloatingLabel label={label}> */}
//       {/*   <Form.Control */}
//       {/*     as={as} */}
//       {/*     name={inputName} */}
//       {/*     value={value} */}
//       {/*     id={inputId} */}
//       {/*     placeholder={label} */}
//       {/*     className={controlClasses} */}
//       {/*     onChange={(e) => setValue(e.target.value)} */}
//       {/*     {...props} */}
//       {/*   /> */}
//       {/* </Form.FloatingLabel> */}
//       {errors && (
//         <Form.Control.Feedback type="invalid">
//           {errors.map((e) => (
//             <div key={e}>{e}</div>
//           ))}
//         </Form.Control.Feedback>
//       )}
//     </Form.Group>
//   );
// }

export function XCheck({ name, model, type, ...props }: XFormCheckProps) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tAm } = useTranslation("activemodel");

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput({
    name,
    model,
  });
  // console.log(form, value)

  const errors = error ? [error] : [];

  const path = `attributes.${form.model}.${name}`;
  const label = tAr(path, tAm(path));

  const controlClasses = cn({
    "is-invalid": errors.length > 0,
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
      {errors.length > 0 && (
        <Form.Control.Feedback type="invalid">
          {errors.map((e) => (
            <div key={e.toString()}>{e}</div>
          ))}
        </Form.Control.Feedback>
      )}
    </Form.Group>
  );
}

type XSelectProps<
  T extends Record<string, unknown>,
  K extends keyof T,
> = XFormControlProps & {
  items?: T[];
  has?: string; // K;
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

  const realError = has ? form.errors[`${form.model}.${has}`] : error;

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    if (has) {
      const defaultItem = get(form.data, [form.model!, has]);
      setSelected(defaultItem);
      // console.log(name, form.model, form.data, has)
    } else {
      const defaultItem = first(
        filteredItems.filter((i) => i[valueField] === value),
      );
      setSelected(defaultItem);
    }
  }, []);

  const errors = realError ? [realError] : [];

  const path = `attributes.${form.model}.${name}`;
  const label = tAr(path, tAm(path));

  const controlClasses = cn({
    "is-invalid": errors.length > 0,
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
      {errors.length > 0 && (
        <Form.Control.Feedback type="invalid">
          {errors.map((e) => (
            <div key={e.toString()}>{e}</div>
          ))}
        </Form.Control.Feedback>
      )}
    </Form.Group>
  );
}

type XStateEventProps = {
  fieldName: string;
};

export function XStateEvent({ fieldName, ...props }: XStateEventProps) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tAm } = useTranslation("activemodel");
  const { t: tCommon } = useTranslation("common");

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput({
    name: `${fieldName}_event`,
    // model,
  });

  const currentState = get(form.data, [form.model!, fieldName]);
  const stateEvents = get(
    form.data,
    [form.model!, `${fieldName}_events`],
    [],
  ) as Array<[string, string]>;

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setValue(e.target.value);
  };

  return (
    <Form.Group className="mb-4">
      <Form.Select name={inputName} id={inputId} onChange={handleChange}>
        <option>{tCommon("state_events")}</option>
        {stateEvents.map(([k, v]) => (
          <option key={v} value={v}>
            {k}
          </option>
        ))}
      </Form.Select>
      <div className="small">
        <span className="me-3">{tCommon("current_state")}:</span>
        <Badge bg="secondary">{currentState}</Badge>
      </div>
    </Form.Group>
  );
}
