import type {
  FormDataKeys,
  FormDataType,
  Method,
  UseFormSubmitOptions,
} from "@inertiajs/core";
import {
  type InertiaFormProps,
  useForm as useInertiaForm,
} from "@inertiajs/react";
import { type ComboboxData, Loader } from "@mantine/core";
import { useDebouncedValue, useDidUpdate, useFetch } from "@mantine/hooks";
import { isNil } from "es-toolkit";
import { get } from "es-toolkit/compat";
import { type BaseSyntheticEvent, useState } from "react";

type UseAppFormOptions = {
  url: string;
  method?: Method;
  rememberKey?: string;
  onSuccess?: UseFormSubmitOptions["onSuccess"];
  onError?: UseFormSubmitOptions["onError"];
  onFinish?: UseFormSubmitOptions["onFinish"];
  onFlash?: UseFormSubmitOptions["onFlash"];
};

const toErrorMessage = (value: unknown): string => {
  if (Array.isArray(value)) {
    return value.filter(Boolean).map(String).join(", ");
  }

  if (typeof value === "string") {
    return value;
  }

  if (isNil(value)) {
    return "";
  }

  return String(value);
};

const isDateObject = (value: unknown): value is Date => value instanceof Date;

function useManagedInertiaForm<TForm extends FormDataType<TForm>>(
  data: TForm,
): InertiaFormProps<TForm>;
function useManagedInertiaForm<TForm extends FormDataType<TForm>>(
  data: TForm,
  rememberKey: string,
): InertiaFormProps<TForm>;
function useManagedInertiaForm<TForm extends FormDataType<TForm>>(
  data: TForm,
  rememberKey?: string,
): InertiaFormProps<TForm>;
function useManagedInertiaForm<TForm extends FormDataType<TForm>>(
  data: TForm,
  rememberKey?: string,
): InertiaFormProps<TForm> {
  // Inertia exposes useForm via overloads: useForm(data) and
  // useForm(rememberKey, data). We keep a single top-level hook call here to
  // satisfy React hooks rules without spreading a union tuple into useForm,
  // which TypeScript rejects.
  const args = rememberKey == null ? [data] : [rememberKey, data];

  return Reflect.apply(useInertiaForm, null, args) as InertiaFormProps<TForm>;
}

export class FormContext<
  T extends object,
  TForm extends FormDataType<TForm>,
  TKey extends keyof T & string,
> {
  private inertiaForm: InertiaFormProps<TForm>;

  constructor({
    inertiaForm,
  }: {
    inertiaForm: InertiaFormProps<TForm>;
  }) {
    this.inertiaForm = inertiaForm;
  }

  private input(key: TKey) {
    const fullKey = this.keyOf(key);
    const raw = this.valueOf(fullKey);
    const value = isNil(raw) ? "" : String(raw);

    return {
      name: fullKey,
      value,
      error: this.errorOf(fullKey),
      onChange: (e: BaseSyntheticEvent) => {
        const nextValue = String(e.currentTarget.value);
        this.setValue(fullKey, nextValue as unknown);
      },
      mb: "sm",
    };
  }

  getInputProps(key: TKey) {
    return this.input(key);
  }

  private phone(key: TKey) {
    const fullKey = this.keyOf(key);
    const rawValue = this.valueOf(fullKey);

    return {
      name: fullKey,
      value: isNil(rawValue) ? "" : String(rawValue),
      error: this.errorOf(fullKey),
      onChange: (nextValue: string) => {
        this.setValue(fullKey, nextValue);
      },
      mb: "sm",
    };
  }

  getPhoneInputProps(key: TKey) {
    return this.phone(key);
  }

  private select(key: TKey, data: ComboboxData) {
    const fullKey = this.keyOf(key);
    const raw = this.valueOf(fullKey);
    const value = isNil(raw) ? null : String(raw);

    return {
      name: fullKey,
      value,
      data,
      error: this.errorOf(fullKey),
      onChange: (val: string | null) => {
        this.setValue(fullKey, val);
      },
      mb: "sm",
    };
  }

  getSelectProps(key: TKey, data: ComboboxData) {
    return this.select(key, data);
  }

  private keyOf(key: TKey & string): FormDataKeys<TForm> {
    return key as FormDataKeys<TForm>;
  }

  private errorOf(fullKey: FormDataKeys<TForm>) {
    const directError = get(this.inertiaForm.errors, fullKey);

    if (!isNil(directError) && directError !== "") {
      return toErrorMessage(directError);
    }

    const associationKey =
      typeof fullKey === "string" && fullKey.endsWith("_id")
        ? fullKey.slice(0, -3)
        : null;

    if (associationKey == null) {
      return "";
    }

    return toErrorMessage(get(this.inertiaForm.errors, associationKey));
  }

  private valueOf(fullKey: FormDataKeys<TForm>) {
    return get(this.inertiaForm.data, fullKey);
  }

  private setValue(fullKey: FormDataKeys<TForm>, value: unknown) {
    this.inertiaForm.setData(fullKey, value as never);
  }
  private checkbox(key: TKey) {
    const fullKey = this.keyOf(key);
    const value = this.valueOf(fullKey);
    const isUnset = isNil(value);
    const checked = Boolean(value);

    if (isUnset) {
      this.setValue(fullKey, false);
    }

    return {
      name: fullKey,
      checked,
      error: this.errorOf(fullKey),
      onChange: (e: BaseSyntheticEvent) => {
        this.setValue(fullKey, e.currentTarget.checked as unknown);
      },
      mb: "sm",
    };
  }

  getCheckboxProps(key: TKey) {
    return this.checkbox(key);
  }

  private tags(key: TKey) {
    const fullKey = this.keyOf(key);
    const raw = this.valueOf(fullKey);
    const value: string[] = Array.isArray(raw)
      ? (raw as string[])
      : isNil(raw)
        ? []
        : String(raw)
            .split(",")
            .map((s) => s.trim())
            .filter((s) => s.length > 0);

    return {
      name: fullKey,
      value,
      error: this.errorOf(fullKey),
      onChange: (next: string[]) => {
        this.setValue(fullKey, next);
      },
      mb: "sm",
    };
  }

  getTagsInputProps(key: TKey) {
    return this.tags(key);
  }

  private multiSelect(key: TKey, data: ComboboxData) {
    const fullKey = this.keyOf(key);
    const raw = this.valueOf(fullKey);
    const value = Array.isArray(raw) ? raw.map((item) => String(item)) : [];

    return {
      name: fullKey,
      value,
      data,
      error: this.errorOf(fullKey),
      onChange: (next: string[]) => {
        this.setValue(fullKey, next);
      },
      mb: "sm",
    };
  }

  getMultiSelectProps(key: TKey, data: ComboboxData) {
    return this.multiSelect(key, data);
  }

  private file(key: TKey) {
    const fullKey = this.keyOf(key);
    return {
      name: fullKey,
      value: this.valueOf(fullKey) as File | null,
      error: this.errorOf(fullKey),
      onChange: (file: File | null) => {
        this.setValue(fullKey, file);
      },
      mb: "sm",
    };
  }

  getFileInputProps(key: TKey) {
    return this.file(key);
  }

  private date(key: TKey) {
    const fullKey = this.keyOf(key);
    const raw = this.valueOf(fullKey);
    const value =
      isNil(raw) || raw === "" ? null : isDateObject(raw) ? raw : String(raw);

    return {
      name: fullKey,
      value,
      error: this.errorOf(fullKey),
      onChange: (nextDate: string | null) => {
        this.setValue(fullKey, nextDate);
      },
      mb: "sm",
    };
  }

  getDateInputProps(key: TKey) {
    return this.date(key);
  }
}

