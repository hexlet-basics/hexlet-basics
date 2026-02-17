import type {
  Errors,
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
import dayjs from "dayjs";
import { isNil, isPlainObject } from "es-toolkit";
import { get } from "es-toolkit/compat";
import type { TFunction } from "i18next";
import { type BaseSyntheticEvent, type ReactNode, useState } from "react";
import { useTranslation } from "react-i18next";

// type AttributeArrayKeys<T> = {
//   [K in keyof T]: K extends `${string}_attributes`
//     ? T[K] extends unknown[]
//       ? K
//       : never
//     : never;
// }[keyof T];

// Тип массива по найденному ключу
// type AttributeArray<T> = T[AttributeArrayKeys<T>];

type UseAppFormOptions<TPayload extends PayloadWithMeta> = {
  url: string;
  method?: Method;
  onSuccess?: () => void;
  onError?: (errors: Errors) => void;
  onFlash?: UseFormSubmitOptions["onFlash"];
};

type ResourceMeta = {
  labels?: Record<string, string | null>;
  model?: string;
  relations?: Record<string, string>;
  [key: string]: unknown;
};

const stripMeta = (value: unknown): unknown => {
  if (Array.isArray(value)) {
    return value.map(stripMeta);
  }

  if (isPlainObject(value)) {
    const { meta, ...rest } = value;
    return Object.fromEntries(
      Object.entries(rest).map(([key, child]) => [key, stripMeta(child)]),
    );
  }

  return value;
};

type ParentChanged<
  T extends object,
  TForm extends FormDataType<TForm>,
  TKey extends keyof T & string,
> = (f: FormContext<T, TForm, TKey>) => void;

export class FormContext<
  T extends object,
  TForm extends FormDataType<TForm>,
  TKey extends keyof T & string,
> {
  index: number;
  data: T;
  private inertiaForm: InertiaFormProps<TForm>;
  private model: string;
  private relation?: string;
  // basePath — строковый префикс к ключам текущей (под)формы,
  // например: 'learning_modules_attributes.0'. Нужен, когда сабформа
  // создаётся для элемента коллекции уже с известным полным путём в data:
  // он позволяет корректно собирать имена полей/ошибок в getPrefix/fullKeyOf
  // независимо от связки relation/index и адресовать вложенные коллекции.
  private basePath?: string;
  private meta?: ResourceMeta;
  private t: TFunction;
  private relations: Record<string, string>;

  private deps = new Map<
    string,
    Map<string, ParentChanged<T, TForm, TKey> | undefined>
  >();

  constructor({
    model,
    index,
    data,
    relation,
    basePath,
    relations,
    inertiaForm,
    t,
    meta,
    // subform,
  }: {
    model: string;
    relation?: string;
    basePath?: string;
    index?: number;
    data: T;
    t: TFunction;
    relations?: Record<string, string>;
    inertiaForm: InertiaFormProps<TForm>;
    meta?: ResourceMeta;
    // subform: Subform;
  }) {
    this.data = data;
    this.meta = meta;
    this.model = model;
    this.relation = relation;
    this.basePath = basePath;
    this.index = index ?? -1;
    this.relations = relations ?? {};
    this.inertiaForm = inertiaForm;
    this.t = t;
  }

  getMetaLabel(key: TKey): string | null | undefined {
    const labels = this.meta?.labels as
      | Record<string, string | null>
      | undefined;
    return labels?.[String(key)] ?? null;
  }

  getInputProps(key: TKey) {
    const fullKey = this.fullKeyOf(key);
    const raw = this.getValue(fullKey);
    const value = isNil(raw) ? "" : String(raw);

    return {
      name: fullKey,
      value,
      error: this.getError(fullKey),
      label: this.getLabel(key),
      onChange: (e: React.BaseSyntheticEvent) => {
        this.setValue(fullKey, e.currentTarget.value);
      },
      // label: getLabel(name, explicitLabel),
      mb: "sm",
    };
  }

  getSelectProps(
    key: TKey,
    data: ComboboxData,
    opts?: {
      parent?: TKey; // базовое имя родителя в том же scope (например 'resource_type')
      parentChanged?: (f: FormContext<T, TForm, TKey>) => void;
    },
  ) {
    const fullKey = this.fullKeyOf(key);
    const raw = this.getValue(fullKey);
    const value = isNil(raw) ? null : String(raw);

    if (opts?.parent) {
      this.registerDep(opts.parent, fullKey, opts.parentChanged);
    }

    return {
      name: fullKey,
      value,
      data,
      error: this.getError(fullKey),
      label: this.getLabel(key),
      onChange: (val: string | null) => {
        this.setValue(fullKey, val);
        this.notify(fullKey);
      },
      mb: "sm",
    };
  }

  private registerDep(
    parentBaseKey: TKey,
    childFullKey: string,
    cb?: ParentChanged<T, TForm, TKey>,
  ) {
    const parentFullKey = this.fullKeyOf(parentBaseKey);
    if (!this.deps.has(parentFullKey)) this.deps.set(parentFullKey, new Map());
    this.deps.get(parentFullKey)!.set(childFullKey, cb);
  }

  private notify(parentFullKey: string) {
    const children = this.deps.get(parentFullKey);
    if (!children) return;
    for (const [childFullKey, cb] of children) {
      // 1) по умолчанию — сбросить в null
      this.inertiaForm.setData(childFullKey as never, null as never);
      // 2) если нужно — кастомная логика внутри колбэка (он сам делает f.set(...))
      cb?.(this);
    }
  }

  private fullKeyOf(key: TKey & string): FormDataKeys<TForm> {
    let result: string;
    if (this.basePath) {
      result = `${this.basePath}.${key}` as FormDataKeys<TForm>;
    } else if (this.relation) {
      result = `${this.relation}.${this.index}.${key}`;
    } else {
      result = key;
    }

    return result as FormDataKeys<TForm>;
  }

  // private getPrefix(key: TKey) {
  //   const prefix = this.basePath
  //     ? `${this.basePath}.${String(key)}`
  //     : this.relation
  //       ? `${this.relation}.${this.index}.${key}`
  //       : key;
  //   return prefix as FormDataKeys<TForm>;
  // }

  private getError(fullKey: FormDataKeys<TForm>) {
    const errors = get(this.inertiaForm.errors, fullKey);
    return Array(errors).flat().join(", ");
  }

  private getValue(fullKey: FormDataKeys<TForm>) {
    const value = get(this.inertiaForm.data, fullKey);
    return value;
  }

  private setValue(fullKey: FormDataKeys<TForm>, value: unknown) {
    this.inertiaForm.setData(fullKey, value as never);
  }

  private getLabel(
    key: TKey,
    explicitLabel?: ReactNode | string | null,
  ): string | ReactNode | undefined {
    if (explicitLabel) return explicitLabel;

    const modelsPath = `models.attributes.${this.model}.${key}`;
    const baseModelsPath = `models.attributes.base.${key}`;

    const baseTranslation = this.t(
      (resources) => get(resources, baseModelsPath) as string,
      { defaultValue: modelsPath },
    );
    const modelTranslation = this.t(
      (resources) => get(resources, modelsPath) as string,
      { defaultValue: baseTranslation },
    );
    return modelTranslation;
  }

  private deriveChildModelName(relationKey: string): string | undefined {
    if (!this.model) return;

    return `${this.model}_${relationKey}`;
  }

  // narrow to collections
  useCollection<N extends { _destroy: boolean }>(relation: TKey) {
    const fullKey = this.fullKeyOf(relation);
    const items = (this.getValue(fullKey) as Array<N>) || [];
    const activeItems = items
      .map((item, rawIndex) => ({ item, rawIndex }))
      .filter(({ item }) => !item._destroy);
    const forms = activeItems.map(({ item, rawIndex }) => {
      const childMeta = isPlainObject(item)
        ? (item as { meta?: ResourceMeta }).meta
        : undefined;
      const childModel =
        (this.relations[relation] as unknown as string) ??
        this.deriveChildModelName(String(relation));
      const form = new FormContext({
        data: item,
        model: childModel as string,
        index: rawIndex,
        relation,
        basePath: `${fullKey}.${rawIndex}` as string,
        relations: this.relations,
        t: this.t,
        inertiaForm: this.inertiaForm,
        meta: childMeta,
      });

      return form;
    });

    const add = (item: N): void => {
      const curRaw = this.getValue(fullKey) as N[] | undefined;
      const cur = Array.isArray(curRaw) ? curRaw : [];
      const next = [...cur, item];
      this.inertiaForm.setData(fullKey, next as never);
    };

    const remove = (index: number): void => {
      const curRaw = this.getValue(fullKey) as N[] | undefined;
      const cur = Array.isArray(curRaw) ? curRaw : [];
      if (index < 0 || index >= cur.length) return;
      const next = cur.map((it, i) =>
        i === index ? { ...it, _destroy: true } : it,
      );
      this.inertiaForm.setData(fullKey, next as never);
    };

    return { forms, add, remove };
  }

  getCheckboxProps(key: TKey) {
    const fullKey = this.fullKeyOf(key);
    const value = this.getValue(fullKey);
    const isUnset = isNil(value);
    const checked = Boolean(value);

    // if (!isNil(value) && typeof value !== 'boolean') {
    //   throw new Error(
    //     `getCheckboxProps supports boolean fields only, got ${typeof value}`,
    //   );
    // }

    if (isUnset) {
      this.setValue(fullKey, false);
    }

    return {
      name: fullKey,
      checked,
      error: this.getError(fullKey),
      label: this.getLabel(key),
      onChange: (e: React.BaseSyntheticEvent) => {
        this.setValue(fullKey, e.currentTarget.checked);
      },
      mb: "sm",
    };
  }

  getTagsInputProps(key: TKey) {
    const fullKey = this.fullKeyOf(key);
    const raw = this.getValue(fullKey);
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
      error: this.getError(fullKey),
      label: this.getLabel(key),
      onChange: (next: string[]) => {
        this.setValue(fullKey, next);
      },
      mb: "sm",
    };
  }

  getMultiSelectProps(key: TKey, data: ComboboxData) {
    const fullKey = this.fullKeyOf(key);
    const raw = this.getValue(fullKey) as unknown[];
    const value = raw.map((item) => String(item));
    // const value = isNil(raw) ? null : String(raw);
    // const value: string[] = Array.isArray(raw)
    //   ? (raw as string[])
    //   : isNil(raw)
    //     ? []
    //     : String(raw)
    //         .split(',')
    //         .map((s) => s.trim())
    //         .filter((s) => s.length > 0);

    return {
      name: fullKey,
      value,
      data,
      error: this.getError(fullKey),
      label: this.getLabel(key),
      onChange: (next: string[]) => {
        this.setValue(fullKey, next);
      },
      mb: "sm",
    };
  }

  getFileInputProps(key: TKey) {
    const fullKey = this.fullKeyOf(key);
    return {
      name: fullKey,
      value: this.getValue(fullKey) as File | null,
      error: this.getError(fullKey),
      label: this.getLabel(key),
      onChange: (file: File | null) => {
        this.setValue(fullKey, file);
      },
      mb: "sm",
    };
  }

  getDateInputProps(key: TKey) {
    const fullKey = this.fullKeyOf(key);
    const raw = this.getValue(fullKey);
    const value = raw ? dayjs(raw as string).toDate() : null;

    return {
      name: fullKey,
      value,
      error: this.getError(fullKey),
      label: this.getLabel(key),
      onChange: (dateStr: string | null) => {
        this.setValue(fullKey, dateStr);
      },
      mb: "sm",
    };
  }
}

