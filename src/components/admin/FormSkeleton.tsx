import { Card, Skeleton, Stack } from "@mantine/core";

// The shape of a CrudForm while the record it edits is still loading: the same
// Card, a few label-and-input pairs and the button row, so the form does not
// Jump when the data arrives.
export default function FormSkeleton({ fields = 4 }: { fields?: number }) {
  return (
    <Card withBorder p="xl" maw={720}>
      <Stack>
        {Array.from({ length: fields }, (_, index) => (
          <Stack key={index} gap={6}>
            <Skeleton h={14} w="30%" />
            <Skeleton h={36} />
          </Stack>
        ))}
        <Skeleton h={36} w={120} />
      </Stack>
    </Card>
  );
}
