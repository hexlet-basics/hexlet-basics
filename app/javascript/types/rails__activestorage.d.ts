declare module "@rails/activestorage" {
  export class DirectUpload {
    constructor(file: File, url: string, delegate?: unknown);
    create(
      callback: (
        error: Error | null,
        attributes: Record<string, unknown>,
      ) => void,
    ): void;
  }
}
