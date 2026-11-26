import { zListQueryPerPage } from "@/client/zod.gen";

// The contract gives a list operation a default page size, and its generated query schema carries
// it (`perPage: ...default(20)`), so the list route's validateSearch fills it in. Stripped at that
// value (`stripSearchParams`), the size stays out of the URL until someone changes it. Only for a
// route whose schema has that default: the audit trail's has none (its page is 25), and there the
// middleware would turn a chosen 20 back into 25.
export const listSearchDefaults = { perPage: zListQueryPerPage.parse(undefined) };
