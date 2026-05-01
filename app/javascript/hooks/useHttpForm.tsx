import type {
  FormDataKeys,
  FormDataType,
  Method,
  UseHttpSubmitOptions,
} from "@inertiajs/core";
import {
  HttpCancelledError,
  HttpNetworkError,
  HttpResponseError,
} from "@inertiajs/core";
import { router, useHttp } from "@inertiajs/react";
import type { ComboboxData } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { get } from "es-toolkit/compat";
import type { BaseSyntheticEvent } from "react";
import { useTranslation } from "react-i18next";

type HttpMethod = Extract<Method, "delete" | "get" | "patch" | "post" | "put">;

type InputOptions<TValue> = {
  format?: (value: TValue) => string;
  parse?: (value: string) => TValue;
};

type SelectOptions<TForm extends object, TValue> = {
  clearValue?: TValue;
  clear?: Partial<TForm>;
  format?: (value: TValue) => string | null;
  parse?: (value: string) => TValue;
};

type SubmitOptions<
  TForm extends object,
  TResponse,
  TReloadOnly extends string,
> = {
  method: HttpMethod;
  onSuccess?: (response: TResponse) => void;
  options?: UseHttpSubmitOptions<TResponse, TForm>;
  reloadOnly?: TReloadOnly[];
  successMessage?: string;
  url: string;
};

type UseHttpFormRequest<TForm extends FormDataType<TForm>> = {
  processing: boolean;
  setData: {
    (data: TForm): void;
    <TKey extends keyof TForm & string>(key: TKey, value: TForm[TKey]): void;
  };
};

type UseHttpFormResult<TForm extends FormDataType<TForm>, TResponse> = {
  getInputProps: <TKey extends keyof TForm & string>(
    key: TKey,
    options?: InputOptions<TForm[TKey]>,
  ) => {
    error: string | undefined;
    name: TKey;
    onChange: (event: BaseSyntheticEvent) => void;
    value: string;
  };
  getSelectProps: <TKey extends keyof TForm & string>(
    key: TKey,
    data: ComboboxData,
    options: SelectOptions<TForm, TForm[TKey]>,
  ) => {
    data: ComboboxData;
    error: string | undefined;
    name: TKey;
    onChange: (value: string | null) => void;
    value: string | null;
  };
  onSubmit: <TReloadOnly extends string = never>(
    event: BaseSyntheticEvent,
    options: SubmitOptions<TForm, TResponse, TReloadOnly>,
  ) => Promise<void>;
  request: UseHttpFormRequest<TForm>;
};

const defaultInputFormat = <TValue,>(value: TValue): string => {
  if (value == null) {
    return "";
  }

  return String(value);
};

const defaultInputParse = <TValue,>(value: string): TValue => value as TValue;

const defaultSelectFormat = <TValue,>(value: TValue): string | null => {
  if (value == null) {
    return null;
  }

  return String(value);
};

const defaultSelectParse = <TValue,>(value: string): TValue => value as TValue;

export function useHttpForm<
  TForm extends FormDataType<TForm>,
  TResponse = TForm,
>(initialData: TForm | (() => TForm)): UseHttpFormResult<TForm, TResponse> {
  const { t } = useTranslation();
  const request = useHttp<TForm, TResponse>(initialData);

  request.transform((data) => ({ data }));

  const getInputProps = <TKey extends keyof TForm & string>(
    key: TKey,
    options: InputOptions<TForm[TKey]> = {},
  ) => {
    const formKey = key as FormDataKeys<TForm>;
    const format = options.format ?? defaultInputFormat<TForm[TKey]>;
    const parse = options.parse ?? defaultInputParse<TForm[TKey]>;
    const rawValue = request.data[key];

    return {
      error: get(request.errors, key) as string | undefined,
      name: key,
      onChange: (event: BaseSyntheticEvent) => {
        request.setData(
          formKey,
          parse(String(event.currentTarget.value)) as never,
        );
      },
      value: format(rawValue),
    };
  };

  const getSelectProps = <TKey extends keyof TForm & string>(
    key: TKey,
    data: ComboboxData,
    options: SelectOptions<TForm, TForm[TKey]>,
  ) => {
    const formKey = key as FormDataKeys<TForm>;
    const format = options.format ?? defaultSelectFormat<TForm[TKey]>;
    const parse = options.parse ?? defaultSelectParse<TForm[TKey]>;
    const rawValue = request.data[key];

    return {
      data,
      error: get(request.errors, key) as string | undefined,
      name: key,
      onChange: (value: string | null) => {
        if (value == null) {
          request.setData(formKey, (options.clearValue ?? null) as never);
          return;
        }

        request.setData((currentData) => ({
          ...currentData,
          ...options.clear,
          [key]: parse(value),
        }));
      },
      value: format(rawValue),
    };
  };

  const onSubmit = async <TReloadOnly extends string = never>(
    event: BaseSyntheticEvent,
    options: SubmitOptions<TForm, TResponse, TReloadOnly>,
  ) => {
    event.preventDefault();

    try {
      const response = await request[options.method](
        options.url,
        options.options,
      );

      options.onSuccess?.(response);

      notifications.show({
        color: "green",
        message: options.successMessage ?? t(($) => $.common.success_message),
      });

      if (options.reloadOnly) {
        const currentUrl = window.location.href;

        router.get(
          currentUrl,
          {},
          {
            only: options.reloadOnly,
            preserveScroll: true,
            preserveState: true,
          },
        );
      }
    } catch (error: unknown) {
      if (error instanceof HttpCancelledError) {
        return;
      }

      if (error instanceof HttpNetworkError) {
        notifications.show({
          color: "red",
          message: t(($) => $.common.network_error),
        });
        return;
      }

      if (error instanceof HttpResponseError && error.response.status === 422) {
        return;
      }

      notifications.show({
        color: "red",
        message: t(($) => $.common.unexpected_error),
      });
    }
  };

  return {
    getInputProps,
    getSelectProps,
    onSubmit,
    request,
  };
}
