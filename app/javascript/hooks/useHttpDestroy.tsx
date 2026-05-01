import type { UseHttpSubmitOptions } from "@inertiajs/core";
import {
  HttpCancelledError,
  HttpNetworkError,
  HttpResponseError,
} from "@inertiajs/core";
import { router, useHttp } from "@inertiajs/react";
import { notifications } from "@mantine/notifications";
import { useTranslation } from "react-i18next";

type DestroyOptions<TResponse, TReloadOnly extends string> = {
  onSuccess?: (response: TResponse) => void;
  options?: UseHttpSubmitOptions<TResponse, Record<string, never>>;
  reloadOnly?: TReloadOnly[];
  successMessage?: string;
  url: string;
};

type UseHttpDestroyResult<TResponse, TReloadOnly extends string> = {
  destroy: (options: DestroyOptions<TResponse, TReloadOnly>) => Promise<void>;
};

export function useHttpDestroy<
  TResponse = unknown,
  TReloadOnly extends string = never,
>(): UseHttpDestroyResult<TResponse, TReloadOnly> {
  const { t } = useTranslation();
  const request = useHttp<Record<string, never>, TResponse>({});

  const destroy = async (options: DestroyOptions<TResponse, TReloadOnly>) => {
    try {
      const response = await request.delete(options.url, options.options);

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

  return { destroy };
}
