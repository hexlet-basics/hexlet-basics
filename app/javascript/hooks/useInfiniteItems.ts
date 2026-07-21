import { useEffect, useRef, useState } from "react";
import type { BlogPost } from "@/types";

interface Return<T> {
  items: T[];
  setContainerRef: (ref: HTMLElement | null, post: T) => void;
  markerRef: React.RefObject<HTMLDivElement | null>;
}

export default function useInfiniteItems<T extends BlogPost>(
  firstPost: T,
  loadNext: (lastPostId: number) => Promise<T>,
): Return<T> {
  const [items, setItems] = useState<T[]>([firstPost]);
  const postMapRef = useRef<Map<HTMLElement, T>>(new Map());

  const setContainerRef = (ref: HTMLElement | null, post: T) => {
    if (ref && !postMapRef.current.has(ref)) {
      postMapRef.current.set(ref, post);
    }
  };

  // Синхронизируем URL с постом, который сейчас находится у верха экрана.
  // Считаем геометрию по живому getBoundingClientRect на скролл, поэтому
  // выбор поста детерминирован и симметричен при прокрутке вверх и вниз
  // (старый вариант брал последний из массива пересечений и не возвращал
  // URL при скролле вверх — issue #587).
  useEffect(() => {
    let frame = 0;

    const syncUrl = () => {
      frame = 0;
      // Опорная линия у верхней части вьюпорта — активен пост, который её пересекает.
      const referenceLine = window.innerHeight * 0.25;

      let active: { post: T; top: number } | null = null;
      for (const [element, post] of postMapRef.current) {
        const { top, bottom } = element.getBoundingClientRect();
        if (top <= referenceLine && bottom > referenceLine) {
          active = { post, top };
          break;
        }
        // Фолбэк на случай зазоров: ближайший пост над линией.
        if (top <= referenceLine && (!active || top > active.top)) {
          active = { post, top };
        }
      }

      if (!active) return;
      const { url } = active.post;
      if (window.location.pathname !== new URL(url).pathname) {
        window.history.replaceState({}, "", url);
      }
    };

    const onScroll = () => {
      if (frame) return;
      frame = requestAnimationFrame(syncUrl);
    };

    window.addEventListener("scroll", onScroll, { passive: true });
    syncUrl();

    return () => {
      window.removeEventListener("scroll", onScroll);
      if (frame) cancelAnimationFrame(frame);
    };
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
          console.log("🟢 Loaded new post:", newPost.url);
          setItems((prev) => [...prev, newPost]);
        }
      },
      {
        rootMargin: "600px 0px 0px 0px",
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
