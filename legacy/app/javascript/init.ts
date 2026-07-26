import type { PageProps } from "@inertiajs/core";
import dayjs from "dayjs";
import "dayjs/locale/ru";
import "dayjs/locale/es";
import duration from "dayjs/plugin/duration";
import localizedFormat from "dayjs/plugin/localizedFormat";
import relativeTime from "dayjs/plugin/relativeTime";
import * as Routes from "@/routes.js";

dayjs.extend(localizedFormat);

dayjs.extend(relativeTime);
dayjs.extend(duration);

export function configureRoutes(suffix: PageProps["suffix"]) {
  Routes.configure({
    default_url_options: {
      ...Routes.config().default_url_options,
      suffix: suffix ?? undefined,
    },
  });
}