type PayloadWithMeta = Record<string, unknown> & { meta?: ResourceMeta };

export function useAppForm<TPayload extends PayloadWithMeta>(
  payload: TPayload,
  {
    url,
    method = "post",
    onSuccess,
    onError,
    onFlash,
  }: UseAppFormOptions<TPayload>,
) {
  const inertiaForm = useInertiaForm(payload as FormDataType<TPayload>);

  // Заворачиваем в data ключ перед отправкой. Соглашение по работе с беком
  inertiaForm.transform((formData) => {
    return { data: stripMeta(formData) };
  });

  const { t } = useTranslation();

  const onSubmit = (e: BaseSyntheticEvent) => {
    e.preventDefault();
    const submitOptions =
      onSuccess || onError || onFlash
        ? {
            onFinish: onSuccess,
            onError,
            onFlash,
          }
        : undefined;

    return inertiaForm.submit(method, url, submitOptions);
  };

  const meta = payload.meta;
  const form = new FormContext({
    data: payload,
    meta,
    model: meta?.model ?? "",
    relations: meta?.relations ?? {},
    inertiaForm,
    t,
  });

  return {
    processing: inertiaForm.processing,
    onSubmit,
    form,
    reset: inertiaForm.reset,
    errors: inertiaForm.errors,
    inertiaForm,
  };
}

export type CollectionFormContext<
  TParent extends Record<string, unknown>,
  TChild extends Record<string, unknown>,
> = FormContext<TChild, FormDataType<TParent>, Extract<keyof TChild, string>>;

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
  },
) {
  const selectProps = form.getSelectProps(key, []);
  const value = selectProps.value;
  const metaLabel = form.getMetaLabel(key);
  const resolvedLabel =
    typeof metaLabel === "string" && metaLabel.length > 0
      ? metaLabel
      : typeof value === "string"
        ? value
        : "";
  const defaultData =
    typeof value === "string" && value.length > 0
      ? [{ label: resolvedLabel, value }]
      : [];

  const delay = opts?.debounceMs ?? 250;
  const minLen = opts?.minLength ?? 2;

  const [search, setSearch] = useState(resolvedLabel);
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
