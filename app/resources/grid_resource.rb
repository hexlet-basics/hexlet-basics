class GridResource
  include Alba::Resource
  include Typelizer::DSL

  # sf = sort field, so = sort order, tr = total records
  attributes :sf, :so, :page, :fields, :tr, :per
  typelize sf: :string,
    so: "SortOrder",
    tr: :number,
    per: :number,
    page: :number,
    fields: "Record<string, string | number | string[] | undefined>"
end
