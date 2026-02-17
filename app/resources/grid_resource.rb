class GridResource < ApplicationResource
  # sf = sort field, so = sort order, tr = total records
  attributes :sf, :so, :page, :fields, :tr, :per
  typelize sf: :string,
    so: '"asc" | "desc"',
    tr: :number,
    per: :number,
    page: :number,
    fields: "Record<string, string | number | undefined>"

  typelize :number
  attribute :first do
    (it.page - 1) * it.per
  end
end
