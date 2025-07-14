import { usePage } from '@inertiajs/react';
import {
  type AutocompleteProps,
  Button,
  Checkbox,
  type CheckboxProps,
  type ComboboxItem,
  Fieldset,
  FileInput,
  type FileInputProps,
  Image,
  PasswordInput,
  type PasswordInputProps,
  Select,
  type SelectProps,
  Stack,
  Textarea,
  type TextareaProps,
  TextInput,
  type TextInputProps,
} from '@mantine/core';
import axios from 'axios';
import { get } from 'es-toolkit/compat';
import { forwardRef, type PropsWithChildren, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Form,
  type FormProps,
  NestedFields,
  type NestedObject,
  useDynamicInputs,
  useInertiaInput,
} from 'use-inertia-form';
import { processHappendEvents } from '@/lib/analytics';
import type { SharedProps } from '@/types';

/* Shared Types */
export type XFormControlProps = {
  model?: string;
  field: string;
};

export interface XInputProps
  extends TextInputProps,
    XFormControlProps,
    Omit<React.ComponentPropsWithoutRef<'input'>, keyof TextInputProps> {}

export interface XPasswordInputProps
  extends PasswordInputProps,
    XFormControlProps,
    Omit<React.ComponentPropsWithoutRef<'input'>, keyof PasswordInputProps> {}

export interface XCheckboxProps
  extends CheckboxProps,
    XFormControlProps,
    Omit<React.ComponentPropsWithoutRef<'input'>, keyof CheckboxProps> {}

export interface XHiddenProps
  extends XFormControlProps,
    Omit<React.ComponentPropsWithoutRef<'input'>, 'type'> {}

export interface XFileProps extends FileInputProps, XFormControlProps {
  metaName: string;
}

export type XStateEventProps = {
  field: string;
};

/* Shared Utilities */
function useLabel(
  model: string | undefined,
  name: string,
  explicitLabel?: React.ReactNode | string | null,
): string | React.ReactNode {
  if (explicitLabel) return explicitLabel;

  // biome-ignore lint/correctness/useHookAtTopLevel: -
  const { t: tAr } = useTranslation('activerecord');
  // biome-ignore lint/correctness/useHookAtTopLevel: -
  const { t: tAm } = useTranslation('activemodel');

  const path = `attributes.${model}.${name}`;
  // @ts-expect-error todo
  const fallback = tAm(path);
  // @ts-expect-error todo
  return tAr(path, fallback);
}

function getFirstError(
  error: string | string[] | undefined,
): string | undefined {
  if (!error) return undefined;
  return Array.isArray(error) ? error[0] : error;
}

export function XForm<TForm extends NestedObject>({
  children,
  railsAttributes = true,
  className,
  ...props
}: FormProps<TForm> & PropsWithChildren) {
  const page = usePage<SharedProps>();
  const { happendEvents } = page.props;

  const [wasSubmitted, setWasSubmitted] = useState(false);

  useEffect(() => {
    if (wasSubmitted && happendEvents) {
      processHappendEvents(happendEvents);
      setWasSubmitted(false);
    }
  }, [wasSubmitted, happendEvents]);

  return (
    <Form
      className={`form ${className ?? ''}`}
      railsAttributes={railsAttributes}
      onSuccess={() => setWasSubmitted(true)}
      {...props}
    >
      {children}
    </Form>
  );
}

/* XInput */
export const XInput = forwardRef<HTMLInputElement, XInputProps>(
  ({ field, model, label, ...props }, ref) => {
    const { inputName, inputId, value, setValue, error, form } =
      useInertiaInput<string | undefined>({
        name: field,
        model,
        errorKey: field,
      });
    const preparedLabel = useLabel(form.model, field, label);
    const errorText = getFirstError(error);

    return (
      <TextInput
        mb="md"
        ref={ref}
        label={preparedLabel}
        name={inputName}
        id={inputId}
        value={value}
        error={errorText}
        onChange={(e) => setValue(e.target.value)}
        {...props}
      />
    );
  },
);
XInput.displayName = 'XInput';

/* XPasswordInput */
export const XPasswordInput = forwardRef<HTMLInputElement, XPasswordInputProps>(
  ({ field, model, label, ...props }, ref) => {
    const { inputName, inputId, value, setValue, error, form } =
      useInertiaInput<string | undefined>({
        name: field,
        model,
        errorKey: field,
      });
    const preparedLabel = useLabel(form.model, field, label);
    const errorText = getFirstError(error);

    return (
      <PasswordInput
        ref={ref}
        mb="md"
        label={preparedLabel}
        name={inputName}
        id={inputId}
        value={value}
        error={errorText}
        onChange={(e) => setValue(e.target.value)}
        {...props}
      />
    );
  },
);
XPasswordInput.displayName = 'XPasswordInput';

