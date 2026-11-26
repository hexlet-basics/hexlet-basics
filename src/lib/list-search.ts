import { zListQueryPerPage } from "@/client/zod.gen";

// The contract gives every list operation a default page size, and the generated query schemas
// carry it (`perPage: ...default(20)`), so a list route's validateSearch fills it in. Stripped
// at that value (`stripSearchParams`), the size stays out of the URL until someone changes it.
export const listSearchDefaults = { perPage: zListQueryPerPage.parse(undefined) };
