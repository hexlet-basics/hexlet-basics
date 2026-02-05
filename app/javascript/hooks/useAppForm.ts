// https://github.com/Trend-Capital/react-hook-form-mantine

import { router, usePage } from '@inertiajs/react';
import type { SelectProps } from '@mantine/core';
import { useEffect } from 'react';
import {
  type ArrayPath,
  type DefaultValues,
  type FieldValues,
  type Path,
  type PathValue,
  type UseControllerProps,
  type UseFormProps,
  useController,
  useFieldArray,
  useForm,
} from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import type { HttpRouterMethod } from '@/types';

interface Container<TForm> {
  data: DefaultValues<TForm>;
  meta?: { modelName: string };
}

interface UseAppFormOptions<
  TForm extends FieldValues,
  TContainer extends Container<TForm> = Container<TForm>,
> extends Omit<UseFormProps<TForm>, 'defaultValues'> {
  url: string;
  method?: HttpRouterMethod;
  container: TContainer;
  onSuccess?: () => void;
}

function normalizeNestedPath(
  modelName: string,
  name: string,
): { mainPath: string; fallbackPath: string } {
  const parts = name.split('.');

  const attribute = parts[parts.length - 1];
  let mainPath: string;

  if (parts.length > 1) {
    const nestedModelRaw = parts[0].replace(/_attributes$/, '');
    const nestedModel = nestedModelRaw.endsWith('s')
      ? nestedModelRaw.slice(0, -1) // qna_items → qna_item
      : nestedModelRaw;

    mainPath = `attributes.${modelName}_${nestedModel}.${attribute}`;
  } else {
    mainPath = `attributes.${modelName}.${attribute}`;
  }

  const fallbackPath = `attributes.base.${attribute}`;
  return { mainPath, fallbackPath };
}

export function useAppForm<
  // biome-ignore lint/suspicious/noExplicitAny: false positive
  TContainer extends Container<any>,
  TForm extends FieldValues = TContainer['data'],
>({
  url,
  method = 'post',
  container,
  onSuccess,
  ...options
}: UseAppFormOptions<TForm, TContainer>) {
  const modelName = container.meta?.modelName ?? '';
  if (!modelName) {
    throw new Error('Missing container.meta.modelName');
  }
  const { errors = {} } = usePage().props;

  const { t: tAr } = useTranslation('activerecord');
  const { t: tAm } = useTranslation('activemodel');

  const form = useForm<TForm>({
    defaultValues: container.data,
    ...options,
  });

  useEffect(() => {
    for (const [field, messages] of Object.entries(errors)) {
      const message = Array.isArray(messages) ? messages[0] : messages;
      form.setError(field as Path<TForm>, { message });
    }
  }, [errors, form]);

  const submit = form.handleSubmit((data) => {
    router[method](url, { [modelName]: data }, { onSuccess });
  });

  function getLabel<Name extends Path<TForm>>(
    name: Name,
    explicitLabel?: React.ReactNode | string | null,
  ): string | React.ReactNode | undefined {
    if (explicitLabel) return explicitLabel;

    const { mainPath, fallbackPath } = normalizeNestedPath(
      modelName,
      name as string,
    );
    const mainFallback = tAm(mainPath, {
      defaultValue: '',
    });
    // console.log(mainFallback)
    const value = tAr(mainPath, { defaultValue: mainFallback });
    if (value) {
      return value;
    }

    const baseFallback = tAm(fallbackPath, {
      defaultValue: '',
    });
    // console.log(baseFallback);
    const baseValue = tAr(fallbackPath, { defaultValue: baseFallback });
    if (baseValue) {
      return baseValue;
    }

    return mainPath;
  }

  function getError<Name extends Path<TForm>>(name: Name): string | undefined {
    const raw = form.formState.errors[name];
    return typeof raw?.message === 'string' ? raw.message : undefined;
  }

  function getInputProps<Name extends Path<TForm>>(
    name: Name,
    explicitLabel?: React.ReactNode,
  ) {
    const { ref, ...rest } = form.register(name);
    return {
      ...rest,
      name,
      ref,
      label: getLabel(name, explicitLabel),
      error: getError(name),
      mb: 'sm',
    };
  }

  function getSelectProps<
    Name extends Path<TForm>,
    Item extends Record<string, unknown>,
  >(
    name: Name,
    items: Item[],
    valueField: keyof Item,
    labelField: keyof Item,
    explicitLabel?: React.ReactNode,
  ) {
    // Получаем регистрацию поля (без onChange, так как мы его переопределим)
    const { name: inputName, onBlur } = form.register(name);
    const value = form.getValues(name); // Текущее значение
    const label = getLabel(name, explicitLabel);
    const error = getError(name);

    // Подготавливаем данные для Select
    const data = items.map((item) => ({
      value: item[valueField] != null ? String(item[valueField]) : '',
      label: item[labelField] != null ? String(item[labelField]) : '',
    }));

    return {
      name: inputName,
      value: value != null ? String(value) : '',
      onChange: (val: string | null) => {
        if (val === null) {
          form.setValue(name, '' as PathValue<TForm, Name>);
        } else {
          form.setValue(name, val as PathValue<TForm, Name>);
        }
        form.trigger(name); // триггерим валидацию
      },
      onBlur,
      error,
      label,
      data,
      mb: 'sm',
    } satisfies Partial<SelectProps>;
  }

  function getTagsInputProps<Name extends Path<TForm>>(
    name: Name,
    explicitLabel?: React.ReactNode,
  ) {
    const { name: inputName, onBlur, onChange } = form.register(name);
    const value = form.getValues(name);
    const error = getError(name);
    const label = getLabel(name, explicitLabel);

    return {
      name: inputName,
      value: Array.isArray(value) ? value : [],
      onChange: (val: string[]) => {
        form.setValue(name, val as PathValue<TForm, Name>);
        form.trigger(name); // триггерим валидацию поля
      },
      onBlur,
      error,
      label,
      mb: 'sm',
    };
  }

  function getFileInputProps<Name extends Path<TForm>>(
    name: Name,
    explicitLabel?: React.ReactNode,
  ) {
    const error = getError(name);
    const label = getLabel(name, explicitLabel);
    const value = form.getValues(name);

    return {
      name,
      value: value ?? null,
      onChange: (file: File | null) => {
        form.setValue(name, file as PathValue<TForm, Name>);
      },
      error,
      label,
      mb: 'sm',
    };
  }

  function useArrayField<Name extends ArrayPath<TForm>>(name: Name) {
    return useFieldArray({
      control: form.control,
      name,
      keyName: '_internalId',
    });
  }

  return {
    ...form,
    submit,
    getInputProps,
    getFileInputProps,
    getTagsInputProps,
    getSelectProps,
    useArrayField,
  };
}
