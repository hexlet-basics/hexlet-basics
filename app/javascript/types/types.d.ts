export type BaseModel = {
  type: "user" | "review" | "language";
};

export type BreadcrumbItem = {
  name: string;
  url: string;
};
