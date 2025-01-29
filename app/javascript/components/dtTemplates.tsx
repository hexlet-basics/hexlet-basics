import dayjs from "dayjs";

export function DTDateTemplate(data: { created_at: string }) {
  return dayjs(data.created_at).format("YYYY-MM-DD hh:mm");
}
