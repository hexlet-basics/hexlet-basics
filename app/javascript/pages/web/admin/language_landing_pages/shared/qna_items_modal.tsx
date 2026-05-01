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
import { buildLandingPageQnaItemFormState } from "@/lib/formBuilder";
import * as Routes from "@/routes.js";
import type { LanguageLandingPageQnaItem } from "@/types";

type Props = {
  landingPageId: number;
  initialItems: LanguageLandingPageQnaItem[];
};

type LandingPageQnaItemFormState = ReturnType<
  typeof buildLandingPageQnaItemFormState
>;

export default function QnaItemsModal({ landingPageId, initialItems }: Props) {
  const { t } = useTranslation();
  const [createOpened, createHandlers] = useDisclosure(false);
  const [editOpened, editHandlers] = useDisclosure(false);
  const [items, setItems] =
    useState<LanguageLandingPageQnaItem[]>(initialItems);
  const [editingItem, setEditingItem] =
    useState<LanguageLandingPageQnaItem | null>(null);
  const createForm = useHttpForm<
    LandingPageQnaItemFormState,
    LanguageLandingPageQnaItem
  >(buildLandingPageQnaItemFormState);
  const editForm = useHttpForm<
    LandingPageQnaItemFormState,
    LanguageLandingPageQnaItem
  >(buildLandingPageQnaItemFormState);
  const destroyForm = useHttpDestroy<LanguageLandingPageQnaItem>();

  const listUrl =
    Routes.admin_api_language_landing_page_qna_items_path(landingPageId);

  const openCreateModal = () => {
    createForm.request.setData(buildLandingPageQnaItemFormState());
    createHandlers.open();
  };

  const openEditModal = (item: LanguageLandingPageQnaItem) => {
    setEditingItem(item);
    editForm.request.setData({
      question: item.question,
      answer: item.answer,
    });
    editHandlers.open();
  };

  const closeCreateModal = () => {
    createForm.request.setData(buildLandingPageQnaItemFormState());
    createHandlers.close();
  };

  const closeEditModal = () => {
    setEditingItem(null);
    editForm.request.setData(buildLandingPageQnaItemFormState());
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
      url: Routes.admin_api_language_landing_page_qna_item_path(
        landingPageId,
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

  const confirmDeleting = useConfirmation<LanguageLandingPageQnaItem>({
    callback: async (_event, item) => {
      if (!item) return;

      await destroyForm.destroy({
        url: Routes.admin_api_language_landing_page_qna_item_path(
          landingPageId,
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
    <Fieldset>
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