export interface XTextareaProps
  extends TextareaProps,
    XFormControlProps,
    Omit<React.ComponentPropsWithoutRef<'textarea'>, keyof TextareaProps> {}

export const XTextarea = forwardRef<HTMLTextAreaElement, XTextareaProps>(
  ({ field, model, label, ...props }, ref) => {
    const { inputName, inputId, value, setValue, error, form } =
      useInertiaInput<string | undefined>({
        name: field,
        model,
        errorKey: field,
      });

    const preparedLabel = useLabel(form.model, field, label);
    const errorText = getFirstError(error);

    return (
      <Textarea
        ref={ref}
        mb="md"
        label={preparedLabel}
        name={inputName}
        id={inputId}
        value={value}
        error={errorText}
        onChange={(e) => setValue(e.target.value)}
        {...props}
      />
    );
  },
);

XTextarea.displayName = 'XTextarea';

/* XCheckbox */
export const XCheckbox = forwardRef<HTMLInputElement, XCheckboxProps>(
  ({ field, model, label, ...props }, ref) => {
    const { inputName, inputId, value, setValue, error, form } =
      useInertiaInput<boolean | undefined>({
        name: field,
        model,
        errorKey: field,
      });
    const preparedLabel = useLabel(form.model, field, label);
    const errorText = getFirstError(error);

    return (
      <Checkbox
        ref={ref}
        mb="md"
        label={preparedLabel}
        name={inputName}
        id={inputId}
        checked={value}
        error={errorText}
        onChange={(e) => setValue(e.target.checked)}
        {...props}
      />
    );
  },
);
XCheckbox.displayName = 'XCheckbox';

export interface XSelectProps<
  T extends Record<string, unknown>,
  LabelKey extends keyof T = keyof T,
  ValueKey extends keyof T = keyof T,
> extends SelectProps,
    XFormControlProps {
  items?: T[];
  labelField: LabelKey;
  valueField: ValueKey;
}

export function XSelect<
  T extends Record<string, unknown>,
  LabelKey extends keyof T,
  ValueKey extends keyof T,
>(props: XSelectProps<T, LabelKey, ValueKey>) {
  const { items = [], field, model, labelField, valueField, ...rest } = props;

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput<
    string | undefined
  >({ name: field, model, errorKey: field });

  const preparedLabel = useLabel(form.model, field);
  const errorText = getFirstError(error);

  const data: ComboboxItem[] = items.map((item) => ({
    value: String(item[valueField]),
    label: String(item[labelField]),
  }));

  const handleChange = (value: string | null, option: ComboboxItem) => {
    setValue(String(option.value));
  };

  return (
    <Select
      name={inputName}
      value={String(value)}
      id={inputId}
      mb="md"
      label={preparedLabel}
      error={errorText}
      data={data}
      onChange={handleChange}
      {...rest}
    />
  );
}

// XSelect.displayName = "XSelect";

export interface XAutocompleteProps<
  T extends Record<string, unknown>,
  LabelKey extends keyof T = keyof T,
  ValueKey extends keyof T = keyof T,
> extends Omit<
      React.ComponentPropsWithoutRef<'input'>,
      keyof AutocompleteProps
    >,
    XFormControlProps {
  has?: string;
  source: string; // теперь обязательный
  labelField: LabelKey;
  valueField: ValueKey;
}

export function XAutocomplete<
  T extends Record<string, unknown>,
  LabelKey extends keyof T,
  ValueKey extends keyof T,
