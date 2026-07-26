import { useEffect, useState } from "react";

type AsyncModuleState<T> = {
  data: T | null;
  isLoading: boolean;
  error: Error | null;
};

export default function useAsyncModule<T>(
  loader: () => Promise<T>,
): AsyncModuleState<T> {
  const [state, setState] = useState<AsyncModuleState<T>>({
    data: null,
    isLoading: true,
    error: null,
  });

  useEffect(() => {
    let cancelled = false;

    setState({
      data: null,
      isLoading: true,
      error: null,
    });

    void loader()
      .then((data) => {
        if (cancelled) return;

        setState({
          data,
          isLoading: false,
          error: null,
        });
      })
      .catch((error: unknown) => {
        if (cancelled) return;

        setState({
          data: null,
          isLoading: false,
          error: error instanceof Error ? error : new Error(String(error)),
        });
      });

    return () => {
      cancelled = true;
    };
  }, [loader]);

  return state;
}
