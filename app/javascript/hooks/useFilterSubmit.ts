import { url } from "@/lib/utils";
import { router } from "@inertiajs/react";
import { merge } from "es-toolkit";
import qs from "qs";
import type { FormEvent } from "react";

const useFilterSubmit = (params: object) => {
  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);
    const formParams = Object.fromEntries(formData.entries());
    const serializedFormParams = qs.stringify(formParams, { encode: false });
    const parsedFormParams = qs.parse(serializedFormParams);
    const allParams = merge(params, parsedFormParams);
    // console.log(params, parsedFormParams, allParams);

    router.get(url(), allParams);
  };

  return handleSubmit;
};

export default useFilterSubmit;
