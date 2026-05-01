import { move } from "@dnd-kit/helpers";
import type { DragDropProvider } from "@dnd-kit/react";
import {
  HttpCancelledError,
  HttpNetworkError,
  HttpResponseError,
} from "@inertiajs/core";
import { useHttp } from "@inertiajs/react";
import { notifications } from "@mantine/notifications";
import { type ComponentProps, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

type ReorderableRecord = {
  id: number;
};

type DragDropProviderOnDragEnd = NonNullable<
  ComponentProps<typeof DragDropProvider>["onDragEnd"]
>;

type UseReorderParams<T extends ReorderableRecord> = {
  items: T[];
  url: string;
};

type UseReorderResult<T extends ReorderableRecord> = {
  onDragEnd: DragDropProviderOnDragEnd;
  records: T[];
};

export function useReorder<T extends ReorderableRecord>({
  items,
  url,
}: UseReorderParams<T>): UseReorderResult<T> {
  const { t } = useTranslation();
  const [records, setRecords] = useState(items);
  const request = useHttp<{ ids: number[] }>({ ids: [] });

  useEffect(() => {
    setRecords(items);
  }, [items]);

  const onDragEnd: DragDropProviderOnDragEnd = (event) => {
    void (async () => {
      const previousRecords = records;
      const reorderedRecords = move(records as never, event) as T[];

      if (reorderedRecords === records) {
        return;
      }

      setRecords(reorderedRecords);

      try {
        request.transform(() => ({
          ids: reorderedRecords.map((record) => record.id),
        }));
        await request.patch(url);

        notifications.show({
          color: "green",
          message: t(($) => $.common.success_message),
        });
      } catch (error) {
        if (error instanceof HttpCancelledError) {
          setRecords(previousRecords);
          return;
        }

        if (
          error instanceof HttpNetworkError ||
          error instanceof HttpResponseError
        ) {
          setRecords(previousRecords);
          notifications.show({
            color: "red",
            message: t(($) => $.common.network_error),
          });
          return;
        }

        throw error;
      }
    })();
  };

  return {
    onDragEnd,
    records,
  };
}
