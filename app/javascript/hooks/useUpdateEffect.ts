import { useEffect, useRef, type EffectCallback, type DependencyList } from "react";

/**
 * A custom hook that works like useEffect but skips the initial render.
 *
 * @param effect - The effect function to run after the initial render.
 * @param dependencies - Dependency array to control when the effect runs.
 */
export default function useUpdateEffect(
  effect: EffectCallback,
  dependencies: DependencyList = [],
): void {
  const isInitialMount = useRef(true);

  // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
  useEffect(() => {
    if (isInitialMount.current) {
      isInitialMount.current = false;
    } else {
      effect();
    }
  }, dependencies);
}
