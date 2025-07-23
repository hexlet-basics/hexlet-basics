import { useIntersection } from '@mantine/hooks';
import { useCallback, useEffect, useRef, useState } from 'react';

type Item = { id: number; url: string };

type UseInfiniteItemsReturn<T> = {
  items: T[];
  markerRef: (el: HTMLElement | null) => void;
  setContainerRef: (el: HTMLElement | null, item: T) => void;
};

export default function useInfiniteItems<T extends Item>(
  initialItem: T,
  loadNext: (lastId: number) => Promise<T>,
): UseInfiniteItemsReturn<T> {
  const [items, setItems] = useState<T[]>([initialItem]);
  const itemMapRef = useRef<Map<Element, T>>(new Map());

  const activeUrlRef = useRef(initialItem.url);

  // Маркер подгрузки
  const { ref: markerRef, entry } = useIntersection({
    threshold: 0.5,
    rootMargin: '0px 0px 20% 0px',
  });

  const observerRef = useRef<IntersectionObserver | null>(null);

  const setContainerRef = useCallback((el: HTMLElement | null, item: T) => {
    if (!el) return;
    itemMapRef.current.set(el, item);
    observerRef.current?.observe(el);
  }, []);

  // biome-ignore lint/correctness/useExhaustiveDependencies: -
  useEffect(() => {
    if (entry?.isIntersecting) {
      (async () => {
        const lastId = items.at(-1)?.id;
        if (!lastId) return;
        try {
          const newItem = await loadNext(lastId);
          setItems((prev) => [...prev, newItem]);
          activeUrlRef.current = newItem.url;
          window.history.replaceState({}, '', newItem.url);
        } catch (e) {
          console.error('Failed to load next item:', e);
        }
      })();
    }
  }, [entry]);

  useEffect(() => {
    if (observerRef.current) return;
    observerRef.current = new IntersectionObserver(
      (entries) => {
        const visibleEntries = entries
          .filter((entry) => entry.isIntersecting)
          .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top);
        if (visibleEntries.length === 0) return;

        visibleEntries.forEach((entry) => {
          if (!entry.isIntersecting) return;
          const el = entry.target;
          const currentItem = itemMapRef.current.get(el);
          console.log('jopA!!!');
          if (currentItem && currentItem.url !== activeUrlRef.current) {
            activeUrlRef.current = currentItem.url;
            window.history.replaceState({}, '', currentItem.url);
          }
        });
      },
      { threshold: 0.5 },
    );

    return () => observerRef.current?.disconnect();
  }, []);

  return { items, markerRef, setContainerRef };
}