type Payload = Record<string, unknown>;

export function useAppForm<TPayload extends Payload>(
  payload: TPayload,
  {
    url,
    method = "post",
    rememberKey,
    onSuccess,
    onError,
    onFinish,
    onFlash,
  }: UseAppFormOptions,
) {
  const formData = payload as FormDataType<TPayload>;
  const inertiaForm = useManagedInertiaForm(formData, rememberKey);

  // Заворачиваем в data ключ перед отправкой. Соглашение по работе с беком
  inertiaForm.transform((formData) => {
    return { data: formData as Record<string, unknown> };
  });

  const onSubmit = (e: BaseSyntheticEvent) => {
    e.preventDefault();
    const submitOptions = {
      ...(onSuccess ? { onSuccess } : {}),
      ...(onError ? { onError } : {}),
      ...(onFinish ? { onFinish } : {}),
      ...(onFlash ? { onFlash } : {}),
    };

    return inertiaForm.submit(method, url, submitOptions);
  };

  const form = new FormContext<
    TPayload,
    FormDataType<TPayload>,
    Extract<keyof TPayload, string>
  >({
    inertiaForm,
  });

  return {
    processing: inertiaForm.processing,
    onSubmit,
    submit: onSubmit,
    form,
    reset: inertiaForm.reset,
    errors: inertiaForm.errors,
    inertiaForm,
  };
}

export function useAjaxSelectProps<
  T extends object,
  TForm extends FormDataType<TForm>,
  TKey extends keyof T & string,
>(
  form: FormContext<T, TForm, TKey>,
  key: TKey,
  url: string,
  opts?: {
    minLength?: number;
    debounceMs?: number;
    clearable?: boolean;
    allowDeselect?: boolean;
    initialValueText?: string;
  },
) {
  const selectProps = form.getSelectProps(key, []);
  const value = selectProps.value;
  const resolvedValueText =
    opts?.initialValueText && opts.initialValueText.length > 0
      ? opts.initialValueText
      : typeof value === "string"
        ? value
        : "";
  const defaultData =
    typeof value === "string" && value.length > 0
      ? [{ label: resolvedValueText, value }]
      : [];

  const delay = opts?.debounceMs ?? 250;
  const minLen = opts?.minLength ?? 2;

  const [search, setSearch] = useState(resolvedValueText);
  const [q] = useDebouncedValue(search, delay);
  const { data, loading, refetch } = useFetch<ComboboxData>(
    `${url}?term=${encodeURIComponent(q)}`,
    { autoInvoke: false },
  );
  useDidUpdate(() => {
    if (q.trim().length >= minLen) {
      refetch();
    }
  }, [q, minLen]);

  return {
    ...selectProps,
    autoComplete: "off",
    data: data ?? defaultData,
    searchable: true,
    searchValue: search,
    onSearchChange: setSearch,
    onChange: (val: string | null) => {
      selectProps.onChange?.(val);
      setSearch("");
    },
    clearable: opts?.clearable ?? true,
    allowDeselect: opts?.allowDeselect ?? true,
    rightSection: loading ? <Loader size="xs" /> : null,
  };
}
