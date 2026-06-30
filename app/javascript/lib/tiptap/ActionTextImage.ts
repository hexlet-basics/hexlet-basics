import Image from "@tiptap/extension-image";

const ActionTextImage = Image.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      "data-blob-sgid": {
        default: null,
      },
    };
  },
});

export default ActionTextImage;
