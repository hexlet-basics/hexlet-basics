import {
  ActionIcon,
  Button,
  Fieldset,
  Group,
  Modal,
  Stack,
  Textarea,
  TextInput,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { IconPencil, IconTrash } from "@tabler/icons-react";
import { DataTable } from "mantine-datatable";
import { type BaseSyntheticEvent, useState } from "react";
import { useTranslation } from "react-i18next";
import useConfirmation from "@/hooks/useConfirmation";
import { useHttpDestroy } from "@/hooks/useHttpDestroy";
import { useHttpForm } from "@/hooks/useHttpForm";
import { buildCategoryQnaItemFormState } from "@/lib/formBuilder";
import * as Routes from "@/routes.js";
import type { LanguageCategoryQnaItem } from "@/types";

type Props = {
  categoryId: number;
  initialItems: LanguageCategoryQnaItem[];
};

type CategoryQnaItemFormState = ReturnType<
  typeof buildCategoryQnaItemFormState
>;

export default function QnaItemsModal({ categoryId, initialItems }: Props) {
  const { t } = useTranslation();
  const [createOpened, createHandlers] = useDisclosure(false);
  const [editOpened, editHandlers] = useDisclosure(false);
  const [items, setItems] = useState<LanguageCategoryQnaItem[]>(initialItems);
  const [editingItem, setEditingItem] =
    useState<LanguageCategoryQnaItem | null>(null);
  const createForm = useHttpForm<
    CategoryQnaItemFormState,
    LanguageCategoryQnaItem
  >(buildCategoryQnaItemFormState);
  const editForm = useHttpForm<
    CategoryQnaItemFormState,
    LanguageCategoryQnaItem
  >(buildCategoryQnaItemFormState);
  const destroyForm = useHttpDestroy<LanguageCategoryQnaItem>();

  const listUrl = Routes.admin_api_language_category_qna_items_path(categoryId);

  const openCreateModal = () => {
    createForm.request.setData(buildCategoryQnaItemFormState());
    createHandlers.open();
  };

  const openEditModal = (item: LanguageCategoryQnaItem) => {
    setEditingItem(item);
    editForm.request.setData({
      question: item.question,
      answer: item.answer,
    });
    editHandlers.open();
  };

  const closeCreateModal = () => {
    createForm.request.setData(buildCategoryQnaItemFormState());
    createHandlers.close();
  };

  const closeEditModal = () => {
    setEditingItem(null);
    editForm.request.setData(buildCategoryQnaItemFormState());
    editHandlers.close();
  };

  const handleCreateSubmit = async (event: BaseSyntheticEvent) => {
    await createForm.onSubmit(event, {
      method: "post",
      url: listUrl,
      onSuccess: (qnaItem) => {
        setItems((currentItems) => {
          if (!qnaItem) return currentItems;
          return [...currentItems, qnaItem];
        });
        closeCreateModal();
      },
    });
  };

  const handleEditSubmit = async (event: BaseSyntheticEvent) => {
    if (!editingItem) return;

    await editForm.onSubmit(event, {
      method: "patch",
      url: Routes.admin_api_language_category_qna_item_path(
        categoryId,
        editingItem.id,
      ),
      onSuccess: (qnaItem) => {
        setItems((currentItems) => {
          if (!qnaItem) return currentItems;

          return currentItems.map((item) =>
            item.id === qnaItem.id ? qnaItem : item,
          );
        });
        closeEditModal();
      },
    });
  };

  const confirmDeleting = useConfirmation<LanguageCategoryQnaItem>({
    callback: async (_event, item) => {
      if (!item) return;

      await destroyForm.destroy({
        url: Routes.admin_api_language_category_qna_item_path(
          categoryId,
          item.id,
        ),
        onSuccess: (qnaItem) => {
          setItems((currentItems) =>
            currentItems.filter((currentItem) => currentItem.id !== qnaItem.id),
          );
        },
      });
    },
  });

  return (
    <Fieldset p="lg" mb="xl">
      <legend>{t(($) => $.admin.language_categories.form.qna_items)}</legend>
      <Stack gap="sm">
        <Group justify="flex-end">
          <Button type="button" variant="light" onClick={openCreateModal}>
            {t(($) => $.helpers.crud.add)}
          </Button>
        </Group>
        <DataTable
          records={items}
          columns={[
            { accessor: "question", title: "Question" },
            { accessor: "answer", title: "Answer" },
            {
              accessor: "actions",
              title: "",
              textAlign: "right",
              render: (item) => (
                <Group justify="flex-end" gap="xs">
                  <ActionIcon
                    variant="subtle"
                    onClick={() => openEditModal(item)}
                  >
                    <IconPencil size={16} />
                  </ActionIcon>
                  <ActionIcon
                    variant="subtle"
                    color="red"
                    onClick={(event) => confirmDeleting(event, item)}
                  >
                    <IconTrash size={16} />
                  </ActionIcon>
                </Group>
              ),
            },
          ]}
        />
      </Stack>
      <Modal opened={createOpened} onClose={closeCreateModal}>
        <form onSubmit={handleCreateSubmit}>
          <TextInput
            label={t(($) => $.models.attributes.base.question)}
            {...createForm.getInputProps("question")}
            mb="sm"
          />
          <Textarea
            label={t(($) => $.models.attributes.base.answer)}
            {...createForm.getInputProps("answer")}
            rows={5}
            mb="sm"
          />
          <Button type="submit" loading={createForm.request.processing}>
            {t(($) => $.helpers.submit.save)}
          </Button>
        </form>
      </Modal>
      <Modal opened={editOpened} onClose={closeEditModal}>
        <form onSubmit={handleEditSubmit}>
          <TextInput
            label={t(($) => $.models.attributes.base.question)}
            {...editForm.getInputProps("question")}
            mb="sm"
          />
          <Textarea
            label={t(($) => $.models.attributes.base.answer)}
            {...editForm.getInputProps("answer")}
            rows={5}
            mb="sm"
          />
          <Button type="submit" loading={editForm.request.processing}>
            {t(($) => $.helpers.submit.save)}
          </Button>
        </form>
      </Modal>
    </Fieldset>
  );
}
