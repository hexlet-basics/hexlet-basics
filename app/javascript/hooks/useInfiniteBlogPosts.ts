import { useEffect, useMemo, useRef, useState } from 'react';
import * as Routes from '@/routes.js';
import type { BlogPost } from '@/types';

export function useInfiniteBlogPosts(
  initialPost: BlogPost,
  loadMore: (lastPostId: number) => Promise<BlogPost>,
) {
  const [posts, setPosts] = useState<BlogPost[]>([initialPost]);
  const [activeSlug, setActiveSlug] = useState(initialPost.slug || '');
  const [isLoading, setIsLoading] = useState(false);
  const [hasError, setHasError] = useState(false);

  const markers = useRef<Record<number, Element | null>>({});
  const lastPostId = useMemo(
    () => (posts.length > 0 ? posts[posts.length - 1].id : null),
    [posts],
  );

  // Отслеживание активного поста
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        const visible = entries
          .filter((e) => e.isIntersecting)
          .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top);

        if (visible.length > 0) {
          const slug = visible[0].target.getAttribute('data-slug');
          if (slug && slug !== activeSlug) {
            setActiveSlug(slug);
            window.history.replaceState({}, '', Routes.blog_post_path(slug));
          }
        } else if (activeSlug !== posts[0].slug) {
          // Если ничего не видно, но мы пролистали вверх до первого поста
          const firstSlug = posts[0].slug;
          if (firstSlug) {
            setActiveSlug(firstSlug);
            window.history.replaceState(
              {},
              '',
              Routes.blog_post_path(firstSlug),
            );
          }
        }
      },
      { threshold: 0.5 },
    );

    const elements = document.querySelectorAll('[data-slug]');
    elements.forEach((el) => observer.observe(el));

    return () => observer.disconnect();
  }, [posts, activeSlug]);

  // Подгрузка следующего поста
  useEffect(() => {
    if (!lastPostId || hasError) return;
    const lastMarker = markers.current[lastPostId];
    if (!lastMarker) return;

    const loadObserver = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !isLoading) {
          setIsLoading(true);
          loadMore(lastPostId)
            .then((newPost) => setPosts((prev) => [...prev, newPost]))
            .catch(() => setHasError(true))
            .finally(() => setIsLoading(false));
        }
      },
      { threshold: 1, rootMargin: '0px 0px 300px 0px' },
    );

    loadObserver.observe(lastMarker);
    return () => loadObserver.disconnect();
  }, [lastPostId, loadMore, isLoading, hasError]);

  const setMarkerRef = (postId: number) => (el: Element | null) => {
    markers.current[postId] = el;
  };

  return { posts, setMarkerRef, isLoading, activeSlug, hasError };
}
