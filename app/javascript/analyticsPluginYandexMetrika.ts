import type { AnalyticsPlugin } from "analytics";

type YMProperties = Record<string, string | number | boolean>;

export interface YandexMetrikaPluginOptions {
  counterId: number;
  enabled?: boolean;
}

interface YandexMetrikaFunction {
  (counterId: number, command: string, ...params: unknown[]): void;
  a?: unknown[];
  l?: number;
}

declare global {
  interface Window {
    ym?: YandexMetrikaFunction;
  }
}

/**
 * Фабрика плагина для интеграции с Яндекс Метрикой.
 * Реализует методы initialize, loaded, track, page и identify.
 *
 * @param options - Опции плагина, включая counterId и флаг enabled.
 * @returns Объект, реализующий интерфейс AnalyticsPlugin.
 */
export default function yandexMetrika(
  options: YandexMetrikaPluginOptions,
): AnalyticsPlugin {
  const defaultOptions: Record<string, boolean> = {
    enabled: true,
  };

  return {
    name: "yandexMetrika",
    config: { ...defaultOptions, ...options },

    /**
     * Метод инициализации плагина.
     */
    initialize({ config }: { config: YandexMetrikaPluginOptions }) {
      const { counterId, enabled } = config;

      if (!counterId) {
        throw new Error("YandexMetrikaPlugin: counterId is required.");
      }

      // Динамически загружаем скрипт Яндекс Метрики
      ((m: Window, e: Document, t: string, r: string, i: "ym") => {
        if (!m[i]) {
          m[i] = function () {
            if (!m.ym) return;
            m.ym.a = m.ym.a || [];
            // biome-ignore lint: Стандартный код ym
            m.ym.a.push(arguments);
          };
        }
        m[i].l = new Date().getTime();
        const k = e.createElement(t) as HTMLScriptElement;
        const a = e.getElementsByTagName(t)[0];
        k.async = true;
        k.src = r;
        a.parentNode?.insertBefore(k, a);
      })(
        window,
        document,
        "script",
        "https://mc.yandex.ru/metrika/tag.js",
        "ym",
      );

      // Если ym не определена, задаём заглушку (она позже заменится загруженным скриптом)
      if (!window.ym) {
        window.ym = () => {
          // biome-ignore lint: Стандартный код ym
          (window.ym!.a = window.ym!.a || []).push(arguments);
        };
      }
      // Инициализируем Яндекс Метрику
      window.ym(counterId, "init", {
        clickmap: true,
        trackLinks: true,
        accurateTrackBounce: true,
        webvisor: true,
      });
    },

    /**
     * Отслеживает, что реальный скрипт Яндекс Метрики был успешно загружен.
     */
    loaded({ config }: { config: YandexMetrikaPluginOptions }) {
      const { enabled } = config;

      if (!enabled) {
        return;
      }

      return !!window.ym && typeof window.ym === "function";
    },
    /**
     * Отслеживает событие через Яндекс Метрику с использованием команды reachGoal.
     *
     * @param event - Название события.
     * @param properties - Дополнительные параметры.
     */
    track(event: string, properties: YMProperties = {}) {
      if (!window.ym || typeof window.ym !== "function") {
        console.warn("YandexMetrikaPlugin: ym function is not available");
        return;
      }
      window.ym(options.counterId, "reachGoal", event, properties);
    },

    /**
     * Регистрирует посещение страницы через команду hit.
     *
     * @param name - Имя страницы (не используется напрямую – используется текущий URL).
     * @param properties - Дополнительные параметры.
     */
    page(name: string, properties: YMProperties = {}) {
      if (!window.ym || typeof window.ym !== "function") {
        console.warn("YandexMetrikaPlugin: ym function is not available");
        return;
      }
      window.ym(options.counterId, "hit", window.location.href, properties);
    },

    /**
     * Устанавливает идентификатор пользователя (setUserID) и, при наличии, дополнительные параметры (userParams).
     *
     * @param userId - Идентификатор пользователя.
     * @param traits - Дополнительные атрибуты пользователя.
     */
    identify(userId: string, traits?: YMProperties) {
      if (!window.ym || typeof window.ym !== "function") {
        console.warn("YandexMetrikaPlugin: ym function is not available");
        return;
      }
      window.ym(options.counterId, "setUserID", userId);
      if (traits && Object.keys(traits).length > 0) {
        window.ym(options.counterId, "userParams", traits);
      }
    },
  };
}