>(props: XAutocompleteProps<T, LabelKey, ValueKey>) {
  const { field, model, has, source, labelField, valueField, ...rest } = props;

  const [term, setTerm] = useState('');

  const { inputName, inputId, value, setValue, error, form } = useInertiaInput<
    string | undefined
  >({ name: field, model, errorKey: field });

  const [loadedItems, setLoadedItems] = useState<T[]>(() => {
    // @ts-expect-error todo
    const labelValue = form.data[form.model][has][labelField];
    const initData = (
      has ? [{ [labelField]: labelValue, [valueField]: value }] : []
    ) as T[];
    return initData;
  });

  const preparedLabel = useLabel(form.model, field);
  const errorText = getFirstError(error);

  useEffect(() => {
    const fn = async () => {
      const result = await axios.get<T[]>(source, { params: { q: term } });
      setLoadedItems(result.data);
    };

    term && term.length > 1 ? fn() : setLoadedItems([]);
  }, [term, source]);

  const data: ComboboxItem[] = loadedItems.map((item) => ({
    value: String(item[valueField]),
    label: String(item[labelField]),
  }));

  const handleChange = (value: string | null, option: ComboboxItem) => {
    setValue(String(option.value));
  };

  const handleSearchChange = (v: string) => {
    setTerm(v);
  };

  // useEffect(() => {
  //   if (has) {
  //     const related = form.data?.[model ?? ""]?.[has] as T | undefined;
  //     const relatedValue = related?.[valueField];
  //     if (relatedValue != null) {
  //       setValue(String(relatedValue));
  //     }
  //   }
  // }, [has, form.data, model, setValue, valueField]);

  return (
    <Select
      id={inputId}
      name={inputName}
      value={String(value)}
      searchable
      onChange={handleChange}
      onSearchChange={handleSearchChange}
      label={preparedLabel}
      error={errorText}
      data={data}
      mb="md"
      // rightSection={loading ? "…" : null}
      {...rest}
    />
  );
}
XAutocomplete.displayName = 'XAutocomplete';

/* XHidden */
export const XHidden = forwardRef<HTMLInputElement, XHiddenProps>(
  ({ field, model, ...props }, ref) => {
    const { inputName, inputId, value } = useInertiaInput<string | undefined>({
      name: field,
      model,
      errorKey: field,
    });

    return (
      <input
        ref={ref}
        type="hidden"
        name={inputName}
        value={value}
        id={inputId}
        {...props}
      />
    );
  },
);
XHidden.displayName = 'XHidden';

export const XFile = forwardRef<HTMLButtonElement, XFileProps>(
  ({ field, model, metaName, label, ...props }, ref) => {
    const {
      props: { railsDirectUploadsUrl },
    } = usePage<SharedProps>();
    // const { t: tAm } = useTranslation("activemodel");
    const { inputName, setValue, error, form } = useInertiaInput<
      File | undefined
    >({
      name: field,
      model,
      errorKey: field,
    });

    const imageUrl = get(form.data, ['meta', metaName]);
    const errorText = getFirstError(error);
    const preparedLabel = useLabel(form.model, field, label);

    const handleChange = (file: File | File[] | null) => {
      if (file instanceof File) {
        setValue(file);
      } else if (Array.isArray(file) && file.length > 0) {
        setValue(file[0]);
      } else {
        setValue(undefined);
      }
    };

    return (
      <Stack gap="xs" mb="md">
        <FileInput
          {...props}
          label={preparedLabel}
          name={inputName}
          data-direct-upload-url={railsDirectUploadsUrl}
          error={errorText}
          onChange={handleChange}
          placeholder={props.placeholder || preparedLabel}
          accept="image/*"
          value={form.data[inputName] as File | undefined}
          ref={ref}
        />
        {imageUrl && <Image src={imageUrl} maw={200} alt={field} />}
      </Stack>
    );
  },
);
XFile.displayName = 'XFile';

type XDynamicInputsProps = PropsWithChildren & {
  model: string;
  emptyData: Record<string, unknown>;
  label?: string;
};

export function XDynamicInputs({
  model,
  emptyData,
  label,
  children,
}: XDynamicInputsProps) {
  const { addInput, paths } = useDynamicInputs({ model, emptyData });

  return (
    <Fieldset legend={label} mb="xl" p="md" radius="md">
      <Stack gap="md">
        {paths.map((path, index) => (
          <NestedFields key={path} model={path}>
            <Stack
              gap="xs"
              p="sm"
              pl="md"
              style={{ borderLeft: '2px solid #ccc' }}
            >
              {children}
            </Stack>
          </NestedFields>
        ))}
        <Button onClick={() => addInput()} variant="light" size="xs">
          +
        </Button>
      </Stack>
    </Fieldset>
  );
}

export interface XCheckProps
  extends CheckboxProps,
    XFormControlProps,
    Omit<React.ComponentPropsWithoutRef<'input'>, keyof CheckboxProps> {}

export const XCheck = forwardRef<HTMLInputElement, XCheckProps>(
  ({ field, model, label, ...props }, ref) => {
    const { inputName, inputId, value, setValue, error, form } =
      useInertiaInput<boolean | undefined>({
        name: field,
        model,
        errorKey: field,
      });

    const preparedLabel = useLabel(form.model, field, label);

    return (
      <Checkbox
        ref={ref}
        mb="md"
        label={preparedLabel}
        name={inputName}
        id={inputId}
        checked={!!value}
        error={getFirstError(error)}
        onChange={(e) => setValue(e.target.checked)}
        {...props}
      />
    );
  },
);

XCheck.displayName = 'XCheck';
