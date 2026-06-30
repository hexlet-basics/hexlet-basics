import {
  ActionIcon,
  Box,
  Input,
  type InputWrapperProps,
  Loader,
  Stack,
  Text,
} from "@mantine/core";
import { getTaskListExtension, Link, RichTextEditor } from "@mantine/tiptap";
import {
  IconBrandYoutube,
  IconFilePlus,
  IconGripVertical,
  IconTablePlus,
} from "@tabler/icons-react";
import type { Editor } from "@tiptap/core";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import Color from "@tiptap/extension-color";
import DragHandle from "@tiptap/extension-drag-handle-react";
import FileHandler from "@tiptap/extension-file-handler";
import Highlight from "@tiptap/extension-highlight";
import Placeholder from "@tiptap/extension-placeholder";
import Subscript from "@tiptap/extension-subscript";
import Superscript from "@tiptap/extension-superscript";
import { TableKit } from "@tiptap/extension-table";
import TaskItem from "@tiptap/extension-task-item";
import TaskList from "@tiptap/extension-task-list";
import TextAlign from "@tiptap/extension-text-align";
import { TextStyle } from "@tiptap/extension-text-style";
import Typography from "@tiptap/extension-typography";
import Underline from "@tiptap/extension-underline";
import Youtube from "@tiptap/extension-youtube";
import { useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import { common, createLowlight } from "lowlight";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { uploadActionTextAttachment } from "@/lib/actionTextDirectUpload";
import ActionTextImage from "@/lib/tiptap/ActionTextImage";

type BlogRichTextEditorProps = InputWrapperProps & {
  value?: string;
  defaultValue?: string;
  name?: string;
  onChange?: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
};

const lowlight = createLowlight(common);

const buildSourceCodeContent = (content: string) => {
  return `<textarea>${content}</textarea>`;
};

const getEditorValue = (editor: Editor, isSourceCodeMode: boolean) => {
  if (isSourceCodeMode) {
    return editor.getText();
  }

  return editor.getHTML();
};

const createChangeEvent = (
  nextValue: string,
): React.ChangeEvent<HTMLTextAreaElement> => {
  return {
    currentTarget: { value: nextValue },
  } as React.ChangeEvent<HTMLTextAreaElement>;
};

const buildEditorExtensions = (
  onInsertUploadedImage: (editor: Editor, file: File, pos?: number) => void,
) => {
  return [
    StarterKit.configure({
      link: false,
      codeBlock: false,
    }),
    Link.configure({
      openOnClick: false,
    }),
    CodeBlockLowlight.configure({
      lowlight,
    }),
    Underline,
    Highlight,
    TextStyle,
    Color,
    TextAlign.configure({
      types: ["heading", "paragraph"],
    }),
    Subscript,
    Superscript,
    Placeholder.configure({
      placeholder: "Write the blog post here.",
    }),
    Typography,
    TableKit.configure({
      table: {
        resizable: true,
      },
    }),
    ActionTextImage,
    getTaskListExtension(TaskList),
    TaskItem.configure({
      nested: true,
    }),
    Youtube.configure({
      controls: true,
      nocookie: true,
    }),
    FileHandler.configure({
      allowedMimeTypes: ["image/jpeg", "image/png", "image/webp", "image/gif"],
      onPaste: (editor, files) => {
        const imageFile = files.find((file) => file.type.startsWith("image/"));
        if (!imageFile) {
          return;
        }

        void onInsertUploadedImage(editor, imageFile);
      },
      onDrop: (editor, files, pos) => {
        const imageFile = files.find((file) => file.type.startsWith("image/"));
        if (!imageFile) {
          return;
        }

        void onInsertUploadedImage(editor, imageFile, pos);
      },
    }),
  ];
};

export default function BlogRichTextEditor({
  label,
  description,
  error,
  value,
  defaultValue,
  onChange,
  required,
  withAsterisk,
  ...wrapperProps
}: BlogRichTextEditorProps) {
  const [internalValue, setInternalValue] = useState(
    String(defaultValue ?? ""),
  );
  const [uploading, setUploading] = useState(false);
  const [isSourceCodeMode, setIsSourceCodeMode] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const isControlled = value !== undefined;
  const currentValue = isControlled ? String(value) : internalValue;

  const updateValue = (nextValue: string) => {
    if (!isControlled) {
      setInternalValue(nextValue);
    }

    onChange?.(createChangeEvent(nextValue));
  };

  const openImageDialog = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  const insertUploadedImage = useCallback(
    async (editorInstance: Editor, file: File, pos?: number) => {
      setUploading(true);

      try {
        const uploadedAttachment = await uploadActionTextAttachment(file);
        const imageNode = {
          type: "image",
          attrs: {
            src: uploadedAttachment.url,
            alt: file.name,
            "data-blob-sgid": uploadedAttachment.sgid,
          },
        };

        if (pos !== undefined) {
          editorInstance.chain().focus().insertContentAt(pos, imageNode).run();
          return;
        }

        editorInstance.chain().focus().insertContent(imageNode).run();
      } finally {
        setUploading(false);
      }
    },
    [],
  );

  const extensions = useMemo(() => {
    return buildEditorExtensions(insertUploadedImage);
  }, [insertUploadedImage]);

  const editor = useEditor({
    immediatelyRender: false,
    shouldRerenderOnTransaction: true,
    extensions,
    content: currentValue,
    onUpdate: ({ editor: currentEditor }) => {
      updateValue(getEditorValue(currentEditor, isSourceCodeMode));
    },
  });

  useEffect(() => {
    if (!editor) {
      return;
    }

    const editorValue = getEditorValue(editor, isSourceCodeMode);
    if (editorValue === currentValue) {
      return;
    }

    const nextContent = isSourceCodeMode
      ? buildSourceCodeContent(currentValue)
      : currentValue;

    editor.commands.setContent(nextContent, {
      emitUpdate: false,
    });
  }, [currentValue, editor, isSourceCodeMode]);

  const handleFiles = async (files: FileList | null) => {
    const file = files?.[0];
    if (!editor || !file?.type.startsWith("image/")) {
      return;
    }

    await insertUploadedImage(editor, file);
  };

  const insertTable = () => {
    editor
      ?.chain()
      .focus()
      .insertTable({ rows: 3, cols: 3, withHeaderRow: true })
      .run();
  };

  const insertYoutube = () => {
    const source = window.prompt("YouTube URL");
    if (!source) {
      return;
    }

    editor?.chain().focus().setYoutubeVideo({ src: source }).run();
  };

  return (
    <Input.Wrapper
      description={description}
      error={error}
      label={label}
      required={required}
      withAsterisk={withAsterisk}
      {...wrapperProps}
    >
      <Stack gap="xs">
        {uploading && (
          <Box c="dimmed">
            <Text component="span" size="sm">
              Uploading image
            </Text>{" "}
            <Loader component="span" size="xs" />
          </Box>
        )}

        <RichTextEditor
          editor={editor}
          onSourceCodeTextSwitch={setIsSourceCodeMode}
        >
          {editor && !isSourceCodeMode ? (
            <DragHandle editor={editor} nested>
              <ActionIcon aria-label="Drag block" size="sm" variant="subtle">
                <IconGripVertical size={16} />
              </ActionIcon>
            </DragHandle>
          ) : null}

          <RichTextEditor.Toolbar
            sticky
            stickyOffset="var(--app-shell-header-offset, 60px)"
          >
            <RichTextEditor.ControlsGroup>
              <RichTextEditor.Bold />
              <RichTextEditor.Italic />
              <RichTextEditor.Underline />
              <RichTextEditor.Strikethrough />
              <RichTextEditor.ClearFormatting />
              <RichTextEditor.Highlight />
              <RichTextEditor.Code />
            </RichTextEditor.ControlsGroup>

            <RichTextEditor.ControlsGroup>
              <RichTextEditor.H1 />
              <RichTextEditor.H2 />
              <RichTextEditor.H3 />
              <RichTextEditor.H4 />
            </RichTextEditor.ControlsGroup>

            <RichTextEditor.ControlsGroup>
              <RichTextEditor.Blockquote />
              <RichTextEditor.BulletList />
              <RichTextEditor.OrderedList />
              <RichTextEditor.TaskList />
              <RichTextEditor.CodeBlock />
              <RichTextEditor.Hr />
            </RichTextEditor.ControlsGroup>

            <RichTextEditor.ControlsGroup>
              <RichTextEditor.Link />
              <RichTextEditor.Unlink />
              <RichTextEditor.Subscript />
              <RichTextEditor.Superscript />
              <RichTextEditor.SourceCode />
            </RichTextEditor.ControlsGroup>

            <RichTextEditor.ControlsGroup>
              <RichTextEditor.AlignLeft />
              <RichTextEditor.AlignCenter />
              <RichTextEditor.AlignJustify />
              <RichTextEditor.AlignRight />
            </RichTextEditor.ControlsGroup>

            <RichTextEditor.ControlsGroup>
              <RichTextEditor.ColorPicker
                colors={[
                  "#1f2937",
                  "#c2410c",
                  "#15803d",
                  "#1d4ed8",
                  "#7c3aed",
                  "#be123c",
                ]}
              />
              <RichTextEditor.UnsetColor />
            </RichTextEditor.ControlsGroup>

            {!isSourceCodeMode && (
              <RichTextEditor.ControlsGroup>
                <RichTextEditor.Control
                  aria-label="Upload image"
                  onClick={openImageDialog}
                >
                  <IconFilePlus size={16} />
                </RichTextEditor.Control>
                <RichTextEditor.Control
                  aria-label="Insert table"
                  onClick={insertTable}
                >
                  <IconTablePlus size={16} />
                </RichTextEditor.Control>
                <RichTextEditor.Control
                  aria-label="Insert YouTube"
                  onClick={insertYoutube}
                >
                  <IconBrandYoutube size={16} />
                </RichTextEditor.Control>
              </RichTextEditor.ControlsGroup>
            )}

            <RichTextEditor.ControlsGroup>
              <RichTextEditor.Undo />
              <RichTextEditor.Redo />
            </RichTextEditor.ControlsGroup>
          </RichTextEditor.Toolbar>

          <RichTextEditor.Content mih={420} />
        </RichTextEditor>

        <input
          accept="image/*"
          hidden
          onChange={(event) => {
            void handleFiles(event.currentTarget.files);
            event.currentTarget.value = "";
          }}
          ref={fileInputRef}
          type="file"
        />
      </Stack>
    </Input.Wrapper>
  );
}
