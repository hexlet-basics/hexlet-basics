import type { SharedProps } from "@/types";
import { LanguageLandingPageCrud } from "@/types/serializers";
import { usePage } from "@inertiajs/react";
import axios from "axios";
import cn from "classnames";
import debug from "debug";
import { get } from "es-toolkit/compat";
import { first } from "es-toolkit/compat";
import {
  AutoComplete,
  type AutoCompleteChangeEvent,
  type AutoCompleteCompleteEvent,
} from "primereact/autocomplete";
// import { Editor } from "primereact/editor";
import {
  type InputHTMLAttributes,
  type PropsWithChildren,
  useEffect,
  useState,
} from "react";
import {
  Badge,
  Button,
  Form,
  type FormCheckProps,
  type FormControlProps,
} from "react-bootstrap";
import type { AsProp } from "react-bootstrap/esm/helpers";
import { useTranslation } from "react-i18next";
import {
  type FormProps,
  Form as InertiaForm,
  NestedFields,
  type NestedObject,
  useDynamicInputs,
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

const debugLog = debug("app:xform");

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
    errorKey: name,
  });

  const errors = error ? [error].flat() : [];

  // console.log(name, model, inputName)
  const path = `attributes.${form.model}.${name}`;
  // @ts-expect-error
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
              <div key={e}>{e}</div>
            ))}
          </Form.Control.Feedback>
        )}
      </Form.FloatingLabel>
    </Form.Group>
  );
}

type XFileProps = InputHTMLAttributes<HTMLInputElement> &
  AsProp & {
    metaName: string;
    model?: string;
    name: string;
  };

export function XFile({ name, model, metaName }: XFileProps) {
  const {
    props: { railsDirectUploadsUrl },
  } = usePage<SharedProps>();
  const { inputName, setValue, error, form } = useInertiaInput<
    undefined | File
  >({
    name,
    errorKey: name,
    model,
  });

  const imageUrl = get(form.data, ["meta", metaName]);
  if (!imageUrl) {
    debugLog(
      "imageUrl is not found",
      ["form.data", "meta", metaName].join("."),
    );
  }

  const errors = error ? [error].flat() : [];

  // const path = `attributes.${form.model}.${name}`;
  // const label = tAr(path, tAm(path));

  const controlClasses = cn("form-control", {
    "is-invalid": errors.length > 0,
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) =>
    setValue(e.target.files?.[0]);

  return (
    <Form.Group className="mb-4">
      {/* <Form.Control type="hidden" name={inputName} value={form.data.meta.cover_signed_id} /> */}
      <input
        type="file"
        name={inputName}
        data-direct-upload-url={railsDirectUploadsUrl}
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
    errorKey: name,
  });
  // console.log(form, value)

  const errors = error ? [error].flat() : [];

  const path = `attributes.${form.model}.${name}`;
  // @ts-expect-error
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
    errorKey: name,
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

  const errors = realError ? [realError].flat() : [];

  const path = `attributes.${form.model}.${name}`;
  // @ts-expect-error
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
      />
      {errors.length > 0 && (
        <Form.Control.Feedback type="invalid">
          {errors.map((e) => (
            <div key={e}>{e}</div>
          ))}
        </Form.Control.Feedback>
      )}
    </Form.Group>
  );
}

type XStateEventProps = {
  fieldName: string;
  // stateEvents: Array<[string, string]>;
};

export function XStateEvent({ fieldName }: XStateEventProps) {
  const { t: tAr } = useTranslation("activerecord");
  const { t: tAm } = useTranslation("activemodel");
  const { t: tCommon } = useTranslation("common");

  const { inputName, inputId, setValue, form } = useInertiaInput({
    name: `${fieldName}_event`,
    // model,
  });

  // @ts-expect-error fix
  const stateEvents = form.data.meta.state_events as Array<[string, string]>;

  const currentState = get(form.data, [form.model!, fieldName]);
  // const stateEvents = get(
  //   form.data,
  //   [form.model!, `${fieldName}_events`],
  //   [],
  // ) as Array<[string, string]>;

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

type XDynamicInputs = PropsWithChildren & {
  model: string;
  emptyData: Record<string, unknown>;
  // label: string;
};

export const XDynamicInputs = ({
  children,
  model,
  // label,
  emptyData,
}: XDynamicInputs) => {
  const { addInput, removeInput, paths } = useDynamicInputs({
    model,
    emptyData,
  });

  return (
    <div className="mb-4">
      {/* <div style={ { display: 'flex' } }> */}
      {/*   <label style={ { flex: 1 } }>{ label }</label> */}
      {/*   <button onClick={ addInput }>+</button> */}
      {/* </div> */}

      <div className="mt-4">
        {paths.map((path, i) => (
          // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
          <NestedFields key={i} model={path}>
            <Form.Group className="mb-4">
              <div>{children}</div>
              <Button variant="outline-primary" onClick={() => removeInput(i)}>
                -
              </Button>
            </Form.Group>
          </NestedFields>
        ))}
      </div>
      <Button onClick={() => addInput()}>+</Button>
    </div>
  );
};
