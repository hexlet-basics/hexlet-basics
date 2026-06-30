import * as Routes from "@/routes.js";

type UploadedAttachment = {
  sgid: string;
  url: string;
};

type UploadAttributes = {
  attachable_sgid: string;
  filename: string;
  signed_id: string;
};

const blobPathTemplate = Routes.rails_service_blob_proxy_path(
  "__signed_id__",
  "__filename__",
);

const createBlobPath = (signedId: string, filename: string): string => {
  return blobPathTemplate
    .replace("__signed_id__", signedId)
    .replace("__filename__", encodeURIComponent(filename));
};

export const uploadActionTextAttachment = async (
  file: File,
): Promise<UploadedAttachment> => {
  if (import.meta.env.SSR) {
    return Promise.reject(
      new Error(
        "ActionText direct upload is available only in browser runtime",
      ),
    );
  }

  const { DirectUpload } = await import("@rails/activestorage");

  return new Promise((resolve, reject) => {
    const upload = new DirectUpload(file, Routes.rails_direct_uploads_path());

    upload.create((error, attributes) => {
      if (error) {
        reject(error);
        return;
      }

      const uploadAttributes = attributes as UploadAttributes;

      resolve({
        sgid: uploadAttributes.attachable_sgid,
        url: createBlobPath(
          uploadAttributes.signed_id,
          uploadAttributes.filename,
        ),
      });
    });
  });
};
