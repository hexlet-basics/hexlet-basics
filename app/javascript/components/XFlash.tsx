import type { SharedProps } from "@/types/types";
import { usePage } from "@inertiajs/react";

export default function XFlash() {
  const { flash } = usePage<SharedProps>().props;

  return (
    <>
      {flash.alert && <div className="alert">{flash.alert}</div>}
      {flash.notice && <div className="notice">{flash.notice}</div>}
    </>
  );
}
