import { usePage } from '@inertiajs/react';
import { useEffect, useRef, useState } from 'react';
import type { BlogPost } from '@/types';

interface Return<T> {
  items: T[];
  setContainerRef: (ref: HTMLElement | null, post: T) => void;
  markerRef: React.RefObject<HTMLDivElement | null>;
}

export default function useInfiniteItems<T extends BlogPost>(
  firstPost: T,
  loadNext: (lastPostId: number) => Promise<T>,
): Return<T> {
  const { url: currentUrl } = usePage();
  const [items, setItems] = useState<T[]>([firstPost]);
  const postMapRef = useRef<Map<HTMLElement, T>>(new Map());
  const observerRef = useRef<IntersectionObserver | null>(null);

  const setContainerRef = (ref: HTMLElement | null, post: T) => {
    if (ref && observerRef.current && !postMapRef.current.has(ref)) {
      postMapRef.current.set(ref, post);
      observerRef.current.observe(ref);
    }
  };

  // IntersectionObserver — отслеживает появление поста в области видимости
  useEffect(() => {
    if (observerRef.current) return;

    observerRef.current = new IntersectionObserver(
      (entries) => {
        const visibleEntries = entries
          .filter((entry) => entry.isIntersecting)
          .map((entry) => ({
            entry,
            post: postMapRef.current.get(entry.target as HTMLElement),
          }));

        const lastVisible = visibleEntries.at(-1);
        if (!lastVisible) return;

        const url = lastVisible.post!.url;
        if (window.location.pathname !== new URL(url).pathname) {
          console.log('🔁 Switched to post:', url);
          window.history.replaceState({}, '', url);
        }
      },
      {
        threshold: 0,
      },
    );

    return () => observerRef.current?.disconnect();
  }, []);

  // Маркер загрузки следующих постов
  const markerRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const el = markerRef.current;
    if (!el) return;

    const observer = new IntersectionObserver(
      async ([entry]) => {
        if (entry.isIntersecting) {
          const lastPost = items[items.length - 1];
          const newPost = await loadNext(lastPost.id);
          console.log('🟢 Loaded new post:', newPost.url);
          setItems((prev) => [...prev, newPost]);
        }
      },
      {
        rootMargin: '600px 0px 0px 0px',
        threshold: 0,
      },
    );

    observer.observe(el);
    return () => observer.disconnect();
  }, [items, loadNext]);

  return {
    items,
    setContainerRef,
    markerRef,
  };
}
